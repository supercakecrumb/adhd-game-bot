package entity

import "time"

// DiscountTier represents a pricing tier that can be applied to shop items
type DiscountTier struct {
	ID              int64
	Name            string
	Description     string
	DiscountPercent float64 // Percentage discount (e.g. 10.0 for 10%)
	MinPurchases    int     // Minimum purchases required to qualify
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
