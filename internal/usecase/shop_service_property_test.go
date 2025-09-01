package usecase_test

import (
	"context"
	"fmt"
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
	"github.com/supercakecrumb/adhd-game-bot/internal/ports"
	"github.com/supercakecrumb/adhd-game-bot/internal/usecase"
)

// Property: User balance should never go negative after a purchase
func TestShopService_BalanceNeverNegative(t *testing.T) {
	f := func(balance, price uint64, stock uint16) bool {
		if balance > 1000000 || price > 1000000 {
			return true // Skip unrealistic values
		}

		ctx := context.Background()

		// Setup mocks
		userRepo := new(MockUserRepository)
		chatConfigRepo := new(MockChatConfigRepository)
		shopItemRepo := new(MockShopItemRepository)
		purchaseRepo := new(MockPurchaseRepository)
		uuidGen := new(MockUUIDGenerator)
		txManager := new(MockTxManager)

		service := usecase.NewShopService(
			shopItemRepo, purchaseRepo, userRepo,
			chatConfigRepo, uuidGen, txManager,
		)

		// Create test data
		user := &entity.User{
			ID:      1,
			ChatID:  100,
			Balance: valueobject.NewDecimal(fmt.Sprintf("%d", balance)),
		}

		stockInt := int(stock)
		item := &entity.ShopItem{
			ID:       1,
			ChatID:   100,
			Code:     "TEST",
			Name:     "Test Item",
			Price:    valueobject.NewDecimal(fmt.Sprintf("%d", price)),
			IsActive: true,
			Stock:    &stockInt,
		}

		// Setup mock expectations
		txManager.On("WithTx", ctx, mock.AnythingOfType("func(context.Context) error")).
			Run(func(args mock.Arguments) {
				fn := args.Get(1).(func(context.Context) error)

				userRepo.On("FindByID", ctx, int64(1)).Return(user, nil).Once()
				shopItemRepo.On("FindByCode", ctx, int64(100), "TEST").Return(item, nil).Once()

				if balance >= price && stock > 0 {
					userRepo.On("UpdateBalance", ctx, int64(1), mock.AnythingOfType("valueobject.Decimal")).Return(nil).Once()
					shopItemRepo.On("Update", ctx, mock.AnythingOfType("*entity.ShopItem")).Return(nil).Once()
					purchaseRepo.On("Create", ctx, mock.AnythingOfType("*entity.Purchase")).Return(nil).Once()
				}

				fn(ctx)
			}).Return(nil).Once()

		// Execute purchase
		_, err := service.PurchaseItem(ctx, 1, "TEST", 1)

		// Property check: if purchase succeeded, balance should not be negative
		if err == nil {
			// The balance after purchase would be: balance - price
			newBalance := int64(balance) - int64(price)
			return newBalance >= 0
		}

		// If purchase failed due to insufficient funds, that's expected
		return err == ports.ErrInsufficientFunds || err == ports.ErrInsufficientStock
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

// Property: Stock should never go negative after a purchase
func TestShopService_StockNeverNegative(t *testing.T) {
	f := func(stock uint16, quantity uint8) bool {
		if quantity == 0 || quantity > 100 {
			return true // Skip invalid quantities
		}

		ctx := context.Background()

		// Setup mocks
		userRepo := new(MockUserRepository)
		chatConfigRepo := new(MockChatConfigRepository)
		shopItemRepo := new(MockShopItemRepository)
		purchaseRepo := new(MockPurchaseRepository)
		uuidGen := new(MockUUIDGenerator)
		txManager := new(MockTxManager)

		service := usecase.NewShopService(
			shopItemRepo, purchaseRepo, userRepo,
			chatConfigRepo, uuidGen, txManager,
		)

		// Create test data
		user := &entity.User{
			ID:      1,
			ChatID:  100,
			Balance: valueobject.NewDecimal("1000000"), // Plenty of money
		}

		stockInt := int(stock)
		item := &entity.ShopItem{
			ID:       1,
			ChatID:   100,
			Code:     "TEST",
			Name:     "Test Item",
			Price:    valueobject.NewDecimal("1"),
			IsActive: true,
			Stock:    &stockInt,
		}

		// Setup mock expectations
		txManager.On("WithTx", ctx, mock.AnythingOfType("func(context.Context) error")).
			Run(func(args mock.Arguments) {
				fn := args.Get(1).(func(context.Context) error)

				userRepo.On("FindByID", ctx, int64(1)).Return(user, nil).Once()
				shopItemRepo.On("FindByCode", ctx, int64(100), "TEST").Return(item, nil).Once()

				if stock >= uint16(quantity) {
					userRepo.On("UpdateBalance", ctx, int64(1), mock.AnythingOfType("valueobject.Decimal")).Return(nil).Once()
					shopItemRepo.On("Update", ctx, mock.AnythingOfType("*entity.ShopItem")).Return(nil).Once()
					purchaseRepo.On("Create", ctx, mock.AnythingOfType("*entity.Purchase")).Return(nil).Once()
				}

				fn(ctx)
			}).Return(nil).Once()

		// Execute purchase
		_, err := service.PurchaseItem(ctx, 1, "TEST", int(quantity))

		// Property check: if purchase succeeded, stock should not be negative
		if err == nil {
			// The stock after purchase would be: stock - quantity
			newStock := int(stock) - int(quantity)
			return newStock >= 0
		}

		// If purchase failed due to insufficient stock, that's expected
		return err == ports.ErrInsufficientStock
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

// Property: Total cost calculation should be correct
func TestShopService_TotalCostCalculation(t *testing.T) {
	f := func(price uint32, quantity uint8) bool {
		if price > 1000000 || quantity == 0 || quantity > 100 {
			return true // Skip unrealistic values
		}

		ctx := context.Background()

		// Setup mocks
		userRepo := new(MockUserRepository)
		chatConfigRepo := new(MockChatConfigRepository)
		shopItemRepo := new(MockShopItemRepository)
		purchaseRepo := new(MockPurchaseRepository)
		uuidGen := new(MockUUIDGenerator)
		txManager := new(MockTxManager)

		service := usecase.NewShopService(
			shopItemRepo, purchaseRepo, userRepo,
			chatConfigRepo, uuidGen, txManager,
		)

		// Calculate expected total
		expectedTotal := uint64(price) * uint64(quantity)

		// Create test data
		user := &entity.User{
			ID:      1,
			ChatID:  100,
			Balance: valueobject.NewDecimal(fmt.Sprintf("%d", expectedTotal+1000)), // Enough balance
		}

		item := &entity.ShopItem{
			ID:       1,
			ChatID:   100,
			Code:     "TEST",
			Name:     "Test Item",
			Price:    valueobject.NewDecimal(fmt.Sprintf("%d", price)),
			IsActive: true,
			Stock:    nil, // Unlimited stock
		}

		var capturedPurchase *entity.Purchase

		// Setup mock expectations
		txManager.On("WithTx", ctx, mock.AnythingOfType("func(context.Context) error")).
			Run(func(args mock.Arguments) {
				fn := args.Get(1).(func(context.Context) error)

				userRepo.On("FindByID", ctx, int64(1)).Return(user, nil).Once()
				shopItemRepo.On("FindByCode", ctx, int64(100), "TEST").Return(item, nil).Once()
				userRepo.On("UpdateBalance", ctx, int64(1), mock.AnythingOfType("valueobject.Decimal")).Return(nil).Once()
				shopItemRepo.On("Update", ctx, mock.AnythingOfType("*entity.ShopItem")).Return(nil).Once()
				purchaseRepo.On("Create", ctx, mock.AnythingOfType("*entity.Purchase")).
					Run(func(args mock.Arguments) {
						capturedPurchase = args.Get(1).(*entity.Purchase)
					}).Return(nil).Once()

				fn(ctx)
			}).Return(nil).Once()

		// Execute purchase
		_, err := service.PurchaseItem(ctx, 1, "TEST", int(quantity))
		require.NoError(t, err)

		// Property check: total cost should equal price * quantity
		expectedTotalDecimal := valueobject.NewDecimal(fmt.Sprintf("%d", expectedTotal))
		return capturedPurchase != nil &&
			capturedPurchase.TotalCost.Cmp(expectedTotalDecimal) == 0 &&
			capturedPurchase.Quantity == int(quantity)
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

// Edge case tests for shop service
// Property: Purchases with same idempotency key should only deduct balance once
func TestShopService_IdempotentPurchase(t *testing.T) {
	f := func(price uint32) bool {
		if price > 1000000 {
			return true // Skip unrealistic values
		}

		ctx := context.Background()

		// Setup mocks
		userRepo := new(MockUserRepository)
		chatConfigRepo := new(MockChatConfigRepository)
		shopItemRepo := new(MockShopItemRepository)
		purchaseRepo := new(MockPurchaseRepository)
		uuidGen := new(MockUUIDGenerator)
		txManager := new(MockTxManager)

		service := usecase.NewShopService(
			shopItemRepo, purchaseRepo, userRepo,
			chatConfigRepo, uuidGen, txManager,
		)

		// Create test data
		user := &entity.User{
			ID:      1,
			ChatID:  100,
			Balance: valueobject.NewDecimal(fmt.Sprintf("%d", price*2)), // Enough for 2 purchases
		}

		item := &entity.ShopItem{
			ID:       1,
			ChatID:   100,
			Code:     "TEST",
			Name:     "Test Item",
			Price:    valueobject.NewDecimal(fmt.Sprintf("%d", price)),
			IsActive: true,
			Stock:    nil, // Unlimited stock
		}

		// Setup mock expectations
		txManager.On("WithTx", ctx, mock.AnythingOfType("func(context.Context) error")).
			Run(func(args mock.Arguments) {
				fn := args.Get(1).(func(context.Context) error)

				userRepo.On("FindByID", ctx, int64(1)).Return(user, nil).Once()
				shopItemRepo.On("FindByCode", ctx, int64(100), "TEST").Return(item, nil).Once()
				userRepo.On("UpdateBalance", ctx, int64(1), mock.AnythingOfType("valueobject.Decimal")).Return(nil).Once()
				shopItemRepo.On("Update", ctx, mock.AnythingOfType("*entity.ShopItem")).Return(nil).Once()
				purchaseRepo.On("Create", ctx, mock.AnythingOfType("*entity.Purchase")).Return(nil).Once()

				fn(ctx)
			}).Return(nil).Twice() // Expect two calls but same mocks

		// Execute first purchase
		_, err := service.PurchaseItem(ctx, 1, "TEST", 1)
		if err != nil {
			return false
		}

		// Execute second identical purchase
		_, err = service.PurchaseItem(ctx, 1, "TEST", 1)
		if err != nil {
			return false
		}

		// Verify balance was only deducted once
		finalBalance := user.Balance.String()
		expectedBalance := fmt.Sprintf("%d", price)
		return finalBalance == expectedBalance
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestShopService_EdgeCases(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "Purchase with exact balance",
			test: func(t *testing.T) {
				ctx := context.Background()

				userRepo := new(MockUserRepository)
				chatConfigRepo := new(MockChatConfigRepository)
				shopItemRepo := new(MockShopItemRepository)
				purchaseRepo := new(MockPurchaseRepository)
				uuidGen := new(MockUUIDGenerator)
				txManager := new(MockTxManager)

				service := usecase.NewShopService(
					shopItemRepo, purchaseRepo, userRepo,
					chatConfigRepo, uuidGen, txManager,
				)

				user := &entity.User{
					ID:      1,
					ChatID:  100,
					Balance: valueobject.NewDecimal("100"),
				}

				item := &entity.ShopItem{
					ID:       1,
					ChatID:   100,
					Code:     "TEST",
					Name:     "Test Item",
					Price:    valueobject.NewDecimal("100"),
					IsActive: true,
					Stock:    nil,
				}

				txManager.On("WithTx", ctx, mock.AnythingOfType("func(context.Context) error")).
					Run(func(args mock.Arguments) {
						fn := args.Get(1).(func(context.Context) error)

						userRepo.On("FindByID", ctx, int64(1)).Return(user, nil).Once()
						shopItemRepo.On("FindByCode", ctx, int64(100), "TEST").Return(item, nil).Once()
						userRepo.On("UpdateBalance", ctx, int64(1), mock.AnythingOfType("valueobject.Decimal")).Return(nil).Once()
						shopItemRepo.On("Update", ctx, mock.AnythingOfType("*entity.ShopItem")).Return(nil).Once()
						purchaseRepo.On("Create", ctx, mock.AnythingOfType("*entity.Purchase")).Return(nil).Once()

						fn(ctx)
					}).Return(nil).Once()

				_, err := service.PurchaseItem(ctx, 1, "TEST", 1)
				require.NoError(t, err)
			},
		},
		{
			name: "Purchase with zero quantity",
			test: func(t *testing.T) {
				ctx := context.Background()

				userRepo := new(MockUserRepository)
				chatConfigRepo := new(MockChatConfigRepository)
				shopItemRepo := new(MockShopItemRepository)
				purchaseRepo := new(MockPurchaseRepository)
				uuidGen := new(MockUUIDGenerator)
				txManager := new(MockTxManager)

				service := usecase.NewShopService(
					shopItemRepo, purchaseRepo, userRepo,
					chatConfigRepo, uuidGen, txManager,
				)

				// Should fail validation before any repository calls
				_, err := service.PurchaseItem(ctx, 1, "TEST", 0)
				require.Error(t, err)
			},
		},
		{
			name: "Purchase inactive item",
			test: func(t *testing.T) {
				ctx := context.Background()

				userRepo := new(MockUserRepository)
				chatConfigRepo := new(MockChatConfigRepository)
				shopItemRepo := new(MockShopItemRepository)
				purchaseRepo := new(MockPurchaseRepository)
				uuidGen := new(MockUUIDGenerator)
				txManager := new(MockTxManager)

				service := usecase.NewShopService(
					shopItemRepo, purchaseRepo, userRepo,
					chatConfigRepo, uuidGen, txManager,
				)

				user := &entity.User{
					ID:      1,
					ChatID:  100,
					Balance: valueobject.NewDecimal("1000"),
				}

				item := &entity.ShopItem{
					ID:       1,
					ChatID:   100,
					Code:     "TEST",
					Name:     "Test Item",
					Price:    valueobject.NewDecimal("100"),
					IsActive: false, // Inactive
					Stock:    nil,
				}

				txManager.On("WithTx", ctx, mock.AnythingOfType("func(context.Context) error")).
					Run(func(args mock.Arguments) {
						fn := args.Get(1).(func(context.Context) error)

						userRepo.On("FindByID", ctx, int64(1)).Return(user, nil).Once()
						shopItemRepo.On("FindByCode", ctx, int64(100), "TEST").Return(item, nil).Once()

						fn(ctx)
					}).Return(nil).Once()

				_, err := service.PurchaseItem(ctx, 1, "TEST", 1)
				require.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
