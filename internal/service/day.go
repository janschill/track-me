package service

import (
	"cmp"
	"log"
	"math"
	"slices"
	"time"

	"github.com/janschill/track-me/internal/repository"
	"github.com/janschill/track-me/internal/utils"
)

type Day struct {
	Date                   string
	AverageSpeed           float64
	MaxSpeed               float64
	DistanceInMeters       float64
	ElevationGain          int64
	ElevationLoss          int64
	AverageAltitude        float64
	MaxAltitude            float64
	MinAltitude            float64
	MovingTimeInSeconds    int64
	NumberOfStops          int64
	TotalStopTimeInSeconds int64
}

type Ride struct {
	Distance      float64
	Progress      float64
	ElevationGain int64
	ElevationLoss int64
	MovingTime    int64
	RestingTime   int64
	ElapsedDays   int
	RemainingDays int
}

type DayService struct {
	daysCache map[string]Day
}

func NewDayService() *DayService {
	return &DayService{
		daysCache: make(map[string]Day),
	}
}

func (s *DayService) calculateDayStats(date string, events []repository.Event) Day {
	movingTime, averageSpeed, maxSpeed := utils.CalculateMovingTimeAndAverageSpeed(events, 0.001)
	gain, loss := utils.CalculateElevationGainAndLoss(events)
	averageAltitude, maxAltitude, minAltitude := utils.CalculateAltitudes(events)
	// maxSpeed := utils.CalculateMaxSpeed(events)
	// numberOfStops, stopTime := utils.CalculateStops(events)

	return Day{
		Date:                   date,
		AverageSpeed:           averageSpeed,
		MaxSpeed:               maxSpeed,
		DistanceInMeters:       utils.DistanceInMeters(events),
		ElevationGain:          gain,
		ElevationLoss:          loss,
		AverageAltitude:        float64(averageAltitude),
		MaxAltitude:            maxAltitude,
		MinAltitude:            minAltitude,
		MovingTimeInSeconds:    int64(movingTime),
		NumberOfStops:          0,
		TotalStopTimeInSeconds: 0,
	}
}

func updateRideStats(ride *Ride, day Day) {
	ride.Distance += day.DistanceInMeters
	ride.ElevationGain += day.ElevationGain
	ride.ElevationLoss += day.ElevationLoss
	ride.MovingTime += day.MovingTimeInSeconds
}

func (s *DayService) GetDays(events []repository.Event) ([]Day, Ride) {
	currentDate := time.Now().Format("2006-01-02")
	eventsByDay := make(map[string][]repository.Event)

	var ride Ride

	for _, event := range events {
		date := time.Unix(event.TimeStamp, 0).Format("2006-01-02")
		eventsByDay[date] = append(eventsByDay[date], event)
	}

	days := make([]Day, 0, len(eventsByDay))

	for date, events := range eventsByDay {
		var day Day

		if date == currentDate {
			// Do not cache the current day
			day = s.calculateDayStats(date, events)
			days = append(days, day)
		} else {
			if day, ok := s.daysCache[date]; ok {
				// Check if the cached day has the same number of events
				log.Printf("Cache hit for %v", date)
				if len(events) == len(eventsByDay[date]) {
					days = append(days, day)
					updateRideStats(&ride, day)
					continue // jump to next element in loop
				}
				// Bust the cache if new events are detected
				delete(s.daysCache, date)
			}
			day = s.calculateDayStats(date, events)
			s.daysCache[date] = day
			days = append(days, day)
		}

		updateRideStats(&ride, day)
	}

	// Sort days by date
	slices.SortFunc(days, func(a, b Day) int { return cmp.Compare(a.Date, b.Date) })

	ride.RestingTime = utils.RestingTime(len(days), ride.MovingTime)
	ride.Progress = utils.Progress(ride.Distance)
	ride.ElapsedDays = len(days)
	ride.RemainingDays = int(math.Max(0, float64(30-ride.ElapsedDays)))

	return days, ride
}
