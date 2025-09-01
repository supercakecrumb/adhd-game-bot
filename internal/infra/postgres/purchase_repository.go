package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
	"github.com/supercakecrumb/adhd-game-bot/internal/ports"
)

type PurchaseRepository struct {
	db *sql.DB
}

func NewPurchaseRepository(db *sql.DB) *PurchaseRepository {
	return &PurchaseRepository{db: db}
}

func (r *PurchaseRepository) Create(ctx context.Context, purchase *entity.Purchase) error {
	var row *sql.Row
	if tx, ok := GetTx(ctx); ok {
		row = tx.QueryRowContext(ctx, `
			INSERT INTO purchases (user_id, item_id, item_name, item_price, quantity, total_cost, status, purchased_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id`,
			purchase.UserID, purchase.ItemID, purchase.ItemName, purchase.ItemPrice.String(),
			purchase.Quantity, purchase.TotalCost.String(), purchase.Status, purchase.PurchasedAt)
	} else {
		row = r.db.QueryRowContext(ctx, `
			INSERT INTO purchases (user_id, item_id, item_name, item_price, quantity, total_cost, status, purchased_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id`,
			purchase.UserID, purchase.ItemID, purchase.ItemName, purchase.ItemPrice.String(),
			purchase.Quantity, purchase.TotalCost.String(), purchase.Status, purchase.PurchasedAt)
	}

	if err := row.Scan(&purchase.ID); err != nil {
		return fmt.Errorf("failed to create purchase: %w", err)
	}

	return nil
}

func (r *PurchaseRepository) FindByID(ctx context.Context, id int64) (*entity.Purchase, error) {
	var purchase entity.Purchase
	var itemPriceStr, totalCostStr string

	var row *sql.Row
	if tx, ok := GetTx(ctx); ok {
		row = tx.QueryRowContext(ctx, `
			SELECT id, user_id, item_id, item_name, item_price, quantity, total_cost, status, purchased_at
			FROM purchases WHERE id = $1`, id)
	} else {
		row = r.db.QueryRowContext(ctx, `
			SELECT id, user_id, item_id, item_name, item_price, quantity, total_cost, status, purchased_at
			FROM purchases WHERE id = $1`, id)
	}

	err := row.Scan(&purchase.ID, &purchase.UserID, &purchase.ItemID, &purchase.ItemName,
		&itemPriceStr, &purchase.Quantity, &totalCostStr, &purchase.Status, &purchase.PurchasedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ports.ErrPurchaseNotFound
		}
		return nil, fmt.Errorf("failed to query purchase: %w", err)
	}

	// Parse decimals
	purchase.ItemPrice = valueobject.NewDecimal(itemPriceStr)
	purchase.TotalCost = valueobject.NewDecimal(totalCostStr)

	return &purchase, nil
}

func (r *PurchaseRepository) FindByUserID(ctx context.Context, userID int64) ([]*entity.Purchase, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, item_id, item_name, item_price, quantity, total_cost, status, purchased_at
		FROM purchases WHERE user_id = $1 ORDER BY purchased_at DESC`, userID)

	if err != nil {
		return nil, fmt.Errorf("failed to query purchases: %w", err)
	}
	defer rows.Close()

	var purchases []*entity.Purchase
	for rows.Next() {
		var purchase entity.Purchase
		var itemPriceStr, totalCostStr string

		err := rows.Scan(
			&purchase.ID, &purchase.UserID, &purchase.ItemID, &purchase.ItemName,
			&itemPriceStr, &purchase.Quantity, &totalCostStr, &purchase.Status, &purchase.PurchasedAt)

		if err != nil {
			return nil, fmt.Errorf("failed to scan purchase: %w", err)
		}

		// Parse decimals
		purchase.ItemPrice = valueobject.NewDecimal(itemPriceStr)
		purchase.TotalCost = valueobject.NewDecimal(totalCostStr)

		purchases = append(purchases, &purchase)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate purchases: %w", err)
	}

	return purchases, nil
}

func (r *PurchaseRepository) FindByItemID(ctx context.Context, itemID int64) ([]*entity.Purchase, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, item_id, item_name, item_price, quantity, total_cost, status, purchased_at
		FROM purchases WHERE item_id = $1 ORDER BY purchased_at DESC`, itemID)

	if err != nil {
		return nil, fmt.Errorf("failed to query purchases: %w", err)
	}
	defer rows.Close()

	var purchases []*entity.Purchase
	for rows.Next() {
		var purchase entity.Purchase
		var itemPriceStr, totalCostStr string

		err := rows.Scan(
			&purchase.ID, &purchase.UserID, &purchase.ItemID, &purchase.ItemName,
			&itemPriceStr, &purchase.Quantity, &totalCostStr, &purchase.Status, &purchase.PurchasedAt)

		if err != nil {
			return nil, fmt.Errorf("failed to scan purchase: %w", err)
		}

		// Parse decimals
		purchase.ItemPrice = valueobject.NewDecimal(itemPriceStr)
		purchase.TotalCost = valueobject.NewDecimal(totalCostStr)

		purchases = append(purchases, &purchase)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate purchases: %w", err)
	}

	return purchases, nil
}
