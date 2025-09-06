package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
)

type ChatConfigRepository struct {
	db *sql.DB
}

func NewChatConfigRepository(db *sql.DB) *ChatConfigRepository {
	return &ChatConfigRepository{db: db}
}

func (r *ChatConfigRepository) Create(ctx context.Context, config *entity.ChatConfig) error {
	// Check if we're in a transaction
	if tx, ok := GetTx(ctx); ok {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO chat_configs (chat_id, currency_name, created_at, updated_at)
			VALUES ($1, $2, $3, $4)`,
			config.ChatID, config.CurrencyName, config.CreatedAt, config.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to create chat config: %w", err)
		}
	} else {
		_, err := r.db.ExecContext(ctx, `
			INSERT INTO chat_configs (chat_id, currency_name, created_at, updated_at)
			VALUES ($1, $2, $3, $4)`,
			config.ChatID, config.CurrencyName, config.CreatedAt, config.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to create chat config: %w", err)
		}
	}

	return nil
}

func (r *ChatConfigRepository) FindByChatID(ctx context.Context, chatID int64) (*entity.ChatConfig, error) {
	var config entity.ChatConfig
	var createdAt time.Time
	var updatedAt time.Time

	var row *sql.Row
	if tx, ok := GetTx(ctx); ok {
		row = tx.QueryRowContext(ctx, `
			SELECT chat_id, currency_name, created_at, updated_at
			FROM chat_configs WHERE chat_id = $1`, chatID)
	} else {
		row = r.db.QueryRowContext(ctx, `
			SELECT chat_id, currency_name, created_at, updated_at
			FROM chat_configs WHERE chat_id = $1`, chatID)
	}

	err := row.Scan(&config.ChatID, &config.CurrencyName, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No config found
		}
		return nil, fmt.Errorf("failed to query chat config: %w", err)
	}

	config.CreatedAt = createdAt
	config.UpdatedAt = updatedAt

	return &config, nil
}

func (r *ChatConfigRepository) Update(ctx context.Context, config *entity.ChatConfig) error {
	// Check if we're in a transaction
	if tx, ok := GetTx(ctx); ok {
		_, err := tx.ExecContext(ctx, `
			UPDATE chat_configs 
			SET currency_name = $1, updated_at = $2
			WHERE chat_id = $3`,
			config.CurrencyName, config.UpdatedAt, config.ChatID)
		if err != nil {
			return fmt.Errorf("failed to update chat config: %w", err)
		}
	} else {
		_, err := r.db.ExecContext(ctx, `
			UPDATE chat_configs 
			SET currency_name = $1, updated_at = $2
			WHERE chat_id = $3`,
			config.CurrencyName, config.UpdatedAt, config.ChatID)
		if err != nil {
			return fmt.Errorf("failed to update chat config: %w", err)
		}
	}

	return nil
}
