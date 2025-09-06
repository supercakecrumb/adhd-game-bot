package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
	"github.com/supercakecrumb/adhd-game-bot/internal/ports"
)

type QuestService struct {
	questRepo       ports.QuestRepository
	userRepo        ports.UserRepository
	uuidGen         ports.UUIDGenerator
	scheduler       ports.Scheduler
	idempotencyRepo ports.IdempotencyRepository
	txManager       ports.TxManager
}

func NewQuestService(
	questRepo ports.QuestRepository,
	userRepo ports.UserRepository,
	uuidGen ports.UUIDGenerator,
	scheduler ports.Scheduler,
	idempotencyRepo ports.IdempotencyRepository,
	txManager ports.TxManager,
) *QuestService {
	return &QuestService{
		questRepo:       questRepo,
		userRepo:        userRepo,
		uuidGen:         uuidGen,
		scheduler:       scheduler,
		idempotencyRepo: idempotencyRepo,
		txManager:       txManager,
	}
}

type CreateQuestInput struct {
	Title            string
	Description      string
	Category         string
	Difficulty       string
	Mode             string
	PointsAward      valueobject.Decimal
	RatePointsPerMin *valueobject.Decimal
	MinMinutes       *int
	MaxMinutes       *int
	DailyPointsCap   *valueobject.Decimal
	CooldownSec      int
	StreakEnabled    bool
	Status           string
	TimeZone         string
}

func (s *QuestService) CreateQuest(ctx context.Context, userID int64, dungeonID string, input CreateQuestInput) (*entity.Quest, error) {
	// Validate user exists
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Create quest entity
	quest := &entity.Quest{
		ID:               s.uuidGen.New(),
		DungeonID:        dungeonID,
		Title:            input.Title,
		Description:      input.Description,
		Category:         input.Category,
		Difficulty:       input.Difficulty,
		Mode:             input.Mode,
		PointsAward:      input.PointsAward,
		RatePointsPerMin: input.RatePointsPerMin,
		MinMinutes:       input.MinMinutes,
		MaxMinutes:       input.MaxMinutes,
		DailyPointsCap:   input.DailyPointsCap,
		CooldownSec:      input.CooldownSec,
		StreakEnabled:    input.StreakEnabled,
		Status:           input.Status,
		TimeZone:         input.TimeZone,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Create quest
	err = s.questRepo.Create(ctx, quest)
	if err != nil {
		return nil, err
	}

	// Schedule recurring quest if needed
	if quest.Category == "daily" || quest.Category == "weekly" {
		err = s.scheduler.ScheduleRecurringTask(ctx, quest)
		if err != nil {
			return nil, err
		}
	}

	return quest, nil
}

func (s *QuestService) GetQuest(ctx context.Context, questID string) (*entity.Quest, error) {
	return s.questRepo.GetByID(ctx, questID)
}

func (s *QuestService) UpdateQuest(ctx context.Context, quest *entity.Quest) error {
	// Verify quest exists
	_, err := s.questRepo.GetByID(ctx, quest.ID)
	if err != nil {
		return err
	}

	quest.UpdatedAt = time.Now()
	return s.questRepo.Update(ctx, quest)
}

type CompleteQuestInput struct {
	IdempotencyKey  string
	CompletionRatio *float64 // For PARTIAL mode
	Minutes         *int     // For PER_MINUTE mode
}

func (s *QuestService) CompleteQuest(ctx context.Context, userID int64, questID string, input CompleteQuestInput) error {
	// Create idempotency key
	idempKey := &entity.IdempotencyKey{
		Key:       input.IdempotencyKey,
		Operation: "quest_complete",
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
		// Get quest
		quest, err := s.questRepo.GetByID(ctx, questID)
		if err != nil {
			return err
		}

		// Verify user exists
		_, err = s.userRepo.FindByID(ctx, userID)
		if err != nil {
			return err
		}

		// Get time in quest's timezone
		var now time.Time
		if quest.TimeZone != "" {
			loc, err := time.LoadLocation(quest.TimeZone)
			if err != nil {
				return err
			}
			now = time.Now().In(loc)
		} else {
			// Use UTC when no timezone is specified
			now = time.Now().UTC()
		}

		// Update quest completion
		quest.LastCompletedAt = &now
		quest.StreakCount++

		err = s.questRepo.Update(ctx, quest)
		if err != nil {
			return err
		}

		// Reschedule recurring quest
		if quest.Category == "daily" || quest.Category == "weekly" {
			err = s.scheduler.ScheduleRecurringTask(ctx, quest)
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

func (s *QuestService) ListQuests(ctx context.Context, userID int64, dungeonID string) ([]*entity.Quest, error) {
	// Verify user exists
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.questRepo.ListByDungeon(ctx, dungeonID)
}
