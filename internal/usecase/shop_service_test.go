package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
	"github.com/supercakecrumb/adhd-game-bot/internal/ports"
	"github.com/supercakecrumb/adhd-game-bot/internal/usecase"
)

// Mock repositories
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

func (m *MockUserRepository) FindByChatID(ctx context.Context, chatID int64) ([]*entity.User, error) {
	args := m.Called(ctx, chatID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockUserRepository) UpdateBalance(ctx context.Context, userID int64, delta valueobject.Decimal) error {
	args := m.Called(ctx, userID, delta)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
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

type MockTxManager struct {
	mock.Mock
}

func (m *MockTxManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	// Execute the function directly (no actual transaction)
	if args.Get(0) == nil {
		return fn(ctx)
	}
	return args.Error(0)
}

func TestShopService_SetCurrencyName(t *testing.T) {
	ctx := context.Background()

	userRepo := new(MockUserRepository)
	chatConfigRepo := new(MockChatConfigRepository)
	shopItemRepo := new(MockShopItemRepository)
	purchaseRepo := new(MockPurchaseRepository)
	uuidGen := new(MockUUIDGenerator)
	txManager := new(MockTxManager)

	service := usecase.NewShopService(shopItemRepo, purchaseRepo, userRepo, chatConfigRepo, uuidGen, txManager)

	t.Run("Create new config", func(t *testing.T) {
		chatConfigRepo.On("FindByChatID", ctx, int64(100)).Return(nil, ports.ErrChatConfigNotFound).Once()
		chatConfigRepo.On("Create", ctx, mock.AnythingOfType("*entity.ChatConfig")).Return(nil).Once()

		err := service.SetCurrencyName(ctx, 100, "Gold Coins")
		assert.NoError(t, err)

		chatConfigRepo.AssertExpectations(t)
	})

	t.Run("Update existing config", func(t *testing.T) {
		existingConfig := &entity.ChatConfig{
			ChatID:       100,
			CurrencyName: "Points",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		chatConfigRepo.On("FindByChatID", ctx, int64(100)).Return(existingConfig, nil).Once()
		chatConfigRepo.On("Update", ctx, mock.AnythingOfType("*entity.ChatConfig")).Return(nil).Once()

		err := service.SetCurrencyName(ctx, 100, "Gold Coins")
		assert.NoError(t, err)

		chatConfigRepo.AssertExpectations(t)
	})
}

func TestShopService_PurchaseItem(t *testing.T) {
	ctx := context.Background()

	userRepo := new(MockUserRepository)
	chatConfigRepo := new(MockChatConfigRepository)
	shopItemRepo := new(MockShopItemRepository)
	purchaseRepo := new(MockPurchaseRepository)
	uuidGen := new(MockUUIDGenerator)
	txManager := new(MockTxManager)

	service := usecase.NewShopService(shopItemRepo, purchaseRepo, userRepo, chatConfigRepo, uuidGen, txManager)

	user := &entity.User{
		ID:      1,
		ChatID:  100,
		Balance: valueobject.NewDecimal("100"),
	}

	t.Run("Successful purchase", func(t *testing.T) {
		item := &entity.ShopItem{
			ID:       1,
			ChatID:   100,
			Code:     "BOOST",
			Name:     "XP Boost",
			Price:    valueobject.NewDecimal("50"),
			Stock:    intPtr(10),
			IsActive: true,
		}

		txManager.On("WithTx", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil).Once()
		userRepo.On("FindByID", ctx, int64(1)).Return(user, nil).Once()
		shopItemRepo.On("FindByCode", ctx, int64(100), "BOOST").Return(item, nil).Once()
		userRepo.On("UpdateBalance", ctx, int64(1), valueobject.NewDecimal("-50")).Return(nil).Once()
		shopItemRepo.On("Update", ctx, mock.AnythingOfType("*entity.ShopItem")).Return(nil).Once()
		purchaseRepo.On("Create", ctx, mock.AnythingOfType("*entity.Purchase")).Return(nil).Once()

		purchase, err := service.PurchaseItem(ctx, 1, "BOOST", 1)
		assert.NoError(t, err)
		assert.NotNil(t, purchase)
		assert.Equal(t, int64(1), purchase.UserID)
		assert.Equal(t, int64(1), purchase.ItemID)
		assert.Equal(t, "50", purchase.TotalCost.String())

		userRepo.AssertExpectations(t)
		shopItemRepo.AssertExpectations(t)
		purchaseRepo.AssertExpectations(t)
	})

	t.Run("Insufficient funds", func(t *testing.T) {
		poorUser := &entity.User{
			ID:      2,
			ChatID:  100,
			Balance: valueobject.NewDecimal("10"),
		}

		item := &entity.ShopItem{
			ID:       1,
			ChatID:   100,
			Code:     "BOOST",
			Name:     "XP Boost",
			Price:    valueobject.NewDecimal("50"),
			Stock:    intPtr(10),
			IsActive: true,
		}

		txManager.On("WithTx", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil).Once()
		userRepo.On("FindByID", ctx, int64(2)).Return(poorUser, nil).Once()
		shopItemRepo.On("FindByCode", ctx, int64(100), "BOOST").Return(item, nil).Once()

		_, err := service.PurchaseItem(ctx, 2, "BOOST", 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insufficient balance")

		userRepo.AssertExpectations(t)
		shopItemRepo.AssertExpectations(t)
	})

	t.Run("Out of stock", func(t *testing.T) {
		item := &entity.ShopItem{
			ID:       1,
			ChatID:   100,
			Code:     "LIMITED",
			Name:     "Limited Item",
			Price:    valueobject.NewDecimal("50"),
			Stock:    intPtr(0),
			IsActive: true,
		}

		txManager.On("WithTx", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil).Once()
		userRepo.On("FindByID", ctx, int64(1)).Return(user, nil).Once()
		shopItemRepo.On("FindByCode", ctx, int64(100), "LIMITED").Return(item, nil).Once()

		_, err := service.PurchaseItem(ctx, 1, "LIMITED", 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insufficient stock")

		userRepo.AssertExpectations(t)
		shopItemRepo.AssertExpectations(t)
	})
}

func intPtr(i int) *int {
	return &i
}
