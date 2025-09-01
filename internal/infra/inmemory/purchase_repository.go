package inmemory

import (
	"context"
	"sync"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/ports"
)

type PurchaseRepository struct {
	mu        sync.RWMutex
	purchases map[int64]*entity.Purchase
	nextID    int64
}

func NewPurchaseRepository() *PurchaseRepository {
	return &PurchaseRepository{
		purchases: make(map[int64]*entity.Purchase),
		nextID:    1,
	}
}

func (r *PurchaseRepository) Create(ctx context.Context, purchase *entity.Purchase) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if purchase.ID == 0 {
		purchase.ID = r.nextID
		r.nextID++
	}

	r.purchases[purchase.ID] = purchase
	return nil
}

func (r *PurchaseRepository) FindByID(ctx context.Context, id int64) (*entity.Purchase, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	purchase, exists := r.purchases[id]
	if !exists {
		return nil, ports.ErrPurchaseNotFound
	}

	// Return a copy to prevent external modifications
	purchaseCopy := *purchase
	return &purchaseCopy, nil
}

func (r *PurchaseRepository) FindByUserID(ctx context.Context, userID int64) ([]*entity.Purchase, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var purchases []*entity.Purchase
	for _, purchase := range r.purchases {
		if purchase.UserID == userID {
			// Create a copy
			purchaseCopy := *purchase
			purchases = append(purchases, &purchaseCopy)
		}
	}

	return purchases, nil
}

func (r *PurchaseRepository) FindByItemID(ctx context.Context, itemID int64) ([]*entity.Purchase, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var purchases []*entity.Purchase
	for _, purchase := range r.purchases {
		if purchase.ItemID == itemID {
			// Create a copy
			purchaseCopy := *purchase
			purchases = append(purchases, &purchaseCopy)
		}
	}

	return purchases, nil
}
