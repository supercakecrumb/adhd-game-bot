package inmemory

import (
	"context"
	"sync"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
	"github.com/supercakecrumb/adhd-game-bot/internal/ports"
)

type UserRepository struct {
	mu    sync.RWMutex
	users map[int64]*entity.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[int64]*entity.User),
	}
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; exists {
		return ports.ErrUserAlreadyExists
	}

	r.users[user.ID] = user
	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, ports.ErrUserNotFound
	}

	return user, nil
}

func (r *UserRepository) UpdateBalance(ctx context.Context, userID int64, currency string, delta valueobject.Decimal) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[userID]
	if !exists {
		return ports.ErrUserNotFound
	}

	if user.Balances == nil {
		user.Balances = make(map[string]valueobject.Decimal)
	}

	current := user.Balances[currency]
	user.Balances[currency] = current.Add(delta)
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return ports.ErrUserNotFound
	}

	delete(r.users, id)
	return nil
}
