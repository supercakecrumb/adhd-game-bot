package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/ports"
)

type ChatConfigRepository struct {
	db *sql.DB
}

func NewChatConfigRepository(db *sql.DB) *ChatConfigRepository {
	return &ChatConfigRepository{db: db}
}

func (r *ChatConfigRepository) Create(ctx context.Context, config *entity.ChatConfig) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO chat_configs (chat_id, currency_name, created_at, updated_at)
		VALUES ($1, $2, $3, $4)`,
		config.ChatID, config.CurrencyName, config.CreatedAt, config.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create chat config: %w", err)
	}

	return nil
}

func (r *ChatConfigRepository) FindByChatID(ctx context.Context, chatID int64) (*entity.ChatConfig, error) {
	var config entity.ChatConfig

	err := r.db.QueryRowContext(ctx, `
		SELECT chat_id, currency_name, created_at, updated_at
		FROM chat_configs WHERE chat_id = $1`, chatID).Scan(
		&config.ChatID, &config.CurrencyName, &config.CreatedAt, &config.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ports.ErrChatConfigNotFound
		}
		return nil, fmt.Errorf("failed to query chat config: %w", err)
	}

	return &config, nil
}

func (r *ChatConfigRepository) Update(ctx context.Context, config *entity.ChatConfig) error {
	config.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, `
		UPDATE chat_configs
		SET currency_name = $2, updated_at = $3
		WHERE chat_id = $1`,
		config.ChatID, config.CurrencyName, config.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update chat config: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ports.ErrChatConfigNotFound
	}

	return nil
}

func (r *ChatConfigRepository) Delete(ctx context.Context, chatID int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM chat_configs WHERE chat_id = $1", chatID)
	if err != nil {
		return fmt.Errorf("failed to delete chat config: %w", err)
	}
	return nil
}
