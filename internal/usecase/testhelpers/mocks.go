package testhelpers

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) UpdateBalance(ctx context.Context, userID int64, delta valueobject.Decimal) error {
	args := m.Called(ctx, userID, delta)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) FindByChatID(ctx context.Context, chatID int64) ([]*entity.User, error) {
	args := m.Called(ctx, chatID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.User), args.Error(1)
}

type MockChatConfigRepository struct {
	mock.Mock
}

func (m *MockChatConfigRepository) Create(ctx context.Context, config *entity.ChatConfig) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

func (m *MockChatConfigRepository) FindByChatID(ctx context.Context, chatID int64) (*entity.ChatConfig, error) {
	args := m.Called(ctx, chatID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ChatConfig), args.Error(1)
}

func (m *MockChatConfigRepository) Update(ctx context.Context, config *entity.ChatConfig) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

type MockShopItemRepository struct {
	mock.Mock
}

func (m *MockShopItemRepository) Create(ctx context.Context, item *entity.ShopItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockShopItemRepository) FindByID(ctx context.Context, id int64) (*entity.ShopItem, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ShopItem), args.Error(1)
}

func (m *MockShopItemRepository) FindByCode(ctx context.Context, chatID int64, code string) (*entity.ShopItem, error) {
	args := m.Called(ctx, chatID, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ShopItem), args.Error(1)
}

func (m *MockShopItemRepository) FindByChatID(ctx context.Context, chatID int64) ([]*entity.ShopItem, error) {
	args := m.Called(ctx, chatID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.ShopItem), args.Error(1)
}

func (m *MockShopItemRepository) Update(ctx context.Context, item *entity.ShopItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockShopItemRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockPurchaseRepository struct {
	mock.Mock
}

func (m *MockPurchaseRepository) Create(ctx context.Context, purchase *entity.Purchase) error {
	args := m.Called(ctx, purchase)
	return args.Error(0)
}

func (m *MockPurchaseRepository) FindByID(ctx context.Context, id int64) (*entity.Purchase, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Purchase), args.Error(1)
}

func (m *MockPurchaseRepository) FindByUserID(ctx context.Context, userID int64) ([]*entity.Purchase, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Purchase), args.Error(1)
}

func (m *MockPurchaseRepository) FindByItemID(ctx context.Context, itemID int64) ([]*entity.Purchase, error) {
	args := m.Called(ctx, itemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Purchase), args.Error(1)
}

type MockUUIDGenerator struct {
	mock.Mock
}

func (m *MockUUIDGenerator) New() string {
	args := m.Called()
	return args.String(0)
}

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(ctx context.Context, task *entity.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) FindByID(ctx context.Context, id string) (*entity.Task, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Task), args.Error(1)
}

func (m *MockTaskRepository) FindActiveByUser(ctx context.Context, userID int64) ([]*entity.Task, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Task), args.Error(1)
}

func (m *MockTaskRepository) FindWithTimers(ctx context.Context, userID int64) ([]*entity.Task, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Task), args.Error(1)
}

func (m *MockTaskRepository) Update(ctx context.Context, task *entity.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTaskRepository) FindByUser(ctx context.Context, userID int64) ([]*entity.Task, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Task), args.Error(1)
}

func (m *MockTaskRepository) FindWithSchedules(ctx context.Context, userID int64) ([]*entity.Task, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Task), args.Error(1)
}

func (m *MockTaskRepository) BulkUpdate(ctx context.Context, tasks []*entity.Task) error {
	args := m.Called(ctx, tasks)
	return args.Error(0)
}

type MockIdempotencyRepository struct {
	mock.Mock
}

func (m *MockIdempotencyRepository) Create(ctx context.Context, key *entity.IdempotencyKey) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockIdempotencyRepository) FindByKey(ctx context.Context, key string) (*entity.IdempotencyKey, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.IdempotencyKey), args.Error(1)
}

func (m *MockIdempotencyRepository) Update(ctx context.Context, key *entity.IdempotencyKey) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockIdempotencyRepository) DeleteExpired(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockIdempotencyRepository) Purge(ctx context.Context, olderThan time.Time) error {
	args := m.Called(ctx, olderThan)
	return args.Error(0)
}

type MockTxManager struct {
	mock.Mock
}

func (m *MockTxManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	if args.Get(0) == nil {
		return fn(ctx)
	}
	return args.Error(0)
}
