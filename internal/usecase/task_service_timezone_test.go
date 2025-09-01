package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/usecase"
)

func TestTimezoneAwareScheduling(t *testing.T) {
	ctx := context.Background()
	taskRepo := &mockTaskRepo{tasks: make(map[string]*entity.Task)}
	userRepo := &mockUserRepo{users: map[int64]*entity.User{1: {ID: 1}}}
	uuidGen := &mockUUIDGen{}
	mockScheduler := new(mockScheduler)

	service := usecase.NewTaskService(taskRepo, userRepo, uuidGen, mockScheduler)

	t.Run("Task with timezone uses local time", func(t *testing.T) {
		task := &entity.Task{
			ID:          "tz-task",
			Title:       "Timezone Test",
			Category:    "daily",
			TimeZone:    "America/New_York",
			StreakCount: 0,
		}
		taskRepo.tasks[task.ID] = task

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

		err := service.CompleteTask(ctx, 1, "tz-task")
		require.NoError(t, err)
		mockScheduler.AssertExpectations(t)
	})

	t.Run("Task without timezone uses UTC", func(t *testing.T) {
		task := &entity.Task{
			ID:          "no-tz-task",
			Title:       "No Timezone",
			Category:    "daily",
			StreakCount: 0,
		}
		taskRepo.tasks[task.ID] = task

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

		err := service.CompleteTask(ctx, 1, "no-tz-task")
		require.NoError(t, err)
		mockScheduler.AssertExpectations(t)
	})
}
