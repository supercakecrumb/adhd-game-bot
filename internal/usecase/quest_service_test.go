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

func TestQuestService(t *testing.T) {
	ctx := context.Background()
	questRepo := new(testhelpers.MockQuestRepository)
	userRepo := new(testhelpers.MockUserRepository)

	uuidGen := &mockUUIDGen{}
	mockScheduler := new(testhelpers.MockScheduler)
	mockIdempotencyRepo := new(testhelpers.MockIdempotencyRepository)
	mockTxManager := new(testhelpers.MockTxManager)

	service := usecase.NewQuestService(questRepo, userRepo, uuidGen, mockScheduler, mockIdempotencyRepo, mockTxManager)

	t.Run("CreateQuest", func(t *testing.T) {
		input := usecase.CreateQuestInput{
			Title:       "Test Quest",
			Description: "Test Description",
			Category:    "daily",
		}

		// Setup mock expectations
		questRepo.On("Create", ctx, mock.AnythingOfType("*entity.Quest")).Return(nil).Once()
		userRepo.On("FindByID", ctx, int64(1)).Return(&entity.User{ID: 1}, nil).Once()

		// Mock the scheduler call for daily quests
		mockScheduler.On("ScheduleRecurringTask", ctx, mock.AnythingOfType("*entity.Quest")).Return(nil).Once()

		created, err := service.CreateQuest(ctx, 1, "dungeon-1", input)
		require.NoError(t, err)
		require.Equal(t, "generated-uuid", created.ID)
		require.Equal(t, "Test Quest", created.Title)

		mockScheduler.AssertExpectations(t)
		mockIdempotencyRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})

	t.Run("CompleteQuest", func(t *testing.T) {
		quest := &entity.Quest{
			ID:          "quest-1",
			Title:       "Complete Me",
			StreakCount: 0,
			Category:    "daily",
		}
		questRepo.On("GetByID", ctx, "quest-1").Return(quest, nil)
		questRepo.On("Update", ctx, mock.AnythingOfType("*entity.Quest")).Return(nil)
		userRepo.On("FindByID", ctx, int64(1)).Return(&entity.User{ID: 1}, nil)

		// Mock idempotency check
		mockIdempotencyRepo.On("FindByKey", ctx, mock.AnythingOfType("string")).Return(nil, ports.ErrIdempotencyKeyNotFound).Once()
		mockIdempotencyRepo.On("Create", ctx, mock.AnythingOfType("*entity.IdempotencyKey")).Return(nil).Once()
		mockIdempotencyRepo.On("Update", ctx, mock.AnythingOfType("*entity.IdempotencyKey")).Return(nil).Once()

		// Mock transaction
		mockTxManager.On("WithTx", ctx, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			// Mock the scheduler call for daily quests
			mockScheduler.On("ScheduleRecurringTask", ctx, mock.AnythingOfType("*entity.Quest")).Return(nil)
			fn(ctx)
		}).Return(nil).Once()

		input := usecase.CompleteQuestInput{
			IdempotencyKey: "key-1",
		}
		err := service.CompleteQuest(ctx, 1, "quest-1", input)
		require.NoError(t, err)

		mockScheduler.AssertExpectations(t)
	})

	t.Run("ListQuests", func(t *testing.T) {
		userRepo.On("FindByID", ctx, int64(1)).Return(&entity.User{ID: 1}, nil).Once()
		// Setup mock to return quests
		questRepo.On("ListByDungeon", ctx, "dungeon-1").Return([]*entity.Quest{
			{ID: "quest-1", Title: "Test Quest 1"},
			{ID: "quest-2", Title: "Test Quest 2"},
		}, nil).Once()

		quests, err := service.ListQuests(ctx, 1, "dungeon-1")
		require.NoError(t, err)
		require.Len(t, quests, 2)
	})
}
