package valueobject

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecimal_New(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"positive", "12.34", "12.34", false},
		{"negative", "-56.78", "-56.78", false},
		{"zero", "0", "0", false},
		{"invalid", "abc", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				assert.Panics(t, func() { NewDecimal(tt.input) })
				return
			}
			got := NewDecimal(tt.input)
			assert.Equal(t, tt.want, got.String())
		})
	}
}

func TestDecimal_Arithmetic(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		add  string
		sub  string
		mul  string
		div  string
	}{
		{"integers", "10", "5", "15", "5", "50", "2"},
		{"decimals", "10.5", "0.25", "10.75", "10.25", "2.625", "42"},
		{"negative", "-3.5", "1.5", "-2", "-5", "-5.25", "-2.3333333333333333"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewDecimal(tt.a)
			b := NewDecimal(tt.b)

			assert.Equal(t, tt.add, a.Add(b).String(), "Add")
			assert.Equal(t, tt.sub, a.Sub(b).String(), "Sub")
			assert.Equal(t, tt.mul, a.Mul(b).String(), "Mul")
			assert.Equal(t, tt.div, a.Div(b).String(), "Div")
		})
	}
}

func TestDecimal_Comparison(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		cmp  int
	}{
		{"equal", "10", "10", 0},
		{"less", "5", "10", -1},
		{"greater", "10", "5", 1},
		{"decimal equal", "10.50", "10.5", 0},
		{"decimal less", "10.49", "10.50", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewDecimal(tt.a)
			b := NewDecimal(tt.b)

			assert.Equal(t, tt.cmp, a.Cmp(b))
		})
	}
}

func TestDecimal_Sign(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		positive bool
		negative bool
		zero     bool
	}{
		{"positive", "10.50", true, false, false},
		{"negative", "-5.25", false, true, false},
		{"zero", "0", false, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDecimal(tt.input)
			assert.Equal(t, tt.positive, d.IsPositive(), "IsPositive")
			assert.Equal(t, tt.negative, d.IsNegative(), "IsNegative")
			assert.Equal(t, tt.zero, d.IsZero(), "IsZero")
		})
	}
}
