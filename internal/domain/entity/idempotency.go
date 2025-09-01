package entity

import (
	"time"
)

// IdempotencyKey represents a unique key for ensuring idempotent operations
type IdempotencyKey struct {
	Key         string // Unique key for the operation
	Operation   string // Operation type (e.g., "task_complete", "purchase_item")
	UserID      int64  // User who initiated the operation
	Status      string // "pending", "completed", "failed"
	Result      string // JSON result of the operation
	CreatedAt   time.Time
	CompletedAt *time.Time
	ExpiresAt   time.Time // Keys expire after a certain time
}

// IsExpired checks if the idempotency key has expired
func (k *IdempotencyKey) IsExpired() bool {
	return time.Now().After(k.ExpiresAt)
}

// IsCompleted checks if the operation has been completed
func (k *IdempotencyKey) IsCompleted() bool {
	return k.Status == "completed"
}
