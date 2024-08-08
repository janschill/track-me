package utils

import (
	"slices"

	"github.com/janschill/track-me/internal/repository"
)


func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func HasMessage(event repository.Event) bool {
	codes := []int{3, 14, 15, 16, 66, 67}
	codes = append(codes, makeRange(24, 63)...)
	return slices.Contains(codes, event.MessageCode)
}
