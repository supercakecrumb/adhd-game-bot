package entity

import (
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
)

// Currency represents a currency within a specific chat/group
type Currency struct {
	ID             int64                          // Auto-generated ID
	ChatID         int64                          // Telegram chat ID this currency belongs to
	Code           string                         // Currency code (e.g., "MM", "BS")
	Name           string                         // Full name (e.g., "Motivation Minutes")
	Decimals       int                            // Number of decimal places
	IsBaseCurrency bool                           // Whether this is the base currency for the chat
	ExchangeRates  map[string]valueobject.Decimal // Exchange rates to other currencies
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// ConvertTo converts an amount from this currency to another currency
func (c *Currency) ConvertTo(amount valueobject.Decimal, targetCurrencyCode string) (valueobject.Decimal, error) {
	if c.Code == targetCurrencyCode {
		return amount, nil
	}

	rate, exists := c.ExchangeRates[targetCurrencyCode]
	if !exists {
		return valueobject.Decimal{}, ErrExchangeRateNotFound
	}

	return amount.Mul(rate), nil
}
