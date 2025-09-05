package inmemory

import (
	"context"
	"sync"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/ports"
)

type InMemoryDiscountTierRepository struct {
	mu    sync.RWMutex
	tiers map[int64]*entity.DiscountTier
}

func NewDiscountTierRepository() *InMemoryDiscountTierRepository {
	return &InMemoryDiscountTierRepository{
		tiers: make(map[int64]*entity.DiscountTier),
	}
}

func (r *InMemoryDiscountTierRepository) Create(ctx context.Context, tier *entity.DiscountTier) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tiers[tier.ID]; exists {
		return ports.ErrDiscountTierExists
	}

	r.tiers[tier.ID] = tier
	return nil
}

func (r *InMemoryDiscountTierRepository) FindByID(ctx context.Context, id int64) (*entity.DiscountTier, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tier, exists := r.tiers[id]
	if !exists {
		return nil, ports.ErrDiscountTierNotFound
	}
	return tier, nil
}

func (r *InMemoryDiscountTierRepository) FindAll(ctx context.Context) ([]*entity.DiscountTier, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tiers := make([]*entity.DiscountTier, 0, len(r.tiers))
	for _, t := range r.tiers {
		tiers = append(tiers, t)
	}
	return tiers, nil
}

func (r *InMemoryDiscountTierRepository) Update(ctx context.Context, tier *entity.DiscountTier) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tiers[tier.ID]; !exists {
		return ports.ErrDiscountTierNotFound
	}

	r.tiers[tier.ID] = tier
	return nil
}

func (r *InMemoryDiscountTierRepository) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.tiers, id)
	return nil
}
