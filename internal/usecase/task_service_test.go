package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
	"github.com/supercakecrumb/adhd-game-bot/internal/ports"
	"github.com/supercakecrumb/adhd-game-bot/internal/usecase"
	"github.com/supercakecrumb/adhd-game-bot/internal/usecase/testhelpers"
)

type mockUserRepo struct {
	users map[int64]*entity.User
}

func (m *mockUserRepo) Create(ctx context.Context, user *entity.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, entity.ErrUserNotFound
	}
	return user, nil
}

func (m *mockUserRepo) FindByChatID(ctx context.Context, chatID int64) ([]*entity.User, error) {
	var users []*entity.User
	for _, user := range m.users {
		if user.ChatID == chatID {
			users = append(users, user)
		}
	}
	return users, nil
}

func (m *mockUserRepo) UpdateBalance(ctx context.Context, userID int64, delta valueobject.Decimal) error {
	// Not needed for these tests
	return nil
}

func (m *mockUserRepo) Delete(ctx context.Context, id int64) error {
	delete(m.users, id)
	return nil
}

type mockUUIDGen struct{}

type mockScheduler struct {
	mock.Mock
}

func (m *mockScheduler) ScheduleRecurringTask(ctx context.Context, task *entity.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *mockScheduler) CancelScheduledTask(ctx context.Context, taskID string) error {
	args := m.Called(ctx, taskID)
	return args.Error(0)
}

func (m *mockScheduler) GetNextOccurrence(ctx context.Context, taskID string) (time.Time, error) {
	args := m.Called(ctx, taskID)
	return args.Get(0).(time.Time), args.Error(1)
}

type mockIdempotencyRepo struct {
	mock.Mock
}

func (m *mockIdempotencyRepo) Create(ctx context.Context, key *entity.IdempotencyKey) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *mockIdempotencyRepo) FindByKey(ctx context.Context, key string) (*entity.IdempotencyKey, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.IdempotencyKey), args.Error(1)
}

func (m *mockIdempotencyRepo) Update(ctx context.Context, key *entity.IdempotencyKey) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *mockIdempotencyRepo) DeleteExpired(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type mockTxManager struct {
	mock.Mock
}

func (m *mockTxManager) WithTx(ctx context.Context, fn func(context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

func (m *mockUUIDGen) New() string {
	return "generated-uuid"
}

func TestTaskService(t *testing.T) {
	ctx := context.Background()
	taskRepo := new(testhelpers.MockTaskRepository)
	userRepo := &mockUserRepo{users: map[int64]*entity.User{1: {ID: 1}}}

	// Setup taskRepo expectations
	taskRepo.On("Create", ctx, mock.AnythingOfType("*entity.Task")).Return(nil)
	taskRepo.On("FindByID", ctx, mock.AnythingOfType("string")).Return(nil, entity.ErrTaskNotFound)
	taskRepo.On("FindByUser", ctx, int64(1)).Return([]*entity.Task{}, nil)
	uuidGen := &mockUUIDGen{}
	mockScheduler := new(mockScheduler)
	mockIdempotencyRepo := new(testhelpers.MockIdempotencyRepository)
	mockTxManager := new(mockTxManager)

	service := usecase.NewTaskService(taskRepo, userRepo, uuidGen, mockScheduler, mockIdempotencyRepo, mockTxManager)

	t.Run("CreateTask", func(t *testing.T) {
		task := &entity.Task{
			Title:       "Test Task",
			Description: "Test Description",
			Category:    "daily",
		}

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

		// Mock idempotency check
		mockIdempotencyRepo.On("FindByKey", ctx, mock.AnythingOfType("string")).Return(nil, ports.ErrIdempotencyKeyNotFound).Once()
		mockIdempotencyRepo.On("Create", ctx, mock.AnythingOfType("*entity.IdempotencyKey")).Return(nil).Once()
		mockIdempotencyRepo.On("Update", ctx, mock.AnythingOfType("*entity.IdempotencyKey")).Return(nil).Once()

		// Mock transaction
		mockTxManager.On("WithTx", ctx, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			// Mock the scheduler call for daily tasks
			mockScheduler.On("ScheduleRecurringTask", ctx, mock.AnythingOfType("*entity.Task")).Return(nil).Once()
			fn(ctx)
		}).Return(nil).Once()

		err := service.CompleteTask(ctx, 1, "task-1")
		require.NoError(t, err)

		updated, err := taskRepo.FindByID(ctx, "task-1")
		require.NoError(t, err)
		require.NotNil(t, updated.LastCompletedAt)
		require.Equal(t, 1, updated.StreakCount)

		mockScheduler.AssertExpectations(t)
	})

	t.Run("ListTasksByUser", func(t *testing.T) {
		tasks, err := service.ListTasksByUser(ctx, 1)
		require.NoError(t, err)
		require.Len(t, tasks, 2) // Created in previous tests
	})
}
