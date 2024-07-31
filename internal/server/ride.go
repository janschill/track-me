package server

import (
	"fmt"
	"strconv"
	"time"

	"github.com/janschill/track-me/internal/db"
)

func oneDecimal(num float64) float64 {
	numStr := fmt.Sprintf("%.1f", num)
	formattedSpeed, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0
	}
	return formattedSpeed
}

func inKm(distance float64) float64 {
	return distance / 1000
}

func movement(events []db.Event) (float64, bool) {
	eventCount := len(events)
	if eventCount < 2 {
		return 0, false
	}

	secondToLastEvent := events[eventCount-2]
	lastEvent := events[eventCount-1]
	speed := db.CalculateSpeed(secondToLastEvent, lastEvent) * 3.6 * 1000

	return oneDecimal(speed), lastEvent.Latitude != secondToLastEvent.Latitude || lastEvent.Longitude != secondToLastEvent.Longitude
}

func distance(days []db.Day, events []db.Event) int64 {
	var totalDistance int64
	for _, d := range days {
		totalDistance += int64(d.TotalDistance)
	}
	return totalDistance / 1000
}

func progress(distanceSoFar int64) float64 {
	distanceInTotal := 3891.0
	percentage := float64(distanceSoFar) / distanceInTotal * 100

	return oneDecimal(percentage)
}

func elevation(days []db.Day, events []db.Event) (int64, int64) {
	var gain, loss int64
	for _, d := range days {
		gain += int64(d.ElevationGain)
		loss += int64(d.ElevationLoss)
	}
	return gain, loss
}

func timeMoving(days []db.Day, events []db.Event) int64 {
	var time int64
	for _, d := range days {
		time += int64(d.MovingTimeInSeconds)
	}
	return time
}

func restingTime(elapsedDays int, movingTime int64) int64 {
	return int64(elapsedDays)*24*60*60 - movingTime
}

func formatTime(seconds int64) string {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	seconds = seconds % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func onDay(ts int64) string {
	t := time.Unix(ts, 0)
	return t.Format("02 January")
}
