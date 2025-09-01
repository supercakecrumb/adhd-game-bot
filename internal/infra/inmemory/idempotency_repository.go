package inmemory

import (
	"context"
	"sync"
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/ports"
)

type InMemoryIdempotencyRepository struct {
	mu   sync.RWMutex
	keys map[string]*entity.IdempotencyKey
}

func NewInMemoryIdempotencyRepository() *InMemoryIdempotencyRepository {
	return &InMemoryIdempotencyRepository{
		keys: make(map[string]*entity.IdempotencyKey),
	}
}

func (r *InMemoryIdempotencyRepository) Create(ctx context.Context, key *entity.IdempotencyKey) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.keys[key.Key]; exists {
		return ports.ErrIdempotencyKeyExists
	}

	r.keys[key.Key] = key
	return nil
}

func (r *InMemoryIdempotencyRepository) FindByKey(ctx context.Context, key string) (*entity.IdempotencyKey, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	idempKey, exists := r.keys[key]
	if !exists {
		return nil, ports.ErrIdempotencyKeyNotFound
	}

	return idempKey, nil
}

func (r *InMemoryIdempotencyRepository) Update(ctx context.Context, key *entity.IdempotencyKey) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.keys[key.Key]; !exists {
		return ports.ErrIdempotencyKeyNotFound
	}

	r.keys[key.Key] = key
	return nil
}

func (r *InMemoryIdempotencyRepository) DeleteExpired(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for k, v := range r.keys {
		if now.After(v.ExpiresAt) {
			delete(r.keys, k)
		}
	}

	return nil
}
