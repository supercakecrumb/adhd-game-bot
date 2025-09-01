package ports

import "errors"

// Repository errors
var (
	ErrUserAlreadyExists      = errors.New("user already exists")
	ErrUserNotFound           = errors.New("user not found")
	ErrTaskNotFound           = errors.New("task not found")
	ErrTimerNotFound          = errors.New("timer not found")
	ErrChatConfigNotFound     = errors.New("chat config not found")
	ErrShopItemNotFound       = errors.New("shop item not found")
	ErrPurchaseNotFound       = errors.New("purchase not found")
	ErrInsufficientStock      = errors.New("insufficient stock")
	ErrInsufficientFunds      = errors.New("insufficient funds")
	ErrIdempotencyKeyExists   = errors.New("idempotency key already exists")
	ErrIdempotencyKeyNotFound = errors.New("idempotency key not found")
)
