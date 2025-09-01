package valueobject_test

import (
	"fmt"
	"math"
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/require"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
)

// Helper function to create decimal from int64
func newDecimalFromInt64(n int64) valueobject.Decimal {
	return valueobject.NewDecimal(fmt.Sprintf("%d", n))
}

// Property: Addition is commutative (a + b = b + a)
func TestDecimal_AddCommutative(t *testing.T) {
	f := func(a, b int64) bool {
		d1 := newDecimalFromInt64(a)
		d2 := newDecimalFromInt64(b)

		result1 := d1.Add(d2)
		result2 := d2.Add(d1)

		return result1.Cmp(result2) == 0
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

// Property: Addition is associative ((a + b) + c = a + (b + c))
func TestDecimal_AddAssociative(t *testing.T) {
	f := func(a, b, c int64) bool {
		d1 := newDecimalFromInt64(a)
		d2 := newDecimalFromInt64(b)
		d3 := newDecimalFromInt64(c)

		result1 := d1.Add(d2).Add(d3)
		result2 := d1.Add(d2.Add(d3))

		return result1.Cmp(result2) == 0
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

// Property: Zero is the identity element for addition (a + 0 = a)
func TestDecimal_AddIdentity(t *testing.T) {
	f := func(a int64) bool {
		d := newDecimalFromInt64(a)
		zero := valueobject.NewDecimal("0")

		result := d.Add(zero)

		return result.Cmp(d) == 0
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

// Property: Subtraction is the inverse of addition (a + b - b = a)
func TestDecimal_SubInverse(t *testing.T) {
	f := func(a, b int64) bool {
		d1 := newDecimalFromInt64(a)
		d2 := newDecimalFromInt64(b)

		result := d1.Add(d2).Sub(d2)

		return result.Cmp(d1) == 0
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

// Property: Multiplication is commutative (a * b = b * a)
func TestDecimal_MulCommutative(t *testing.T) {
	f := func(a, b int32) bool {
		// Use smaller values to avoid overflow
		d1 := newDecimalFromInt64(int64(a))
		d2 := newDecimalFromInt64(int64(b))

		result1 := d1.Mul(d2)
		result2 := d2.Mul(d1)

		return result1.Cmp(result2) == 0
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

// Property: One is the identity element for multiplication (a * 1 = a)
func TestDecimal_MulIdentity(t *testing.T) {
	f := func(a int64) bool {
		d := newDecimalFromInt64(a)
		one := valueobject.NewDecimal("1")

		result := d.Mul(one)

		return result.Cmp(d) == 0
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

// Property: Comparison is transitive (if a < b and b < c, then a < c)
func TestDecimal_CompareTransitive(t *testing.T) {
	f := func(a, b, c int64) bool {
		if a >= b || b >= c {
			return true // Skip non-ordered inputs
		}

		d1 := newDecimalFromInt64(a)
		d2 := newDecimalFromInt64(b)
		d3 := newDecimalFromInt64(c)

		return d1.Cmp(d2) < 0 && d2.Cmp(d3) < 0 && d1.Cmp(d3) < 0
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

// Property: Zero comparison
func TestDecimal_ZeroComparison(t *testing.T) {
	f := func(a int64) bool {
		d := newDecimalFromInt64(a)
		zero := valueobject.NewDecimal("0")

		if a > 0 {
			return d.IsPositive() && d.Cmp(zero) > 0
		} else if a < 0 {
			return d.IsNegative() && d.Cmp(zero) < 0
		} else {
			return d.IsZero() && d.Cmp(zero) == 0
		}
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

// Property: String parsing and formatting are inverse operations
func TestDecimal_StringRoundTrip(t *testing.T) {
	f := func(a int64) bool {
		if a == math.MinInt64 {
			return true // Skip edge case that might overflow
		}

		d1 := newDecimalFromInt64(a)
		str := d1.String()
		d2 := valueobject.NewDecimal(str)

		return d1.Cmp(d2) == 0
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

// Edge case tests
func TestDecimal_EdgeCases(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "MaxInt64 operations",
			test: func(t *testing.T) {
				max := newDecimalFromInt64(math.MaxInt64)
				one := valueobject.NewDecimal("1")

				// Adding 1 to MaxInt64 should not panic
				require.NotPanics(t, func() {
					_ = max.Add(one)
				})
			},
		},
		{
			name: "MinInt64 operations",
			test: func(t *testing.T) {
				min := newDecimalFromInt64(math.MinInt64)
				one := valueobject.NewDecimal("1")

				// Subtracting 1 from MinInt64 should not panic
				require.NotPanics(t, func() {
					_ = min.Sub(one)
				})
			},
		},
		{
			name: "Division by zero",
			test: func(t *testing.T) {
				d := valueobject.NewDecimal("100")
				zero := valueobject.NewDecimal("0")

				// Division by zero should panic (as per shopspring/decimal behavior)
				require.Panics(t, func() {
					_ = d.Div(zero)
				})
			},
		},
		{
			name: "Very small decimal precision",
			test: func(t *testing.T) {
				// Test with very small decimals
				small := valueobject.NewDecimal("0.00000001")

				// Operations should maintain precision
				result := small.Add(small)
				expected := valueobject.NewDecimal("0.00000002")
				require.Equal(t, 0, result.Cmp(expected))
			},
		},
		{
			name: "Large decimal operations",
			test: func(t *testing.T) {
				// Test with large numbers
				large1 := valueobject.NewDecimal("999999999999999999.999999999999999999")
				large2 := valueobject.NewDecimal("1.000000000000000001")

				// Addition should work correctly
				result := large1.Add(large2)
				expected := valueobject.NewDecimal("1000000000000000001")

				// Due to rounding, we check if they're very close
				diff := result.Sub(expected)
				require.True(t, diff.Cmp(valueobject.NewDecimal("0.000000000000001")) < 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
