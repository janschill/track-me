package utils

import (
	"encoding/json"
	"os"

	"github.com/janschill/track-me/internal/repository"
)

type TestPoint struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Altitude  float64 `json:"altitude"`
	TimeStamp int64   `json:"timestamp"`
}

type GPXData struct {
	Distance      string      `json:"distance"`
	MovingTime    string      `json:"movingTime"`
	AverageSpeed  string      `json:"averageSpeed"`
	ElevationGain string      `json:"elevationGain"`
	ElevationLoss string      `json:"elevationLoss"`
	Points        []TestPoint `json:"points"`
}

func ReadGPXDataFromFile(path string) (GPXData, error) {
	var data GPXData
	file, err := os.Open(path)
	if err != nil {
		return data, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	return data, err
}

func ConvertPointsToEvents(points []TestPoint) []repository.Event {
	events := make([]repository.Event, len(points))
	for i, point := range points {
		events[i] = repository.Event{
			Longitude: point.Longitude,
			Latitude:  point.Latitude,
			Altitude:  point.Altitude,
			TimeStamp: point.TimeStamp,
		}
	}
	return events
}
