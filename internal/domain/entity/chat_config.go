package entity

import "time"

// ChatConfig stores configuration for each chat/group
type ChatConfig struct {
	ChatID       int64
	CurrencyName string // Configurable currency name for this chat
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
