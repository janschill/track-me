package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/janschill/track-me/internal/repository"
)

func OneDecimal(num float64) float64 {
	numStr := fmt.Sprintf("%.1f", num)
	formattedSpeed, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0
	}
	return formattedSpeed
}

func InKm(distance float64) float64 {
	return distance / 1000
}

func FormatTime(seconds int64) string {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	seconds = seconds % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func OnDay(ts int64) string {
	t := time.Unix(ts, 0).UTC()
	return t.Format("02 January")
}

func WroteOnTime(ts int64) string {
	t := time.Unix(ts, 0).UTC()
	return t.Format("on 02 January at 15:04")
}

func OnDayFromString(dateStr string) (string, error) {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", err
	}

	return t.Format("02 January"), nil
}

func FindKudos(kudos []repository.Kudos, day string) int {
	for _, k := range kudos {
		if k.Day == day {
			return k.Count
		}
	}
	return 0
}
