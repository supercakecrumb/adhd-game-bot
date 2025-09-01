package usecase

import (
	"context"
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/ports"
)

type TaskService struct {
	taskRepo  ports.TaskRepository
	userRepo  ports.UserRepository
	uuidGen   ports.UUIDGenerator
	scheduler ports.Scheduler
}

func NewTaskService(
	taskRepo ports.TaskRepository,
	userRepo ports.UserRepository,
	uuidGen ports.UUIDGenerator,
	scheduler ports.Scheduler,
) *TaskService {
	return &TaskService{
		taskRepo:  taskRepo,
		userRepo:  userRepo,
		uuidGen:   uuidGen,
		scheduler: scheduler,
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
		now = time.Now()
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
		return s.scheduler.ScheduleRecurringTask(ctx, task)
	}

	return nil
}

func (s *TaskService) ListTasksByUser(ctx context.Context, userID int64) ([]*entity.Task, error) {
	// Verify user exists
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.taskRepo.FindByUser(ctx, userID)
}
