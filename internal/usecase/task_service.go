package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/ports"
)

type TaskService struct {
	taskRepo        ports.TaskRepository
	userRepo        ports.UserRepository
	uuidGen         ports.UUIDGenerator
	scheduler       ports.Scheduler
	idempotencyRepo ports.IdempotencyRepository
	txManager       ports.TxManager
}

func NewTaskService(
	taskRepo ports.TaskRepository,
	userRepo ports.UserRepository,
	uuidGen ports.UUIDGenerator,
	scheduler ports.Scheduler,
	idempotencyRepo ports.IdempotencyRepository,
	txManager ports.TxManager,
) *TaskService {
	return &TaskService{
		taskRepo:        taskRepo,
		userRepo:        userRepo,
		uuidGen:         uuidGen,
		scheduler:       scheduler,
		idempotencyRepo: idempotencyRepo,
		txManager:       txManager,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, userID int64, task *entity.Task) (*entity.Task, error) {
	// Validate user exists
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Generate task ID
	task.ID = s.uuidGen.New()

	// Create task
	err = s.taskRepo.Create(ctx, task)
	if err != nil {
		return nil, err
	}

	// Schedule recurring task if needed
	if task.Category == "daily" || task.Category == "weekly" {
		err = s.scheduler.ScheduleRecurringTask(ctx, task)
		if err != nil {
			return nil, err
		}
	}

	return task, nil
}

func (s *TaskService) GetTask(ctx context.Context, taskID string) (*entity.Task, error) {
	return s.taskRepo.FindByID(ctx, taskID)
}

func (s *TaskService) UpdateTask(ctx context.Context, task *entity.Task) error {
	// Verify task exists
	_, err := s.taskRepo.FindByID(ctx, task.ID)
	if err != nil {
		return err
	}

	return s.taskRepo.Update(ctx, task)
}

func (s *TaskService) CompleteTask(ctx context.Context, userID int64, taskID string) error {
	// Create idempotency key
	idempKey := &entity.IdempotencyKey{
		Key:       s.generateIdempotencyKey(userID, taskID, "complete"),
		Operation: "task_complete",
		UserID:    userID,
		Status:    "pending",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour), // Expire after 24 hours
	}

	// Check if operation already exists
	existingKey, err := s.idempotencyRepo.FindByKey(ctx, idempKey.Key)
	if err == nil && existingKey != nil {
		if existingKey.IsCompleted() {
			// Operation already completed, return success
			return nil
		}
		if !existingKey.IsExpired() {
			// Operation is still pending
			return errors.New("operation in progress")
		}
	}

	// Create idempotency key
	err = s.idempotencyRepo.Create(ctx, idempKey)
	if err != nil && err != ports.ErrIdempotencyKeyExists {
		return err
	}

	// Execute the operation in a transaction
	err = s.txManager.WithTx(ctx, func(ctx context.Context) error {
		// Get task
		task, err := s.taskRepo.FindByID(ctx, taskID)
		if err != nil {
			return err
		}

		// Verify user exists
		_, err = s.userRepo.FindByID(ctx, userID)
		if err != nil {
			return err
		}

		// Get time in task's timezone
		var now time.Time
		if task.TimeZone != "" {
			loc, err := time.LoadLocation(task.TimeZone)
			if err != nil {
				return err
			}
			now = time.Now().In(loc)
		} else {
			// Use UTC when no timezone is specified
			now = time.Now().UTC()
		}

		// Update task completion
		task.LastCompletedAt = &now
		task.StreakCount++

		err = s.taskRepo.Update(ctx, task)
		if err != nil {
			return err
		}

		// Reschedule recurring task
		if task.Category == "daily" || task.Category == "weekly" {
			err = s.scheduler.ScheduleRecurringTask(ctx, task)
			if err != nil {
				return err
			}
		}

		return nil
	})

	// Update idempotency key status
	completedAt := time.Now()
	idempKey.CompletedAt = &completedAt
	if err != nil {
		idempKey.Status = "failed"
		idempKey.Result = err.Error()
	} else {
		idempKey.Status = "completed"
		idempKey.Result = "success"
	}

	updateErr := s.idempotencyRepo.Update(ctx, idempKey)
	if updateErr != nil {
		// Log error but don't fail the operation
		// In production, you'd want to log this
	}

	return err
}

func (s *TaskService) generateIdempotencyKey(userID int64, taskID string, operation string) string {
	return fmt.Sprintf("%d-%s-%s-%d", userID, taskID, operation, time.Now().Unix())
}

func (s *TaskService) ListTasksByUser(ctx context.Context, userID int64) ([]*entity.Task, error) {
	// Verify user exists
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.taskRepo.FindByUser(ctx, userID)
}
