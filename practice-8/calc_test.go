package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSubtract(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"Both positive", 10, 5, 5},
		{"Positive minus zero", 5, 0, 5},
		{"Negative minus positive", -1, 4, -5},
		{"Both negative", -2, -3, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Subtract(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Subtract(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		got, err := Divide(10, 2)
		if err != nil || got != 5 {
			t.Errorf("Divide(10, 2) = %d; want 5", got)
		}
	})

	t.Run("Division by zero", func(t *testing.T) {
		_, err := Divide(10, 0)
		if err == nil {
			t.Error("Expected error for division by zero")
		}
	})
}

func TestAddTableDriven(t *testing.T) {
	cases := map[string]struct {
		val1, val2 int
		expected   int
	}{
		"both_positive":      {2, 3, 5},
		"positive_plus_zero": {5, 0, 5},
		"negative_plus_pos":  {-1, 4, 3},
		"both_negative":      {-2, -3, -5},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			result := Add(tc.val1, tc.val2)
			require.Equal(t, tc.expected, result)
		})
	}
}