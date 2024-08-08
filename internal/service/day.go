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
	LastPing      int64
	Distance      float64
	Progress      float64
	ElevationGain int64
	ElevationLoss int64
	MovingTime    string
	RestingTime   string
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
	movingTime, averageSpeed := utils.CalculateMovingTimeAndAverageSpeed(events, 0.001)
	gain, loss := utils.CalculateElevationGainAndLoss(events)
	averageAltitude, maxAltitude, minAltitude := utils.CalculateAltitudes(events)
	maxSpeed := utils.CalculateMaxSpeed(events)

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

func (s *DayService) GetDays(events []repository.Event) ([]Day, Ride) {
	currentDate := time.Now().Format("2006-01-02")
	eventsByDay := make(map[string][]repository.Event)

	var totalDistance float64
	var totalElevationGain, totalElevationLoss, totalMovingTime, totalNumberOfStops, totalStopTimeInSeconds int64
	var totalAverageSpeed, totalAverageAltitude float64
	var maxAltitude, minAltitude float64

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
					continue
				}
				// Bust the cache if new events are detected
				delete(s.daysCache, date)
			}
			day = s.calculateDayStats(date, events)
			s.daysCache[date] = day
			days = append(days, day)
		}

		totalDistance += day.DistanceInMeters
		totalElevationGain += day.ElevationGain
		totalElevationLoss += day.ElevationLoss
		totalMovingTime += day.MovingTimeInSeconds
		totalNumberOfStops += day.NumberOfStops
		totalStopTimeInSeconds += day.TotalStopTimeInSeconds
		totalAverageSpeed += day.AverageSpeed
		totalAverageAltitude += day.AverageAltitude
		if day.MaxAltitude > maxAltitude {
			maxAltitude = day.MaxAltitude
		}
		if minAltitude == 0 || day.MinAltitude < minAltitude {
			minAltitude = day.MinAltitude
		}
	}

	movingTimeFormatted := utils.FormatTime(totalMovingTime)
	restingTimeFormatted := utils.FormatTime(utils.RestingTime(len(days), totalMovingTime))

	slices.SortFunc(days, func(a, b Day) int { return cmp.Compare(a.Date, b.Date) })

	return days, Ride{
		LastPing:      0,
		Distance:      totalDistance,
		Progress:      utils.Progress(totalDistance),
		ElevationGain: totalElevationGain,
		ElevationLoss: totalElevationLoss,
		MovingTime:    movingTimeFormatted,
		RestingTime:   restingTimeFormatted,
		ElapsedDays:   len(days),
		RemainingDays: int(math.Max(0, float64(30-len(days)))), // just making sure we are not going under 0
	}
}
