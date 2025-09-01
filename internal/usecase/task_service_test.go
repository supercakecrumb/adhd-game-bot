package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
	"github.com/supercakecrumb/adhd-game-bot/internal/usecase"
)

type mockTaskRepo struct {
	tasks map[string]*entity.Task
}

func (m *mockTaskRepo) Create(ctx context.Context, task *entity.Task) error {
	m.tasks[task.ID] = task
	return nil
}

func (m *mockTaskRepo) FindByID(ctx context.Context, id string) (*entity.Task, error) {
	task, exists := m.tasks[id]
	if !exists {
		return nil, entity.ErrTaskNotFound
	}
	return task, nil
}

func (m *mockTaskRepo) Update(ctx context.Context, task *entity.Task) error {
	if _, exists := m.tasks[task.ID]; !exists {
		return entity.ErrTaskNotFound
	}
	m.tasks[task.ID] = task
	return nil
}

func (m *mockTaskRepo) Delete(ctx context.Context, id string) error {
	delete(m.tasks, id)
	return nil
}

func (m *mockTaskRepo) FindByUser(ctx context.Context, userID int64) ([]*entity.Task, error) {
	var result []*entity.Task
	for _, task := range m.tasks {
		// In real implementation we would filter by userID
		result = append(result, task)
	}
	return result, nil
}

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

func (m *mockUUIDGen) New() string {
	return "generated-uuid"
}

func TestTaskService(t *testing.T) {
	ctx := context.Background()
	taskRepo := &mockTaskRepo{tasks: make(map[string]*entity.Task)}
	userRepo := &mockUserRepo{users: map[int64]*entity.User{1: {ID: 1}}}
	uuidGen := &mockUUIDGen{}

	mockScheduler := new(mockScheduler)
	service := usecase.NewTaskService(taskRepo, userRepo, uuidGen, mockScheduler)

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
	})

	t.Run("CompleteTask", func(t *testing.T) {
		task := &entity.Task{
			ID:          "task-1",
			Title:       "Complete Me",
			StreakCount: 0,
			Category:    "daily",
		}
		taskRepo.tasks[task.ID] = task

		// Mock the scheduler call for daily tasks
		mockScheduler.On("ScheduleRecurringTask", ctx, mock.AnythingOfType("*entity.Task")).Return(nil).Once()

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
