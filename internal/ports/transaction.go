package ports

import "context"

// TxManager handles database transactions
type TxManager interface {
	// WithTx executes the given function within a transaction
	// If the function returns an error, the transaction is rolled back
	// Otherwise, the transaction is committed
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

// TxKey is used to store transaction in context
type txKey struct{}

// TxFromContext retrieves a transaction from context
func TxFromContext(ctx context.Context) interface{} {
	return ctx.Value(txKey{})
}

// ContextWithTx returns a new context with the transaction
func ContextWithTx(ctx context.Context, tx interface{}) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}
