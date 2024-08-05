package service

import (
	"github.com/janschill/track-me/internal/repository"
	"github.com/janschill/track-me/internal/utils"
)

type Ride struct {
	IsMoving      bool
	LastPing      int64
	Distance      int64
	Progress      float64
	ElevationGain int64
	ElevationLoss int64
	MovingTime    string
	RestingTime   string
	ElapsedDays   int
	RemainingDays int
	CurrentSpeed  float64
}

func NewRide(events []repository.Event, days []repository.Day) Ride {
	var (
		lastEvent            repository.Event
		currentSpeed         float64
		isMoving             bool
		dist                 int64
		gain                 int64
		loss                 int64
		movingTime           int64
		movingTimeFormatted  string
		restingTimeFormatted string
	)

	if len(events) > 1 {
		lastEvent = events[len(events)-1]
		currentSpeed, isMoving = utils.Movement(events)
		dist = utils.Distance(days, events)
		gain, loss = utils.Elevation(days, events)
		movingTime = utils.TimeMoving(days, events)
		movingTimeFormatted = utils.FormatTime(movingTime)
		restingTimeFormatted = utils.FormatTime(utils.RestingTime(len(days), movingTime))
	}

	return Ride{
		IsMoving:      isMoving,
		LastPing:      lastEvent.TimeStamp,
		Distance:      dist,
		Progress:      utils.Progress(dist),
		CurrentSpeed:  currentSpeed,
		ElevationGain: gain,
		ElevationLoss: loss,
		MovingTime:    movingTimeFormatted,
		RestingTime:   restingTimeFormatted,
		ElapsedDays:   len(days),
		RemainingDays: 30 - len(days),
	}
}
