package entity

import (
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
)

// ShopItem represents an item that can be purchased in the shop
type ShopItem struct {
	ID             int64
	ChatID         int64  // 0 for global items, specific chat ID for chat-specific items
	Code           string // Unique code for the item
	Name           string
	Description    string
	Price          valueobject.Decimal
	Category       string // e.g., "rewards", "boosts", "cosmetics"
	IsActive       bool
	Stock          *int   // nil for unlimited stock, number for limited stock
	DiscountTierID *int64 // Optional discount tier
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Purchase represents a user's purchase of a shop item
type Purchase struct {
	ID             int64
	UserID         int64
	ItemID         int64
	ItemName       string              // Denormalized for history
	ItemPrice      valueobject.Decimal // Price at time of purchase
	Quantity       int
	TotalCost      valueobject.Decimal
	Status         string // "pending", "completed", "refunded"
	DiscountTierID *int64 // Discount tier applied (if any)
	PurchasedAt    time.Time
}
