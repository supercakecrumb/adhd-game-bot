package inmemory

import (
	"context"
)

// TxManager is a no-op transaction manager for in-memory repositories
type TxManager struct{}

// NewTxManager creates a new in-memory transaction manager
func NewTxManager() *TxManager {
	return &TxManager{}
}

// WithTx executes the given function without any actual transaction
// This is suitable for in-memory repositories that don't need transactions
func (m *TxManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	// Simply execute the function without any transaction
	return fn(ctx)
}
