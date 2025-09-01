package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
	"github.com/supercakecrumb/adhd-game-bot/internal/ports"
)

type ShopItemRepository struct {
	db *sql.DB
}

func NewShopItemRepository(db *sql.DB) *ShopItemRepository {
	return &ShopItemRepository{db: db}
}

func (r *ShopItemRepository) Create(ctx context.Context, item *entity.ShopItem) error {
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO shop_items (chat_id, code, name, description, price, category, is_active, stock, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`,
		item.ChatID, item.Code, item.Name, item.Description, item.Price.String(),
		item.Category, item.IsActive, item.Stock, item.CreatedAt, item.UpdatedAt).Scan(&item.ID)

	if err != nil {
		return fmt.Errorf("failed to create shop item: %w", err)
	}

	return nil
}

func (r *ShopItemRepository) FindByID(ctx context.Context, id int64) (*entity.ShopItem, error) {
	var item entity.ShopItem
	var priceStr string

	err := r.db.QueryRowContext(ctx, `
		SELECT id, chat_id, code, name, description, price, category, is_active, stock, created_at, updated_at
		FROM shop_items WHERE id = $1`, id).Scan(
		&item.ID, &item.ChatID, &item.Code, &item.Name, &item.Description,
		&priceStr, &item.Category, &item.IsActive, &item.Stock, &item.CreatedAt, &item.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ports.ErrShopItemNotFound
		}
		return nil, fmt.Errorf("failed to query shop item: %w", err)
	}

	// Parse price
	item.Price = valueobject.NewDecimal(priceStr)

	return &item, nil
}

func (r *ShopItemRepository) FindByCode(ctx context.Context, chatID int64, code string) (*entity.ShopItem, error) {
	var item entity.ShopItem
	var priceStr string

	err := r.db.QueryRowContext(ctx, `
		SELECT id, chat_id, code, name, description, price, category, is_active, stock, created_at, updated_at
		FROM shop_items WHERE chat_id = $1 AND code = $2`, chatID, code).Scan(
		&item.ID, &item.ChatID, &item.Code, &item.Name, &item.Description,
		&priceStr, &item.Category, &item.IsActive, &item.Stock, &item.CreatedAt, &item.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ports.ErrShopItemNotFound
		}
		return nil, fmt.Errorf("failed to query shop item: %w", err)
	}

	// Parse price
	item.Price = valueobject.NewDecimal(priceStr)

	return &item, nil
}

func (r *ShopItemRepository) FindByChatID(ctx context.Context, chatID int64) ([]*entity.ShopItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, chat_id, code, name, description, price, category, is_active, stock, created_at, updated_at
		FROM shop_items WHERE chat_id = $1 ORDER BY created_at DESC`, chatID)

	if err != nil {
		return nil, fmt.Errorf("failed to query shop items: %w", err)
	}
	defer rows.Close()

	var items []*entity.ShopItem
	for rows.Next() {
		var item entity.ShopItem
		var priceStr string

		err := rows.Scan(
			&item.ID, &item.ChatID, &item.Code, &item.Name, &item.Description,
			&priceStr, &item.Category, &item.IsActive, &item.Stock, &item.CreatedAt, &item.UpdatedAt)

		if err != nil {
			return nil, fmt.Errorf("failed to scan shop item: %w", err)
		}

		// Parse price
		item.Price = valueobject.NewDecimal(priceStr)

		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate shop items: %w", err)
	}

	return items, nil
}

func (r *ShopItemRepository) Update(ctx context.Context, item *entity.ShopItem) error {
	item.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, `
		UPDATE shop_items
		SET code = $2, name = $3, description = $4, price = $5, category = $6, is_active = $7, stock = $8, updated_at = $9
		WHERE id = $1`,
		item.ID, item.Code, item.Name, item.Description, item.Price.String(),
		item.Category, item.IsActive, item.Stock, item.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update shop item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ports.ErrShopItemNotFound
	}

	return nil
}

func (r *ShopItemRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM shop_items WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete shop item: %w", err)
	}
	return nil
}
