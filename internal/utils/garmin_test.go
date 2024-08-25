package utils

import (
	"slices"
	"testing"

	"github.com/janschill/track-me/internal/repository"
)

func TestMakeRange(t *testing.T) {
	tests := []struct {
		min, max int
		want     []int
	}{
		{1, 5, []int{1, 2, 3, 4, 5}},
		{0, 2, []int{0, 1, 2}},
	}

	for _, tt := range tests {
		got := makeRange(tt.min, tt.max)
		if !slices.Equal(got, tt.want) {
			t.Errorf("makeRange(%d, %d) = %v, want %v", tt.min, tt.max, got, tt.want)
		}
	}
}

func TestHasMessage(t *testing.T) {
	tests := []struct {
		code int
		want bool
	}{
		{3, true},
		{50, true},
		{1, false},
	}

	for _, tt := range tests {
		event := repository.Event{MessageCode: tt.code}
		got := HasMessage(event)
		if got != tt.want {
			t.Errorf("hasMessage(%d) = %v, want %v", tt.code, got, tt.want)
		}
	}
}
