package ports

import (
	"context"
	"time"
)

// Clock provides time-related functionality that can be mocked for testing
type Clock interface {
	Now() time.Time
	Sleep(d time.Duration)
	After(d time.Duration) <-chan time.Time
	NewTicker(d time.Duration) *time.Ticker
}

// ContextProvider provides context management functionality
type ContextProvider interface {
	WithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc)
	WithCancel(parent context.Context) (context.Context, context.CancelFunc)
}
