package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/supercakecrumb/adhd-game-bot/internal/ports"
)

// TxManager implements transaction management for PostgreSQL
type TxManager struct {
	db *sql.DB
}

// NewTxManager creates a new PostgreSQL transaction manager
func NewTxManager(db *sql.DB) *TxManager {
	return &TxManager{db: db}
}

// WithTx executes the given function within a transaction
func (m *TxManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	// Check if we're already in a transaction
	if tx := ports.TxFromContext(ctx); tx != nil {
		// Already in a transaction, just execute the function
		return fn(ctx)
	}

	// Start a new transaction
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Create a new context with the transaction
	txCtx := ports.ContextWithTx(ctx, tx)

	// Execute the function
	err = fn(txCtx)
	if err != nil {
		// Rollback on error
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("failed to rollback transaction: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetTx retrieves the transaction from context
func GetTx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ports.TxFromContext(ctx).(*sql.Tx)
	return tx, ok
}
