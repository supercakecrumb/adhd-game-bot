package usecase

import (
	"context"
	"fmt"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
	"github.com/supercakecrumb/adhd-game-bot/internal/ports"
)

type ShopService struct {
	shopItemRepo   ports.ShopItemRepository
	purchaseRepo   ports.PurchaseRepository
	userRepo       ports.UserRepository
	chatConfigRepo ports.ChatConfigRepository
	uuidGen        ports.UUIDGenerator
}

func NewShopService(
	shopItemRepo ports.ShopItemRepository,
	purchaseRepo ports.PurchaseRepository,
	userRepo ports.UserRepository,
	chatConfigRepo ports.ChatConfigRepository,
	uuidGen ports.UUIDGenerator,
) *ShopService {
	return &ShopService{
		shopItemRepo:   shopItemRepo,
		purchaseRepo:   purchaseRepo,
		userRepo:       userRepo,
		chatConfigRepo: chatConfigRepo,
		uuidGen:        uuidGen,
	}
}

// CreateShopItem creates a new item in the shop
func (s *ShopService) CreateShopItem(ctx context.Context, item *entity.ShopItem) error {
	return s.shopItemRepo.Create(ctx, item)
}

// GetShopItems returns all available items for a chat (including global items)
func (s *ShopService) GetShopItems(ctx context.Context, chatID int64) ([]*entity.ShopItem, error) {
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

// PurchaseItem handles a user purchasing an item from the shop
func (s *ShopService) PurchaseItem(ctx context.Context, userID int64, itemCode string, quantity int) (*entity.Purchase, error) {
	// Get user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Get item
	item, err := s.shopItemRepo.FindByCode(ctx, user.ChatID, itemCode)
	if err != nil {
		// Try global items
		item, err = s.shopItemRepo.FindByCode(ctx, 0, itemCode)
		if err != nil {
			return nil, fmt.Errorf("item not found: %w", err)
		}
	}

	// Check if item is active
	if !item.IsActive {
		return nil, fmt.Errorf("item is not available")
	}

	// Check stock
	if item.Stock != nil && *item.Stock < quantity {
		return nil, fmt.Errorf("insufficient stock")
	}

	// Calculate total cost
	totalCost := item.Price.Mul(valueobject.NewDecimal(fmt.Sprintf("%d", quantity)))

	// Check user balance
	if user.Balance.Cmp(totalCost) < 0 {
		return nil, fmt.Errorf("insufficient balance")
	}

	// Create purchase record
	purchase := &entity.Purchase{
		UserID:    userID,
		ItemID:    item.ID,
		ItemName:  item.Name,
		ItemPrice: item.Price,
		Quantity:  quantity,
		TotalCost: totalCost,
		Status:    "completed",
	}

	// In a real implementation, this would be done in a transaction
	// 1. Deduct user balance
	negativeAmount := totalCost.Mul(valueobject.NewDecimal("-1"))
	err = s.userRepo.UpdateBalance(ctx, userID, negativeAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	// 2. Update stock if limited
	if item.Stock != nil {
		newStock := *item.Stock - quantity
		item.Stock = &newStock
		err = s.shopItemRepo.Update(ctx, item)
		if err != nil {
			// Rollback balance update
			s.userRepo.UpdateBalance(ctx, userID, totalCost)
			return nil, fmt.Errorf("failed to update stock: %w", err)
		}
	}

	// 3. Create purchase record
	err = s.purchaseRepo.Create(ctx, purchase)
	if err != nil {
		// Rollback previous changes
		s.userRepo.UpdateBalance(ctx, userID, totalCost)
		if item.Stock != nil {
			newStock := *item.Stock + quantity
			item.Stock = &newStock
			s.shopItemRepo.Update(ctx, item)
		}
		return nil, fmt.Errorf("failed to create purchase: %w", err)
	}

	return purchase, nil
}

// GetUserPurchases returns all purchases for a user
func (s *ShopService) GetUserPurchases(ctx context.Context, userID int64) ([]*entity.Purchase, error) {
	return s.purchaseRepo.FindByUserID(ctx, userID)
}

// GetCurrencyName returns the currency name for a chat
func (s *ShopService) GetCurrencyName(ctx context.Context, chatID int64) (string, error) {
	config, err := s.chatConfigRepo.FindByChatID(ctx, chatID)
	if err != nil {
		// Return default if not configured
		return "Points", nil
	}
	return config.CurrencyName, nil
}

// SetCurrencyName sets the currency name for a chat
func (s *ShopService) SetCurrencyName(ctx context.Context, chatID int64, currencyName string) error {
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
