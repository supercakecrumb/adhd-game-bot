package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
)

type PurchaseRepository struct {
	db *sql.DB
}

func NewPurchaseRepository(db *sql.DB) *PurchaseRepository {
	return &PurchaseRepository{db: db}
}

func (r *PurchaseRepository) Create(ctx context.Context, purchase *entity.Purchase) error {
	// Check if we're in a transaction
	if tx, ok := GetTx(ctx); ok {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO purchases (id, user_id, item_id, dungeon_id, item_name, item_price, quantity, total_cost, status, discount_tier_id, purchased_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
			purchase.ID, purchase.UserID, purchase.ItemID, purchase.DungeonID, purchase.ItemName,
			purchase.ItemPrice.String(), purchase.Quantity, purchase.TotalCost.String(), purchase.Status,
			purchase.DiscountTierID, purchase.PurchasedAt)
		if err != nil {
			return fmt.Errorf("failed to create purchase: %w", err)
		}
	} else {
		_, err := r.db.ExecContext(ctx, `
			INSERT INTO purchases (id, user_id, item_id, dungeon_id, item_name, item_price, quantity, total_cost, status, discount_tier_id, purchased_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
			purchase.ID, purchase.UserID, purchase.ItemID, purchase.DungeonID, purchase.ItemName,
			purchase.ItemPrice.String(), purchase.Quantity, purchase.TotalCost.String(), purchase.Status,
			purchase.DiscountTierID, purchase.PurchasedAt)
		if err != nil {
			return fmt.Errorf("failed to create purchase: %w", err)
		}
	}

	return nil
}

func (r *PurchaseRepository) FindByID(ctx context.Context, id int64) (*entity.Purchase, error) {
	var purchase entity.Purchase
	var itemPriceStr string
	var totalCostStr string
	var purchasedAt time.Time

	var row *sql.Row
	if tx, ok := GetTx(ctx); ok {
		row = tx.QueryRowContext(ctx, `
			SELECT id, user_id, item_id, dungeon_id, item_name, item_price, quantity, total_cost, status, discount_tier_id, purchased_at
			FROM purchases WHERE id = $1`, id)
	} else {
		row = r.db.QueryRowContext(ctx, `
			SELECT id, user_id, item_id, dungeon_id, item_name, item_price, quantity, total_cost, status, discount_tier_id, purchased_at
			FROM purchases WHERE id = $1`, id)
	}

	err := row.Scan(&purchase.ID, &purchase.UserID, &purchase.ItemID, &purchase.DungeonID, &purchase.ItemName,
		&itemPriceStr, &purchase.Quantity, &totalCostStr, &purchase.Status, &purchase.DiscountTierID, &purchasedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("purchase not found: %w", ErrPurchaseNotFound)
		}
		return nil, fmt.Errorf("failed to query purchase: %w", err)
	}

	// Parse the prices
	itemPrice := valueobject.NewDecimal(itemPriceStr)
	totalCost := valueobject.NewDecimal(totalCostStr)

	purchase.ItemPrice = itemPrice
	purchase.TotalCost = totalCost
	purchase.PurchasedAt = purchasedAt

	return &purchase, nil
}

func (r *PurchaseRepository) FindByUserID(ctx context.Context, userID int64) ([]*entity.Purchase, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, item_id, dungeon_id, item_name, item_price, quantity, total_cost, status, discount_tier_id, purchased_at
		FROM purchases WHERE user_id = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query purchases: %w", err)
	}
	defer rows.Close()

	var purchases []*entity.Purchase
	for rows.Next() {
		var purchase entity.Purchase
		var itemPriceStr string
		var totalCostStr string
		var purchasedAt time.Time

		err := rows.Scan(&purchase.ID, &purchase.UserID, &purchase.ItemID, &purchase.DungeonID, &purchase.ItemName,
			&itemPriceStr, &purchase.Quantity, &totalCostStr, &purchase.Status, &purchase.DiscountTierID, &purchasedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan purchase: %w", err)
		}

		// Parse the prices
		itemPrice := valueobject.NewDecimal(itemPriceStr)
		totalCost := valueobject.NewDecimal(totalCostStr)

		purchase.ItemPrice = itemPrice
		purchase.TotalCost = totalCost
		purchase.PurchasedAt = purchasedAt

		purchases = append(purchases, &purchase)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over purchase rows: %w", err)
	}

	return purchases, nil
}

func (r *PurchaseRepository) FindByItemID(ctx context.Context, itemID int64) ([]*entity.Purchase, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, item_id, dungeon_id, item_name, item_price, quantity, total_cost, status, discount_tier_id, purchased_at
		FROM purchases WHERE item_id = $1`, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to query purchases: %w", err)
	}
	defer rows.Close()

	var purchases []*entity.Purchase
	for rows.Next() {
		var purchase entity.Purchase
		var itemPriceStr string
		var totalCostStr string
		var purchasedAt time.Time

		err := rows.Scan(&purchase.ID, &purchase.UserID, &purchase.ItemID, &purchase.DungeonID, &purchase.ItemName,
			&itemPriceStr, &purchase.Quantity, &totalCostStr, &purchase.Status, &purchase.DiscountTierID, &purchasedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan purchase: %w", err)
		}

		// Parse the prices
		itemPrice := valueobject.NewDecimal(itemPriceStr)
		totalCost := valueobject.NewDecimal(totalCostStr)

		purchase.ItemPrice = itemPrice
		purchase.TotalCost = totalCost
		purchase.PurchasedAt = purchasedAt

		purchases = append(purchases, &purchase)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over purchase rows: %w", err)
	}

	return purchases, nil
}
