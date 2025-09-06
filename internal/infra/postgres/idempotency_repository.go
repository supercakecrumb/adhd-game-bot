package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
)

type IdempotencyRepository struct {
	db *sql.DB
}

func NewIdempotencyRepository(db *sql.DB) *IdempotencyRepository {
	return &IdempotencyRepository{db: db}
}

func (r *IdempotencyRepository) Create(ctx context.Context, key *entity.IdempotencyKey) error {
	// Check if we're in a transaction
	if tx, ok := GetTx(ctx); ok {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO idempotency_keys (key, operation, user_id, status, result, created_at, completed_at, expires_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			key.Key, key.Operation, key.UserID, key.Status, key.Result, key.CreatedAt, key.CompletedAt, key.ExpiresAt)
		if err != nil {
			return fmt.Errorf("failed to create idempotency key: %w", err)
		}
	} else {
		_, err := r.db.ExecContext(ctx, `
			INSERT INTO idempotency_keys (key, operation, user_id, status, result, created_at, completed_at, expires_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			key.Key, key.Operation, key.UserID, key.Status, key.Result, key.CreatedAt, key.CompletedAt, key.ExpiresAt)
		if err != nil {
			return fmt.Errorf("failed to create idempotency key: %w", err)
		}
	}

	return nil
}

func (r *IdempotencyRepository) FindByKey(ctx context.Context, key string) (*entity.IdempotencyKey, error) {
	var idempotencyKey entity.IdempotencyKey
	var createdAt time.Time
	var completedAt *time.Time
	var expiresAt time.Time

	var row *sql.Row
	if tx, ok := GetTx(ctx); ok {
		row = tx.QueryRowContext(ctx, `
			SELECT key, operation, user_id, status, result, created_at, completed_at, expires_at
			FROM idempotency_keys WHERE key = $1 AND expires_at > NOW()`, key)
	} else {
		row = r.db.QueryRowContext(ctx, `
			SELECT key, operation, user_id, status, result, created_at, completed_at, expires_at
			FROM idempotency_keys WHERE key = $1 AND expires_at > NOW()`, key)
	}

	err := row.Scan(&idempotencyKey.Key, &idempotencyKey.Operation, &idempotencyKey.UserID, &idempotencyKey.Status,
		&idempotencyKey.Result, &createdAt, &completedAt, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No key found or expired
		}
		return nil, fmt.Errorf("failed to query idempotency key: %w", err)
	}

	idempotencyKey.CreatedAt = createdAt
	idempotencyKey.CompletedAt = completedAt
	idempotencyKey.ExpiresAt = expiresAt

	return &idempotencyKey, nil
}

func (r *IdempotencyRepository) Update(ctx context.Context, key *entity.IdempotencyKey) error {
	// Check if we're in a transaction
	if tx, ok := GetTx(ctx); ok {
		_, err := tx.ExecContext(ctx, `
			UPDATE idempotency_keys 
			SET status = $1, result = $2, completed_at = $3, expires_at = $4
			WHERE key = $5`,
			key.Status, key.Result, key.CompletedAt, key.ExpiresAt, key.Key)
		if err != nil {
			return fmt.Errorf("failed to update idempotency key: %w", err)
		}
	} else {
		_, err := r.db.ExecContext(ctx, `
			UPDATE idempotency_keys 
			SET status = $1, result = $2, completed_at = $3, expires_at = $4
			WHERE key = $5`,
			key.Status, key.Result, key.CompletedAt, key.ExpiresAt, key.Key)
		if err != nil {
			return fmt.Errorf("failed to update idempotency key: %w", err)
		}
	}

	return nil
}

func (r *IdempotencyRepository) DeleteExpired(ctx context.Context) error {
	// Check if we're in a transaction
	if tx, ok := GetTx(ctx); ok {
		_, err := tx.ExecContext(ctx, `DELETE FROM idempotency_keys WHERE expires_at <= NOW()`)
		if err != nil {
			return fmt.Errorf("failed to delete expired idempotency keys: %w", err)
		}
	} else {
		_, err := r.db.ExecContext(ctx, `DELETE FROM idempotency_keys WHERE expires_at <= NOW()`)
		if err != nil {
			return fmt.Errorf("failed to delete expired idempotency keys: %w", err)
		}
	}

	return nil
}

func (r *IdempotencyRepository) Purge(ctx context.Context, olderThan time.Time) error {
	// Check if we're in a transaction
	if tx, ok := GetTx(ctx); ok {
		_, err := tx.ExecContext(ctx, `DELETE FROM idempotency_keys WHERE created_at < $1`, olderThan)
		if err != nil {
			return fmt.Errorf("failed to purge idempotency keys: %w", err)
		}
	} else {
		_, err := r.db.ExecContext(ctx, `DELETE FROM idempotency_keys WHERE created_at < $1`, olderThan)
		if err != nil {
			return fmt.Errorf("failed to purge idempotency keys: %w", err)
		}
	}

	return nil
}
