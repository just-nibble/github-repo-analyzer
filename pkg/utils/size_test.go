// pkg/utils/size_test.go
package utils

import (
	"testing"
)

func TestBytesToHumanReadable(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int64
		expected string
	}{
		{"Bytes", 500, "500 B"},
		{"Kilobytes", 1024, "1.00 KB"},
		{"Megabytes", 1024 * 1024, "1.00 MB"},
		{"Gigabytes", 1024 * 1024 * 1024, "1.00 GB"},
		{"Terabytes", 1024 * 1024 * 1024 * 1024, "1.00 TB"},
		{"Zero", 0, "0 B"},
		{"Large Number", 1024 * 1024 * 1024 * 1024 * 1024, "1.00 PB"},
		{"Partial Units", 1536, "1.50 KB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BytesToHumanReadable(tt.bytes)
			if result != tt.expected {
				t.Errorf("BytesToHumanReadable(%d) = %s; want %s", tt.bytes, result, tt.expected)
			}
		})
	}
}

func TestBytesToMB(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int64
		expected float64
	}{
		{"One MB", 1024 * 1024, 1.0},
		{"Two MB", 2 * 1024 * 1024, 2.0},
		{"Half MB", 512 * 1024, 0.5},
		{"Zero", 0, 0.0},
		{"Small Bytes", 1024, 0.0009765625},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BytesToMB(tt.bytes)
			if result != tt.expected {
				t.Errorf("BytesToMB(%d) = %f; want %f", tt.bytes, result, tt.expected)
			}
		})
	}
}

func TestRoundToTwoDecimals(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{"Integer", 1.0, 1.0},
		{"One Decimal", 1.5, 1.50},
		{"Two Decimals", 1.25, 1.25},
		{"Round Down", 1.249, 1.25},
		{"Round Up", 1.251, 1.25},
		{"Zero", 0.0, 0.0},
		{"Negative", -1.23456, -1.23},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RoundToTwoDecimals(tt.input)
			if result != tt.expected {
				t.Errorf("RoundToTwoDecimals(%f) = %f; want %f", tt.input, result, tt.expected)
			}
		})
	}
}
