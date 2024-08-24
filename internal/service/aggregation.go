package service

import (
	"github.com/janschill/track-me/internal/repository"
)

type AggregationService struct {
	repo *repository.Repository
}

func NewAggregationService(repo *repository.Repository) *AggregationService {
	return &AggregationService{repo: repo}
}

func (s *AggregationService) Aggregate(day string) {
	// events, err := s.repo.Events.AllByDay(day)
	// if err != nil {
	// 	log.Fatal("Failed to get events by day:", err)
	// 	return
	// }

	// Perform your aggregation logic here
	// Example: log the number of events
	// log.Printf("Aggregated %d events for day %s", len(events), day)
	// movingTime, averageSpeed := utils.CalculateMovingTimeAndAverageSpeed(events, 0.001)
	// gain, loss := utils.CalculateElevationGainAndLoss(events)
	// averageAltitude, maxAltitude, minAltitude := utils.CalculateAltitudes(events)
	// stops, stopTime := CalculateStops(events)

	// log.Printf("Moving time: %v", movingTime)
	// log.Printf("Average Speed: %v", averageSpeed)
	// log.Printf("Gain: %v", gain)
	// log.Printf("Loss: %v", loss)
	// log.Printf("MaxSpeed: %v", maxSpeed)
	// log.Printf("Average Alt: %v", averageAltitude)
	// log.Printf("Max Alt: %v", maxAltitude)
	// log.Printf("Min Alt: %v", minAltitude)
	// log.Printf("Stops: %v", stops)
	// log.Printf("Stop time: %v", stopTime)

	// log.Printf("Reducing events from %v", len(events))
	// reduced_events := utils.Rdp(events, 0.0002) // roughly 1500 -> 321
	// log.Printf("Number of reduced_events after Rdp %v", len(reduced_events))

	// distance := CalculateDistance(events)
	// log.Printf("Distance before aggregate events: %v", distance)
	// log.Printf("Reducing events from %v", len(events))
	// reduced_events := Rdp(copyEvents(events), 0.0002) // roughly 1500 -> 321
	// log.Printf("Number of reduced_events after Rdp %v", len(reduced_events))
	// log.Printf("Number of events after Rdp %v", len(events))
	// distance = CalculateDistance(reduced_events)
	// log.Printf("Distance on aggregated events %v", distance)
	// distance = CalculateDistance(events)
	// log.Printf("Distance all after aggregate: %v", distance)

	// reduced_events_json, err := json.Marshal(reduced_events)
	// if err != nil {
	// 	log.Fatal("Failed to marshal points data:", err)
	// 	return
	// }

	// d := .Day{
	// 	TimeStamp:              events[0].TimeStamp,
	// 	Points:                 string(reduced_events_json),
	// 	TripID:                 reduced_events[0].TripID,
	// 	AverageSpeed:           averageSpeed,
	// 	MaxSpeed:               0,
	// 	MinSpeed:               0,
	// 	TotalDistance:          utils.CalculateDistance(reduced_events),
	// 	ElevationGain:          gain,
	// 	ElevationLoss:          loss,
	// 	AverageAltitude:        float64(averageAltitude),
	// 	MaxAltitude:            maxAltitude,
	// 	MinAltitude:            minAltitude,
	// 	MovingTimeInSeconds:    int64(movingTime),
	// 	NumberOfStops:          0,
	// 	TotalStopTimeInSeconds: 0,
	// }

	// err = s.repo.EventsCache.Create(d)
	// if err != nil {
	// 	log.Fatal("Failed to save day")
	// }
}
