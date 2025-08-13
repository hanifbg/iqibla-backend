package entity

import (
	"testing"
)

func TestFormatToIndonesianCurrency(t *testing.T) {
	tests := []struct {
		name     string
		amount   float64
		expected string
	}{
		{"zero", 0, "0"},
		{"single digit", 5, "5"},
		{"double digits", 42, "42"},
		{"triple digits", 123, "123"},
		{"thousand", 1000, "1.000"},
		{"ten thousand", 10000, "10.000"},
		{"hundred thousand", 100000, "100.000"},
		{"million", 1000000, "1.000.000"},
		{"with decimal", 1234.56, "1.234,56"},
		{"large number", 1234567890, "1.234.567.890"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatToIndonesianCurrency(tt.amount)
			if result != tt.expected {
				t.Errorf("FormatToIndonesianCurrency(%v) = %v, want %v", tt.amount, result, tt.expected)
			}
		})
	}
}