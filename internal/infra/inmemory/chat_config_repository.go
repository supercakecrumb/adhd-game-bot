package inmemory

import (
	"context"
	"errors"
	"sync"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/ports"
)

type ChatConfigRepository struct {
	mu      sync.RWMutex
	configs map[int64]*entity.ChatConfig
}

func NewChatConfigRepository() *ChatConfigRepository {
	return &ChatConfigRepository{
		configs: make(map[int64]*entity.ChatConfig),
	}
}

func (r *ChatConfigRepository) Create(ctx context.Context, config *entity.ChatConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.configs[config.ChatID]; exists {
		return errors.New("chat config already exists")
	}

	r.configs[config.ChatID] = config
	return nil
}

func (r *ChatConfigRepository) FindByChatID(ctx context.Context, chatID int64) (*entity.ChatConfig, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	config, exists := r.configs[chatID]
	if !exists {
		return nil, ports.ErrChatConfigNotFound
	}

	// Return a copy to prevent external modifications
	configCopy := *config
	return &configCopy, nil
}

func (r *ChatConfigRepository) Update(ctx context.Context, config *entity.ChatConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.configs[config.ChatID]; !exists {
		return ports.ErrChatConfigNotFound
	}

	r.configs[config.ChatID] = config
	return nil
}
