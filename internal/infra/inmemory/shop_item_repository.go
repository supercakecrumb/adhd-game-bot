package inmemory

import (
	"context"
	"sync"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/ports"
)

type ShopItemRepository struct {
	mu     sync.RWMutex
	items  map[int64]*entity.ShopItem
	nextID int64
}

func NewShopItemRepository() *ShopItemRepository {
	return &ShopItemRepository{
		items:  make(map[int64]*entity.ShopItem),
		nextID: 1,
	}
}

func (r *ShopItemRepository) Create(ctx context.Context, item *entity.ShopItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if item.ID == 0 {
		item.ID = r.nextID
		r.nextID++
	}

	r.items[item.ID] = item
	return nil
}

func (r *ShopItemRepository) FindByID(ctx context.Context, id int64) (*entity.ShopItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.items[id]
	if !exists {
		return nil, ports.ErrShopItemNotFound
	}

	// Return a copy to prevent external modifications
	itemCopy := *item
	if item.Stock != nil {
		stockCopy := *item.Stock
		itemCopy.Stock = &stockCopy
	}
	return &itemCopy, nil
}

func (r *ShopItemRepository) FindByCode(ctx context.Context, chatID int64, code string) (*entity.ShopItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, item := range r.items {
		if item.ChatID == chatID && item.Code == code {
			// Return a copy
			itemCopy := *item
			if item.Stock != nil {
				stockCopy := *item.Stock
				itemCopy.Stock = &stockCopy
			}
			return &itemCopy, nil
		}
	}

	return nil, ports.ErrShopItemNotFound
}

func (r *ShopItemRepository) FindByChatID(ctx context.Context, chatID int64) ([]*entity.ShopItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var items []*entity.ShopItem
	for _, item := range r.items {
		if item.ChatID == chatID {
			// Create a copy
			itemCopy := *item
			if item.Stock != nil {
				stockCopy := *item.Stock
				itemCopy.Stock = &stockCopy
			}
			items = append(items, &itemCopy)
		}
	}

	return items, nil
}

func (r *ShopItemRepository) Update(ctx context.Context, item *entity.ShopItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.items[item.ID]
	if !exists {
		return ports.ErrShopItemNotFound
	}

	// Check stock constraints
	if item.Stock != nil && existing.Stock != nil {
		if *item.Stock < 0 {
			return ports.ErrInsufficientStock
		}
	}

	// Create a copy to store
	itemCopy := *item
	if item.Stock != nil {
		stockCopy := *item.Stock
		itemCopy.Stock = &stockCopy
	}
	r.items[item.ID] = &itemCopy
	return nil
}

func (r *ShopItemRepository) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[id]; !exists {
		return ports.ErrShopItemNotFound
	}

	delete(r.items, id)
	return nil
}
