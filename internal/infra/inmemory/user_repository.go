package inmemory

import (
	"context"
	"sync"
	"time"

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

	// Set timestamps if not set
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}
	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = time.Now()
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

func (r *UserRepository) FindByChatID(ctx context.Context, chatID int64) ([]*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var users []*entity.User
	for _, u := range r.users {
		if u.ChatID == chatID {
			users = append(users, u)
		}
	}

	return users, nil
}

func (r *UserRepository) UpdateBalance(ctx context.Context, userID int64, delta valueobject.Decimal) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[userID]
	if !exists {
		return ports.ErrUserNotFound
	}

	user.Balance = user.Balance.Add(delta)
	user.UpdatedAt = time.Now()
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
