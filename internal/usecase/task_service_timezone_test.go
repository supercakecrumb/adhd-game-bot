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

	t.Run("Task with timezone uses local time", func(t *testing.T) {
		taskRepo := new(testhelpers.MockTaskRepository)
		userRepo := &mockUserRepo{users: map[int64]*entity.User{1: {ID: 1, ChatID: 1}}}
		uuidGen := &mockUUIDGen{}
		mockScheduler := new(mockScheduler)
		mockIdempotencyRepo := new(testhelpers.MockIdempotencyRepository)
		mockTxManager := new(mockTxManager)

		task := &entity.Task{
			ID:          "tz-task",
			Title:       "Timezone Test",
			Category:    "daily",
			TimeZone:    "America/New_York",
			StreakCount: 0,
		}

		taskRepo.On("FindByID", ctx, "tz-task").Return(task, nil)
		taskRepo.On("Update", ctx, mock.AnythingOfType("*entity.Task")).Return(nil)

		service := usecase.NewTaskService(taskRepo, userRepo, uuidGen, mockScheduler, mockIdempotencyRepo, mockTxManager)

		// Mock idempotency
		mockIdempotencyRepo.On("FindByKey", ctx, mock.AnythingOfType("string")).Return(nil, ports.ErrIdempotencyKeyNotFound).Once()
		mockIdempotencyRepo.On("Create", ctx, mock.AnythingOfType("*entity.IdempotencyKey")).Return(nil).Once()
		mockIdempotencyRepo.On("Update", ctx, mock.AnythingOfType("*entity.IdempotencyKey")).Return(nil).Once()

		// Mock transaction
		mockTxManager.On("WithTx", ctx, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			// Expect scheduler to be called with timezone-adjusted task
			mockScheduler.On("ScheduleRecurringTask", ctx, mock.AnythingOfType("*entity.Task")).
				Return(nil).
				Run(func(args mock.Arguments) {
					scheduledTask := args.Get(1).(*entity.Task)
					require.Equal(t, "America/New_York", scheduledTask.TimeZone)
					require.NotNil(t, scheduledTask.LastCompletedAt)
					_, offset := scheduledTask.LastCompletedAt.Zone()
					nyLoc, _ := time.LoadLocation("America/New_York")
					_, expectedOffset := time.Now().In(nyLoc).Zone()
					require.Equal(t, expectedOffset, offset)
				})
			fn(ctx)
		}).Return(nil).Once()

		err := service.CompleteTask(ctx, 1, "tz-task")
		require.NoError(t, err)
		mockScheduler.AssertExpectations(t)
		mockIdempotencyRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})

	t.Run("Task without timezone uses UTC", func(t *testing.T) {
		taskRepo := new(testhelpers.MockTaskRepository)
		userRepo := &mockUserRepo{users: map[int64]*entity.User{1: {ID: 1, ChatID: 1}}}
		uuidGen := &mockUUIDGen{}
		mockScheduler := new(mockScheduler)
		mockIdempotencyRepo := new(testhelpers.MockIdempotencyRepository)
		mockTxManager := new(mockTxManager)

		task := &entity.Task{
			ID:          "no-tz-task",
			Title:       "No Timezone",
			Category:    "daily",
			StreakCount: 0,
		}

		taskRepo.On("FindByID", ctx, "no-tz-task").Return(task, nil)
		taskRepo.On("Update", ctx, mock.AnythingOfType("*entity.Task")).Return(nil)

		service := usecase.NewTaskService(taskRepo, userRepo, uuidGen, mockScheduler, mockIdempotencyRepo, mockTxManager)

		// Mock idempotency
		mockIdempotencyRepo.On("FindByKey", ctx, mock.AnythingOfType("string")).Return(nil, ports.ErrIdempotencyKeyNotFound).Once()
		mockIdempotencyRepo.On("Create", ctx, mock.AnythingOfType("*entity.IdempotencyKey")).Return(nil).Once()
		mockIdempotencyRepo.On("Update", ctx, mock.AnythingOfType("*entity.IdempotencyKey")).Return(nil).Once()

		// Mock transaction
		mockTxManager.On("WithTx", ctx, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			// Expect scheduler to be called with UTC time
			mockScheduler.On("ScheduleRecurringTask", ctx, mock.AnythingOfType("*entity.Task")).
				Return(nil).
				Run(func(args mock.Arguments) {
					scheduledTask := args.Get(1).(*entity.Task)
					require.Empty(t, scheduledTask.TimeZone)
					require.NotNil(t, scheduledTask.LastCompletedAt)
					_, offset := scheduledTask.LastCompletedAt.Zone()
					require.Equal(t, 0, offset) // UTC
				})
			fn(ctx)
		}).Return(nil).Once()

		err := service.CompleteTask(ctx, 1, "no-tz-task")
		require.NoError(t, err)
		mockScheduler.AssertExpectations(t)
		mockIdempotencyRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})
}
