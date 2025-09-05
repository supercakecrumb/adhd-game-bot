package usecase

import (
	"context"
	"fmt"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
	"github.com/supercakecrumb/adhd-game-bot/internal/ports"
)

type ShopService struct {
	shopItemRepo    ports.ShopItemRepository
	purchaseRepo    ports.PurchaseRepository
	userRepo        ports.UserRepository
	chatConfigRepo  ports.ChatConfigRepository
	rewardTierRepo  ports.RewardTierRepository
	uuidGen         ports.UUIDGenerator
	txManager       ports.TxManager
	idempotencyRepo ports.IdempotencyRepository
	// Deprecated fields for backward compatibility
	legacyMode bool
}

// NewShopService creates a new ShopService instance
// For backward compatibility, idempotencyRepo is optional (will use no-op implementation if nil)
// NewShopService creates a new ShopService instance with backward compatibility
func NewShopService(
	shopItemRepo ports.ShopItemRepository,
	purchaseRepo ports.PurchaseRepository,
	userRepo ports.UserRepository,
	chatConfigRepo ports.ChatConfigRepository,
	args ...interface{},
) *ShopService {
	var (
		rewardTierRepo  ports.RewardTierRepository
		uuidGen         ports.UUIDGenerator
		txManager       ports.TxManager
		idempotencyRepo ports.IdempotencyRepository
	)

	// Handle backward compatibility
	switch len(args) {
	case 3: // Old format (uuidGen, txManager)
		uuidGen = args[0].(ports.UUIDGenerator)
		txManager = args[1].(ports.TxManager)
		idempotencyRepo = &noopIdempotencyRepo{}
	case 4: // New format (rewardTierRepo, uuidGen, txManager, idempotencyRepo)
		rewardTierRepo = args[0].(ports.RewardTierRepository)
		uuidGen = args[1].(ports.UUIDGenerator)
		txManager = args[2].(ports.TxManager)
		idempotencyRepo = args[3].(ports.IdempotencyRepository)
	default:
		panic("invalid number of arguments")
	}
	svc := &ShopService{
		shopItemRepo:    shopItemRepo,
		purchaseRepo:    purchaseRepo,
		userRepo:        userRepo,
		chatConfigRepo:  chatConfigRepo,
		rewardTierRepo:  rewardTierRepo,
		uuidGen:         uuidGen,
		txManager:       txManager,
		idempotencyRepo: idempotencyRepo,
	}
	

	if idempotencyRepo == nil {
		svc.legacyMode = true
		// Initialize with no-op implementation
		svc.idempotencyRepo = &noopIdempotencyRepo{}
	}

	return svc
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
		
		// noopIdempotencyRepo provides a no-op implementation for backward compatibility
		type noopIdempotencyRepo struct{}
		
		func (r *noopIdempotencyRepo) Create(ctx context.Context, key *entity.IdempotencyKey) error {
			return nil
		}
		
		func (r *noopIdempotencyRepo) FindByKey(ctx context.Context, key string) (*entity.IdempotencyKey, error) {
			return nil, nil
		}
		
		func (r *noopIdempotencyRepo) Update(ctx context.Context, key *entity.IdempotencyKey) error {
			return nil
		}
		
		func (r *noopIdempotencyRepo) DeleteExpired(ctx context.Context) error {
			return nil
		}
		
		func (r *noopIdempotencyRepo) Purge(ctx context.Context, olderThan time.Time) error {
			return nil
		}
	}

	return allItems, nil
}

// PurchaseItem handles a user purchasing an item from the shop with idempotency
func (s *ShopService) PurchaseItem(ctx context.Context, userID int64, itemCode string, quantity int, idempotencyKey string) (*entity.Purchase, error) {
	var purchase *entity.Purchase

	// Check idempotency first
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
			// Try global items
			item, err = s.shopItemRepo.FindByCode(txCtx, 0, itemCode)
			if err != nil {
				return fmt.Errorf("item not found: %w", err)
			}
		}

		// Check if item is active
		if !item.IsActive {
			return fmt.Errorf("item is not available")
		}

		// Check stock
		if item.Stock != nil && *item.Stock < quantity {
			return fmt.Errorf("insufficient stock")
		}

		// Calculate total cost
		totalCost := item.Price.Mul(valueobject.NewDecimal(fmt.Sprintf("%d", quantity)))

		// Check user balance
		if user.Balance.Cmp(totalCost) < 0 {
			return fmt.Errorf("insufficient balance")
		}

		// Create purchase record
		purchase = &entity.Purchase{
			UserID:    userID,
			ItemID:    item.ID,
			ItemName:  item.Name,
			ItemPrice: item.Price,
			Quantity:  quantity,
			TotalCost: totalCost,
			Status:    "completed",
		}

		// 1. Deduct user balance
		negativeAmount := totalCost.Mul(valueobject.NewDecimal("-1"))
		err = s.userRepo.UpdateBalance(txCtx, userID, negativeAmount)
		if err != nil {
			return fmt.Errorf("failed to update balance: %w", err)
		}

		// 2. Update stock if limited
		if item.Stock != nil {
			newStock := *item.Stock - quantity
			item.Stock = &newStock
			err = s.shopItemRepo.Update(txCtx, item)
			if err != nil {
				return fmt.Errorf("failed to update stock: %w", err)
			}
		}

		// 3. Create purchase record
		err = s.purchaseRepo.Create(txCtx, purchase)
		if err != nil {
			return fmt.Errorf("failed to create purchase: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
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
