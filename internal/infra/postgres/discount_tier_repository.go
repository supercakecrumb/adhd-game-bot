package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
)

type DiscountTierRepository struct {
	db *sql.DB
}

func NewDiscountTierRepository(db *sql.DB) *DiscountTierRepository {
	return &DiscountTierRepository{db: db}
}

func (r *DiscountTierRepository) Create(ctx context.Context, tier *entity.DiscountTier) error {
	// Check if we're in a transaction
	if tx, ok := GetTx(ctx); ok {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO discount_tiers (id, name, description, discount_percent, min_purchases, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			tier.ID, tier.Name, tier.Description, tier.DiscountPercent, tier.MinPurchases, tier.CreatedAt, tier.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to create discount tier: %w", err)
		}
	} else {
		_, err := r.db.ExecContext(ctx, `
			INSERT INTO discount_tiers (id, name, description, discount_percent, min_purchases, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			tier.ID, tier.Name, tier.Description, tier.DiscountPercent, tier.MinPurchases, tier.CreatedAt, tier.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to create discount tier: %w", err)
		}
	}

	return nil
}

func (r *DiscountTierRepository) FindByID(ctx context.Context, id int64) (*entity.DiscountTier, error) {
	var tier entity.DiscountTier
	var createdAt time.Time
	var updatedAt time.Time

	var row *sql.Row
	if tx, ok := GetTx(ctx); ok {
		row = tx.QueryRowContext(ctx, `
			SELECT id, name, description, discount_percent, min_purchases, created_at, updated_at
			FROM discount_tiers WHERE id = $1`, id)
	} else {
		row = r.db.QueryRowContext(ctx, `
			SELECT id, name, description, discount_percent, min_purchases, created_at, updated_at
			FROM discount_tiers WHERE id = $1`, id)
	}

	err := row.Scan(&tier.ID, &tier.Name, &tier.Description, &tier.DiscountPercent, &tier.MinPurchases, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("discount tier not found: %w", ErrDiscountTierNotFound)
		}
		return nil, fmt.Errorf("failed to query discount tier: %w", err)
	}

	tier.CreatedAt = createdAt
	tier.UpdatedAt = updatedAt

	return &tier, nil
}

func (r *DiscountTierRepository) FindAll(ctx context.Context) ([]*entity.DiscountTier, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, description, discount_percent, min_purchases, created_at, updated_at
		FROM discount_tiers ORDER BY discount_percent ASC`)
	if err != nil {
		return nil, fmt.Errorf("failed to query discount tiers: %w", err)
	}
	defer rows.Close()

	var tiers []*entity.DiscountTier
	for rows.Next() {
		var tier entity.DiscountTier
		var createdAt time.Time
		var updatedAt time.Time

		err := rows.Scan(&tier.ID, &tier.Name, &tier.Description, &tier.DiscountPercent, &tier.MinPurchases, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan discount tier: %w", err)
		}

		tier.CreatedAt = createdAt
		tier.UpdatedAt = updatedAt

		tiers = append(tiers, &tier)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over discount tier rows: %w", err)
	}

	return tiers, nil
}

func (r *DiscountTierRepository) Update(ctx context.Context, tier *entity.DiscountTier) error {
	// Check if we're in a transaction
	if tx, ok := GetTx(ctx); ok {
		_, err := tx.ExecContext(ctx, `
			UPDATE discount_tiers 
			SET name = $1, description = $2, discount_percent = $3, min_purchases = $4, updated_at = $5
			WHERE id = $6`,
			tier.Name, tier.Description, tier.DiscountPercent, tier.MinPurchases, tier.UpdatedAt, tier.ID)
		if err != nil {
			return fmt.Errorf("failed to update discount tier: %w", err)
		}
	} else {
		_, err := r.db.ExecContext(ctx, `
			UPDATE discount_tiers 
			SET name = $1, description = $2, discount_percent = $3, min_purchases = $4, updated_at = $5
			WHERE id = $6`,
			tier.Name, tier.Description, tier.DiscountPercent, tier.MinPurchases, tier.UpdatedAt, tier.ID)
		if err != nil {
			return fmt.Errorf("failed to update discount tier: %w", err)
		}
	}

	return nil
}

func (r *DiscountTierRepository) Delete(ctx context.Context, id int64) error {
	// Check if we're in a transaction
	if tx, ok := GetTx(ctx); ok {
		_, err := tx.ExecContext(ctx, `DELETE FROM discount_tiers WHERE id = $1`, id)
		if err != nil {
			return fmt.Errorf("failed to delete discount tier: %w", err)
		}
	} else {
		_, err := r.db.ExecContext(ctx, `DELETE FROM discount_tiers WHERE id = $1`, id)
		if err != nil {
			return fmt.Errorf("failed to delete discount tier: %w", err)
		}
	}

	return nil
}
