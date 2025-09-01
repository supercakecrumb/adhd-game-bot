package valueobject

import (
	"github.com/shopspring/decimal"
)

// Decimal represents a fixed-point decimal number with high precision.
// It wraps shopspring/decimal to provide domain-specific functionality.
type Decimal struct {
	value decimal.Decimal
}

// NewDecimal creates a new Decimal from a string representation.
// Panics if the string is not a valid decimal number.
func NewDecimal(s string) Decimal {
	d, err := decimal.NewFromString(s)
	if err != nil {
		panic(err) // Should be caught by validation
	}
	return Decimal{value: d}
}

// Add returns a new Decimal that is the sum of d and other.
func (d Decimal) Add(other Decimal) Decimal {
	return Decimal{value: d.value.Add(other.value)}
}

// Sub returns a new Decimal that is the difference of d and other.
func (d Decimal) Sub(other Decimal) Decimal {
	return Decimal{value: d.value.Sub(other.value)}
}

// Mul returns a new Decimal that is the product of d and other.
func (d Decimal) Mul(other Decimal) Decimal {
	return Decimal{value: d.value.Mul(other.value)}
}

// Div returns a new Decimal that is the quotient of d and other.
// Uses 16 decimal places precision for division operations.
func (d Decimal) Div(other Decimal) Decimal {
	return Decimal{value: d.value.Div(other.value).Round(16)}
}

// Cmp compares d and other and returns:
// -1 if d < other
// 0 if d == other
// 1 if d > other
func (d Decimal) Cmp(other Decimal) int {
	return d.value.Cmp(other.value)
}

// IsPositive returns true if d is greater than zero.
func (d Decimal) IsPositive() bool {
	return d.value.IsPositive()
}

// IsNegative returns true if d is less than zero.
func (d Decimal) IsNegative() bool {
	return d.value.IsNegative()
}

// IsZero returns true if d is zero.
func (d Decimal) IsZero() bool {
	return d.value.IsZero()
}

// String returns the string representation of the decimal.
func (d Decimal) String() string {
	return d.value.String()
}

// Float64 returns the float64 representation of the decimal.
// Note: This may lose precision.
func (d Decimal) Float64() float64 {
	f, _ := d.value.Float64()
	return f
}

// MarshalJSON implements json.Marshaler.
func (d Decimal) MarshalJSON() ([]byte, error) {
	return d.value.MarshalJSON()
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *Decimal) UnmarshalJSON(data []byte) error {
	return d.value.UnmarshalJSON(data)
}
