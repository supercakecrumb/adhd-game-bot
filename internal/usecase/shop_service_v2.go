package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
	"github.com/supercakecrumb/adhd-game-bot/internal/ports"
)

type ShopServiceV2 struct {
	shopItemRepo     ports.ShopItemRepository
	purchaseRepo     ports.PurchaseRepository
	userRepo         ports.UserRepository
	chatConfigRepo   ports.ChatConfigRepository
	discountTierRepo ports.DiscountTierRepository
	uuidGen          ports.UUIDGenerator
	txManager        ports.TxManager
	idempotencyRepo  ports.IdempotencyRepository
}

func NewShopServiceV2(
	shopItemRepo ports.ShopItemRepository,
	purchaseRepo ports.PurchaseRepository,
	userRepo ports.UserRepository,
	chatConfigRepo ports.ChatConfigRepository,
	discountTierRepo ports.DiscountTierRepository,
	uuidGen ports.UUIDGenerator,
	txManager ports.TxManager,
	idempotencyRepo ports.IdempotencyRepository,
) *ShopServiceV2 {
	return &ShopServiceV2{
		shopItemRepo:     shopItemRepo,
		purchaseRepo:     purchaseRepo,
		userRepo:         userRepo,
		chatConfigRepo:   chatConfigRepo,
		discountTierRepo: discountTierRepo,
		uuidGen:          uuidGen,
		txManager:        txManager,
		idempotencyRepo:  idempotencyRepo,
	}
}

// PurchaseItemWithIdempotency handles purchases with idempotency key support
func (s *ShopServiceV2) PurchaseItemWithIdempotency(
	ctx context.Context,
	userID int64,
	itemCode string,
	quantity int,
	idempotencyKey string,
) (*entity.Purchase, error) {
	var purchase *entity.Purchase

	if idempotencyKey != "" {
		existing, err := s.idempotencyRepo.FindByKey(ctx, idempotencyKey)
		if err == nil && existing != nil {
			return nil, ports.ErrDuplicateRequest
		}
	}

	err := s.txManager.WithTx(ctx, func(txCtx context.Context) error {
		// Get user
		user, err := s.userRepo.FindByID(txCtx, userID)
		if err != nil {
			return fmt.Errorf("user not found: %w", err)
		}

		// Get item
		item, err := s.shopItemRepo.FindByCode(txCtx, user.ChatID, itemCode)
		if err != nil {
			item, err = s.shopItemRepo.FindByCode(txCtx, 0, itemCode)
			if err != nil {
				return fmt.Errorf("item not found: %w", err)
			}
		}

		// Check item availability
		if !item.IsActive {
			return fmt.Errorf("item is not available")
		}
		if item.Stock != nil && *item.Stock < quantity {
			return fmt.Errorf("insufficient stock")
		}

		// Calculate total cost
		totalCost := item.Price.Mul(valueobject.NewDecimal(fmt.Sprintf("%d", quantity)))
		if user.Balance.Cmp(totalCost) < 0 {
			return fmt.Errorf("insufficient balance")
		}

		// Handle reward tiers
		if item.DiscountTierID != nil {
			discountTier, err := s.discountTierRepo.FindByID(txCtx, *item.DiscountTierID)
			if err == nil && discountTier != nil {
				totalCost = totalCost.Mul(valueobject.NewDecimal(fmt.Sprintf("%.2f", 1-discountTier.DiscountPercent/100)))
			}
		}

		// Create purchase
		purchase = &entity.Purchase{
			UserID:         userID,
			ItemID:         item.ID,
			ItemName:       item.Name,
			ItemPrice:      item.Price,
			Quantity:       quantity,
			TotalCost:      totalCost,
			Status:         "completed",
			DiscountTierID: item.DiscountTierID,
		}

		// Process transaction
		negativeAmount := totalCost.Mul(valueobject.NewDecimal("-1"))
		if err := s.userRepo.UpdateBalance(txCtx, userID, negativeAmount); err != nil {
			return fmt.Errorf("failed to update balance: %w", err)
		}

		if item.Stock != nil {
			newStock := *item.Stock - quantity
			item.Stock = &newStock
			if err := s.shopItemRepo.Update(txCtx, item); err != nil {
				return fmt.Errorf("failed to update stock: %w", err)
			}
		}

		if err := s.purchaseRepo.Create(txCtx, purchase); err != nil {
			return fmt.Errorf("failed to create purchase: %w", err)
		}

		// Record idempotency key if provided
		if idempotencyKey != "" {
			key := &entity.IdempotencyKey{
				Key:       idempotencyKey,
				ExpiresAt: time.Now().Add(24 * time.Hour),
			}
			if err := s.idempotencyRepo.Create(txCtx, key); err != nil {
				return fmt.Errorf("failed to record idempotency key: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return purchase, nil
}

// GetShopItems returns all available items for a chat (including global items)
func (s *ShopServiceV2) GetShopItems(ctx context.Context, chatID int64) ([]*entity.ShopItem, error) {
	// Get chat-specific items
	chatItems, err := s.shopItemRepo.FindByChatID(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat items: %w", err)
	}

	// Get global items (chatID = 0)
	globalItems, err := s.shopItemRepo.FindByChatID(ctx, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get global items: %w", err)
	}

	// Combine and filter active items
	var allItems []*entity.ShopItem
	for _, item := range append(chatItems, globalItems...) {
		if item.IsActive {
			allItems = append(allItems, item)
		}
	}

	return allItems, nil
}

// GetUserPurchases returns all purchases for a user
func (s *ShopServiceV2) GetUserPurchases(ctx context.Context, userID int64) ([]*entity.Purchase, error) {
	return s.purchaseRepo.FindByUserID(ctx, userID)
}

// GetCurrencyName returns the currency name for a chat
func (s *ShopServiceV2) GetCurrencyName(ctx context.Context, chatID int64) (string, error) {
	config, err := s.chatConfigRepo.FindByChatID(ctx, chatID)
	if err != nil {
		// Return default if not configured
		return "Points", nil
	}
	return config.CurrencyName, nil
}

// SetCurrencyName sets the currency name for a chat
func (s *ShopServiceV2) SetCurrencyName(ctx context.Context, chatID int64, currencyName string) error {
	config, err := s.chatConfigRepo.FindByChatID(ctx, chatID)
	if err != nil {
		// Create new config
		config = &entity.ChatConfig{
			ChatID:       chatID,
			CurrencyName: currencyName,
		}
		return s.chatConfigRepo.Create(ctx, config)
	}

	// Update existing config
	config.CurrencyName = currencyName
	return s.chatConfigRepo.Update(ctx, config)
}
