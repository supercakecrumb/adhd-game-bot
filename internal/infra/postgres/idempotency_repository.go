package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/ports"
)

type PgIdempotencyRepository struct {
	db *sql.DB
}

func NewPgIdempotencyRepository(db *sql.DB) *PgIdempotencyRepository {
	return &PgIdempotencyRepository{db: db}
}

func (r *PgIdempotencyRepository) Create(ctx context.Context, key *entity.IdempotencyKey) error {
	query := `
		INSERT INTO idempotency_keys (key, operation, user_id, status, result, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		key.Key,
		key.Operation,
		key.UserID,
		key.Status,
		key.Result,
		key.CreatedAt,
		key.ExpiresAt,
	)

	if err != nil {
		// Check if it's a unique constraint violation
		if isUniqueViolation(err) {
			return ports.ErrIdempotencyKeyExists
		}
		return err
	}

	return nil
}

func (r *PgIdempotencyRepository) FindByKey(ctx context.Context, key string) (*entity.IdempotencyKey, error) {
	query := `
		SELECT key, operation, user_id, status, result, created_at, completed_at, expires_at
		FROM idempotency_keys
		WHERE key = $1
	`

	var idempKey entity.IdempotencyKey
	err := r.db.QueryRowContext(ctx, query, key).Scan(
		&idempKey.Key,
		&idempKey.Operation,
		&idempKey.UserID,
		&idempKey.Status,
		&idempKey.Result,
		&idempKey.CreatedAt,
		&idempKey.CompletedAt,
		&idempKey.ExpiresAt,
	)

	if err == sql.ErrNoRows {
		return nil, ports.ErrIdempotencyKeyNotFound
	}
	if err != nil {
		return nil, err
	}

	return &idempKey, nil
}

func (r *PgIdempotencyRepository) Update(ctx context.Context, key *entity.IdempotencyKey) error {
	query := `
		UPDATE idempotency_keys
		SET status = $2, result = $3, completed_at = $4
		WHERE key = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		key.Key,
		key.Status,
		key.Result,
		key.CompletedAt,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ports.ErrIdempotencyKeyNotFound
	}

	return nil
}

func (r *PgIdempotencyRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM idempotency_keys WHERE expires_at < $1`
	_, err := r.db.ExecContext(ctx, query, time.Now())
	return err
}

func isUniqueViolation(err error) bool {
	// PostgreSQL unique violation error code is 23505
	// This is a simplified check - in production you'd want to use pq.Error
	return err != nil && err.Error() == "pq: duplicate key value violates unique constraint"
}
