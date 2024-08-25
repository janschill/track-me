package utils

import (
	"testing"
)

func TestOneDecimal(t *testing.T) {
	tests := []struct {
		name     string
		num      float64
		expected float64
	}{
		{"Round down", 3.14, 3.1},
		{"Round up", 2.75, 2.8},
		{"No rounding needed", 5.0, 5.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := OneDecimal(tt.num)
			if result != tt.expected {
				t.Errorf("OneDecimal(%f) = %f, want %f", tt.num, result, tt.expected)
			}
		})
	}
}

func TestInKm(t *testing.T) {
	tests := []struct {
		name     string
		distance float64
		expected float64
	}{
		{"Meters to kilometers", 1000, 1},
		{"Half kilometer", 500, 0.5},
		{"No distance", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := InKm(tt.distance)
			if result != tt.expected {
				t.Errorf("InKm(%f) = %f, want %f", tt.distance, result, tt.expected)
			}
		})
	}
}

func TestFormatTime(t *testing.T) {
	tests := []struct {
		name     string
		seconds  int64
		expected string
	}{
		{"One hour", 3600, "01:00:00"},
		{"One minute and one second", 61, "00:01:01"},
		{"All units", 3661, "01:01:01"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatTime(tt.seconds)
			if result != tt.expected {
				t.Errorf("FormatTime(%d) = %s, want %s", tt.seconds, result, tt.expected)
			}
		})
	}
}

func TestOnDay(t *testing.T) {
	tests := []struct {
		name     string
		ts       int64
		expected string
	}{
		{"New Year's Day 2020", 1577836800, "01 January"},
		{"Leap Day 2020", 1582934400, "29 February"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := OnDay(tt.ts)
			if result != tt.expected {
				t.Errorf("OnDay(%d) = %s, want %s", tt.ts, result, tt.expected)
			}
		})
	}
}

func TestWroteOnTime(t *testing.T) {
	tests := []struct {
		name     string
		ts       int64
		expected string
	}{
		{"Example timestamp", 1609459200, "on 01 January at 00:00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WroteOnTime(tt.ts)
			if result != tt.expected {
				t.Errorf("WroteOnTime(%d) = %s, want %s", tt.ts, result, tt.expected)
			}
		})
	}
}
