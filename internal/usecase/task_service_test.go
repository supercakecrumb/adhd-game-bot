package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/ports"
	"github.com/supercakecrumb/adhd-game-bot/internal/usecase"
	"github.com/supercakecrumb/adhd-game-bot/internal/usecase/testhelpers"
)

type mockUUIDGen struct{}

func (m *mockUUIDGen) New() string {
	return "generated-uuid"
}

func TestTaskService(t *testing.T) {
	ctx := context.Background()
	taskRepo := new(testhelpers.MockTaskRepository)
	userRepo := new(testhelpers.MockUserRepository)

	uuidGen := &mockUUIDGen{}
	mockScheduler := new(testhelpers.MockScheduler)
	mockIdempotencyRepo := new(testhelpers.MockIdempotencyRepository)
	mockTxManager := new(testhelpers.MockTxManager)

	service := usecase.NewTaskService(taskRepo, userRepo, uuidGen, mockScheduler, mockIdempotencyRepo, mockTxManager)

	t.Run("CreateTask", func(t *testing.T) {
		task := &entity.Task{
			Title:       "Test Task",
			Description: "Test Description",
			Category:    "daily",
		}

		// Setup mock expectations
		taskRepo.On("Create", ctx, mock.AnythingOfType("*entity.Task")).Return(nil).Once()
		userRepo.On("FindByID", ctx, int64(1)).Return(&entity.User{ID: 1}, nil).Once()

		// Mock the scheduler call for daily tasks
		mockScheduler.On("ScheduleRecurringTask", ctx, mock.AnythingOfType("*entity.Task")).Return(nil).Once()

		created, err := service.CreateTask(ctx, 1, task)
		require.NoError(t, err)
		require.Equal(t, "generated-uuid", created.ID)
		require.Equal(t, "Test Task", created.Title)

		mockScheduler.AssertExpectations(t)
		mockIdempotencyRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})

	t.Run("CompleteTask", func(t *testing.T) {
		task := &entity.Task{
			ID:          "task-1",
			Title:       "Complete Me",
			StreakCount: 0,
			Category:    "daily",
		}
		taskRepo.On("FindByID", ctx, "task-1").Return(task, nil)
		taskRepo.On("Update", ctx, mock.AnythingOfType("*entity.Task")).Return(nil)
		userRepo.On("FindByID", ctx, int64(1)).Return(&entity.User{ID: 1}, nil)

		// Mock idempotency check
		mockIdempotencyRepo.On("FindByKey", ctx, mock.AnythingOfType("string")).Return(nil, ports.ErrIdempotencyKeyNotFound).Once()
		mockIdempotencyRepo.On("Create", ctx, mock.AnythingOfType("*entity.IdempotencyKey")).Return(nil).Once()
		mockIdempotencyRepo.On("Update", ctx, mock.AnythingOfType("*entity.IdempotencyKey")).Return(nil).Once()

		// Mock transaction
		mockTxManager.On("WithTx", ctx, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			// Mock the scheduler call for daily tasks
			mockScheduler.On("ScheduleRecurringTask", ctx, mock.AnythingOfType("*entity.Task")).Return(nil)
			fn(ctx)
		}).Return(nil).Once()

		err := service.CompleteTask(ctx, 1, "task-1")
		require.NoError(t, err)

		updated, err := taskRepo.FindByID(ctx, "task-1")
		require.NoError(t, err)
		require.NotNil(t, updated.LastCompletedAt)
		require.Equal(t, 2, updated.StreakCount)

		mockScheduler.AssertExpectations(t)
	})

	t.Run("ListTasksByUser", func(t *testing.T) {
		// Setup mock to return tasks
		taskRepo.On("FindByUser", ctx, int64(1)).Return([]*entity.Task{
			{ID: "task-1", Title: "Test Task 1", ChatID: 1},
			{ID: "task-2", Title: "Test Task 2", ChatID: 1},
		}, nil).Once()

		tasks, err := service.ListTasksByUser(ctx, 1)
		require.NoError(t, err)
		require.Len(t, tasks, 2)
	})
}
