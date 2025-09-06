package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
)

type ShopItemRepository struct {
	db *sql.DB
}

func NewShopItemRepository(db *sql.DB) *ShopItemRepository {
	return &ShopItemRepository{db: db}
}

func (r *ShopItemRepository) Create(ctx context.Context, item *entity.ShopItem) error {
	// Check if we're in a transaction
	if tx, ok := GetTx(ctx); ok {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO shop_items (id, chat_id, dungeon_id, code, name, description, price, category, is_active, stock, discount_tier_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
			item.ID, item.ChatID, item.DungeonID, item.Code, item.Name, item.Description, item.Price.String(),
			item.Category, item.IsActive, item.Stock, item.DiscountTierID, item.CreatedAt, item.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to create shop item: %w", err)
		}
	} else {
		_, err := r.db.ExecContext(ctx, `
			INSERT INTO shop_items (id, chat_id, dungeon_id, code, name, description, price, category, is_active, stock, discount_tier_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
			item.ID, item.ChatID, item.DungeonID, item.Code, item.Name, item.Description, item.Price.String(),
			item.Category, item.IsActive, item.Stock, item.DiscountTierID, item.CreatedAt, item.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to create shop item: %w", err)
		}
	}

	return nil
}

func (r *ShopItemRepository) FindByID(ctx context.Context, id int64) (*entity.ShopItem, error) {
	var item entity.ShopItem
	var priceStr string
	var createdAt time.Time
	var updatedAt time.Time

	var row *sql.Row
	if tx, ok := GetTx(ctx); ok {
		row = tx.QueryRowContext(ctx, `
			SELECT id, chat_id, dungeon_id, code, name, description, price, category, is_active, stock, discount_tier_id, created_at, updated_at
			FROM shop_items WHERE id = $1`, id)
	} else {
		row = r.db.QueryRowContext(ctx, `
			SELECT id, chat_id, dungeon_id, code, name, description, price, category, is_active, stock, discount_tier_id, created_at, updated_at
			FROM shop_items WHERE id = $1`, id)
	}

	err := row.Scan(&item.ID, &item.ChatID, &item.DungeonID, &item.Code, &item.Name, &item.Description, &priceStr,
		&item.Category, &item.IsActive, &item.Stock, &item.DiscountTierID, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("shop item not found: %w", ErrItemNotFound)
		}
		return nil, fmt.Errorf("failed to query shop item: %w", err)
	}

	// Parse the price
	price := valueobject.NewDecimal(priceStr)

	item.Price = price
	item.CreatedAt = createdAt
	item.UpdatedAt = updatedAt

	return &item, nil
}

func (r *ShopItemRepository) FindByCode(ctx context.Context, chatID int64, code string) (*entity.ShopItem, error) {
	var item entity.ShopItem
	var priceStr string
	var createdAt time.Time
	var updatedAt time.Time

	var row *sql.Row
	if tx, ok := GetTx(ctx); ok {
		row = tx.QueryRowContext(ctx, `
			SELECT id, chat_id, dungeon_id, code, name, description, price, category, is_active, stock, discount_tier_id, created_at, updated_at
			FROM shop_items WHERE chat_id = $1 AND code = $2`, chatID, code)
	} else {
		row = r.db.QueryRowContext(ctx, `
			SELECT id, chat_id, dungeon_id, code, name, description, price, category, is_active, stock, discount_tier_id, created_at, updated_at
			FROM shop_items WHERE chat_id = $1 AND code = $2`, chatID, code)
	}

	err := row.Scan(&item.ID, &item.ChatID, &item.DungeonID, &item.Code, &item.Name, &item.Description, &priceStr,
		&item.Category, &item.IsActive, &item.Stock, &item.DiscountTierID, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No item found
		}
		return nil, fmt.Errorf("failed to query shop item: %w", err)
	}

	// Parse the price
	price := valueobject.NewDecimal(priceStr)

	item.Price = price
	item.CreatedAt = createdAt
	item.UpdatedAt = updatedAt

	return &item, nil
}

func (r *ShopItemRepository) FindByChatID(ctx context.Context, chatID int64) ([]*entity.ShopItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, chat_id, dungeon_id, code, name, description, price, category, is_active, stock, discount_tier_id, created_at, updated_at
		FROM shop_items WHERE chat_id = $1`, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to query shop items: %w", err)
	}
	defer rows.Close()

	var items []*entity.ShopItem
	for rows.Next() {
		var item entity.ShopItem
		var priceStr string
		var createdAt time.Time
		var updatedAt time.Time

		err := rows.Scan(&item.ID, &item.ChatID, &item.DungeonID, &item.Code, &item.Name, &item.Description, &priceStr,
			&item.Category, &item.IsActive, &item.Stock, &item.DiscountTierID, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan shop item: %w", err)
		}

		// Parse the price
		price := valueobject.NewDecimal(priceStr)

		item.Price = price
		item.CreatedAt = createdAt
		item.UpdatedAt = updatedAt

		items = append(items, &item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over shop item rows: %w", err)
	}

	return items, nil
}

func (r *ShopItemRepository) Update(ctx context.Context, item *entity.ShopItem) error {
	// Check if we're in a transaction
	if tx, ok := GetTx(ctx); ok {
		_, err := tx.ExecContext(ctx, `
			UPDATE shop_items 
			SET code = $1, name = $2, description = $3, price = $4, category = $5, is_active = $6, stock = $7, discount_tier_id = $8, updated_at = $9
			WHERE id = $10`,
			item.Code, item.Name, item.Description, item.Price.String(), item.Category, item.IsActive,
			item.Stock, item.DiscountTierID, item.UpdatedAt, item.ID)
		if err != nil {
			return fmt.Errorf("failed to update shop item: %w", err)
		}
	} else {
		_, err := r.db.ExecContext(ctx, `
			UPDATE shop_items 
			SET code = $1, name = $2, description = $3, price = $4, category = $5, is_active = $6, stock = $7, discount_tier_id = $8, updated_at = $9
			WHERE id = $10`,
			item.Code, item.Name, item.Description, item.Price.String(), item.Category, item.IsActive,
			item.Stock, item.DiscountTierID, item.UpdatedAt, item.ID)
		if err != nil {
			return fmt.Errorf("failed to update shop item: %w", err)
		}
	}

	return nil
}

func (r *ShopItemRepository) Delete(ctx context.Context, id int64) error {
	// Check if we're in a transaction
	if tx, ok := GetTx(ctx); ok {
		_, err := tx.ExecContext(ctx, `DELETE FROM shop_items WHERE id = $1`, id)
		if err != nil {
			return fmt.Errorf("failed to delete shop item: %w", err)
		}
	} else {
		_, err := r.db.ExecContext(ctx, `DELETE FROM shop_items WHERE id = $1`, id)
		if err != nil {
			return fmt.Errorf("failed to delete shop item: %w", err)
		}
	}

	return nil
}
