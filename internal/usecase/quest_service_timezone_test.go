package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/ports"
	"github.com/supercakecrumb/adhd-game-bot/internal/usecase"
	"github.com/supercakecrumb/adhd-game-bot/internal/usecase/testhelpers"
)

func TestTimezoneAwareScheduling(t *testing.T) {
	ctx := context.Background()

	t.Run("Quest with timezone uses local time", func(t *testing.T) {
		questRepo := new(testhelpers.MockQuestRepository)
		userRepo := new(testhelpers.MockUserRepository)
		uuidGen := new(testhelpers.MockUUIDGenerator)
		mockScheduler := new(testhelpers.MockScheduler)
		mockIdempotencyRepo := new(testhelpers.MockIdempotencyRepository)
		mockTxManager := new(testhelpers.MockTxManager)

		quest := &entity.Quest{
			ID:          "tz-quest",
			Title:       "Timezone Test",
			Category:    "daily",
			TimeZone:    "America/New_York",
			StreakCount: 0,
		}

		questRepo.On("GetByID", ctx, "tz-quest").Return(quest, nil)
		questRepo.On("Update", ctx, mock.AnythingOfType("*entity.Quest")).Return(nil)
		userRepo.On("FindByID", ctx, int64(1)).Return(&entity.User{ID: 1}, nil)

		service := usecase.NewQuestService(questRepo, userRepo, uuidGen, mockScheduler, mockIdempotencyRepo, mockTxManager)

		// Mock idempotency
		mockIdempotencyRepo.On("FindByKey", ctx, mock.AnythingOfType("string")).Return(nil, ports.ErrIdempotencyKeyNotFound).Once()
		mockIdempotencyRepo.On("Create", ctx, mock.AnythingOfType("*entity.IdempotencyKey")).Return(nil).Once()
		mockIdempotencyRepo.On("Update", ctx, mock.AnythingOfType("*entity.IdempotencyKey")).Return(nil).Once()

		// Mock transaction
		mockTxManager.On("WithTx", ctx, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			// Expect scheduler to be called with timezone-adjusted quest
			mockScheduler.On("ScheduleRecurringTask", ctx, mock.AnythingOfType("*entity.Quest")).
				Return(nil).
				Run(func(args mock.Arguments) {
					scheduledQuest := args.Get(1).(*entity.Quest)
					require.Equal(t, "America/New_York", scheduledQuest.TimeZone)
					require.NotNil(t, scheduledQuest.LastCompletedAt)
					_, offset := scheduledQuest.LastCompletedAt.Zone()
					nyLoc, _ := time.LoadLocation("America/New_York")
					_, expectedOffset := time.Now().In(nyLoc).Zone()
					require.Equal(t, expectedOffset, offset)
				})
			fn(ctx)
		}).Return(nil).Once()

		input := usecase.CompleteQuestInput{
			IdempotencyKey: "key-1",
		}
		err := service.CompleteQuest(ctx, 1, "tz-quest", input)
		require.NoError(t, err)
		mockScheduler.AssertExpectations(t)
		mockIdempotencyRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})

	t.Run("Quest without timezone uses UTC", func(t *testing.T) {
		questRepo := new(testhelpers.MockQuestRepository)
		userRepo := new(testhelpers.MockUserRepository)
		uuidGen := new(testhelpers.MockUUIDGenerator)
		mockScheduler := new(testhelpers.MockScheduler)
		mockIdempotencyRepo := new(testhelpers.MockIdempotencyRepository)
		mockTxManager := new(testhelpers.MockTxManager)

		quest := &entity.Quest{
			ID:          "no-tz-quest",
			Title:       "No Timezone",
			Category:    "daily",
			StreakCount: 0,
		}

		questRepo.On("GetByID", ctx, "no-tz-quest").Return(quest, nil)
		questRepo.On("Update", ctx, mock.AnythingOfType("*entity.Quest")).Return(nil)
		userRepo.On("FindByID", ctx, int64(1)).Return(&entity.User{ID: 1}, nil)

		service := usecase.NewQuestService(questRepo, userRepo, uuidGen, mockScheduler, mockIdempotencyRepo, mockTxManager)

		// Mock idempotency
		mockIdempotencyRepo.On("FindByKey", ctx, mock.AnythingOfType("string")).Return(nil, ports.ErrIdempotencyKeyNotFound).Once()
		mockIdempotencyRepo.On("Create", ctx, mock.AnythingOfType("*entity.IdempotencyKey")).Return(nil).Once()
		mockIdempotencyRepo.On("Update", ctx, mock.AnythingOfType("*entity.IdempotencyKey")).Return(nil).Once()

		// Mock transaction
		mockTxManager.On("WithTx", ctx, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			// Expect scheduler to be called with UTC time
			mockScheduler.On("ScheduleRecurringTask", ctx, mock.AnythingOfType("*entity.Quest")).
				Return(nil).
				Run(func(args mock.Arguments) {
					scheduledQuest := args.Get(1).(*entity.Quest)
					require.Empty(t, scheduledQuest.TimeZone)
					require.NotNil(t, scheduledQuest.LastCompletedAt)
					_, offset := scheduledQuest.LastCompletedAt.Zone()
					require.Equal(t, 0, offset) // UTC
				})
			fn(ctx)
		}).Return(nil).Once()

		input := usecase.CompleteQuestInput{
			IdempotencyKey: "key-2",
		}
		err := service.CompleteQuest(ctx, 1, "no-tz-quest", input)
		require.NoError(t, err)
		mockScheduler.AssertExpectations(t)
		mockIdempotencyRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})
}
