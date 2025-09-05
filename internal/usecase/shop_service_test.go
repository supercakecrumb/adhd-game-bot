package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
)

// Mock implementations for core test cases
type mockShopItemRepo struct{ mock.Mock }

func (m *mockShopItemRepo) FindByCode(ctx context.Context, chatID int64, code string) (*entity.ShopItem, error) {
	args := m.Called(ctx, chatID, code)
	return args.Get(0).(*entity.ShopItem), args.Error(1)
}

type mockPurchaseRepo struct{ mock.Mock }

func (m *mockPurchaseRepo) Create(ctx context.Context, p *entity.Purchase) error {
	return m.Called(ctx, p).Error(0)
}

type mockUserRepo struct{ mock.Mock }

func (m *mockUserRepo) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.User), args.Error(1)
}
func (m *mockUserRepo) UpdateBalance(ctx context.Context, id int64, amount valueobject.Decimal) error {
	return m.Called(ctx, id, amount).Error(0)
}

type mockIdempotencyRepo struct{ mock.Mock }

func (m *mockIdempotencyRepo) FindByKey(ctx context.Context, key string) (*entity.IdempotencyKey, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(*entity.IdempotencyKey), args.Error(1)
}
func (m *mockIdempotencyRepo) Create(ctx context.Context, key *entity.IdempotencyKey) error {
	return m.Called(ctx, key).Error(0)
}

func TestShopServiceV2(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	shopItemRepo := &mockShopItemRepo{}
	purchaseRepo := &mockPurchaseRepo{}
	userRepo := &mockUserRepo{}
	idempotencyRepo := &mockIdempotencyRepo{}

	service := NewShopServiceV2(
		shopItemRepo,
		purchaseRepo,
		userRepo,
		nil, // chatConfigRepo
		nil, // discountTierRepo
		nil, // uuidGen
		nil, // txManager
		idempotencyRepo,
	)

	t.Run("PurchaseItemWithIdempotency succeeds", func(t *testing.T) {
		user := &entity.User{ID: 1, Balance: valueobject.NewDecimal("100.00")}
		item := &entity.ShopItem{ID: 1, Price: valueobject.NewDecimal("10.00"), IsActive: true}

		userRepo.On("FindByID", ctx, int64(1)).Return(user, nil)
		shopItemRepo.On("FindByCode", ctx, int64(0), "ITEM").Return(item, nil)
		purchaseRepo.On("Create", ctx, mock.Anything).Return(nil)
		userRepo.On("UpdateBalance", ctx, int64(1), mock.Anything).Return(nil)
		idempotencyRepo.On("Create", ctx, mock.Anything).Return(nil)

		purchase, err := service.PurchaseItemWithIdempotency(ctx, 1, "ITEM", 1, "key123")
		assert.NoError(t, err)
		assert.NotNil(t, purchase)
	})
}
