package utils

import (
	"math"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/janschill/track-me/internal/repository"
)

func walkTestFiles(t *testing.T, fileHandler func(path string, data GPXData)) {
	dirPath := "./test_data"
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			data, err := ReadGPXDataFromFile(path)
			if err != nil {
				t.Fatalf("Failed to read GPX data from file %s: %v", path, err)
			}
			fileHandler(path, data)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to walk through directory: %v", err)
	}
}

func TestDistanceInMeters(t *testing.T) {
	walkTestFiles(t, func(path string, data GPXData) {
		expectedDistance, err := strconv.ParseFloat(data.Distance, 64)
		if err != nil {
			t.Fatalf("Failed to parse expected distance from file %s: %v", path, err)
		}

		events := ConvertPointsToEvents(data.Points)
		result := DistanceInMeters(events)
		relativeThreshold := 0.05

		if math.Abs(result-expectedDistance) > relativeThreshold*expectedDistance {
			t.Errorf("DistanceInMeters = %f; expected %f (file: %s)", result, expectedDistance, path)
		}
	})
}

func TestCalculateMovingTimeAndAverageSpeed(t *testing.T) {
	walkTestFiles(t, func(path string, data GPXData) {
		expectedMovingTime, err := strconv.ParseFloat(data.MovingTime, 64)
		if err != nil {
			t.Fatalf("Failed to parse expected moving time from file %s: %v", path, err)
		}

		expectedAverageSpeed, err := strconv.ParseFloat(data.AverageSpeed, 64)
		if err != nil {
			t.Fatalf("Failed to parse expected average speed from file %s: %v", path, err)
		}

		events := ConvertPointsToEvents(data.Points)
		movingTime, averageSpeed, _ := CalculateMovingTimeAndAverageSpeed(events, 0.001)
		relativeThreshold := 0.05

		if math.Abs(movingTime-expectedMovingTime) > relativeThreshold*expectedMovingTime {
			t.Errorf("CalculateMovingTimeAndAverageSpeed moving time = %f; expected %f (file: %s)", movingTime, expectedMovingTime, path)
		}

		if math.Abs(averageSpeed-expectedAverageSpeed) > relativeThreshold*expectedAverageSpeed {
			t.Errorf("CalculateMovingTimeAndAverageSpeed average speed = %f; expected %f (file: %s)", averageSpeed, expectedAverageSpeed, path)
		}
	})
}

func TestCalculateAltitudes(t *testing.T) {
	tests := []struct {
		name           string
		events         []repository.Event
		expectedAvgAlt float64
		expectedMaxAlt float64
		expectedMinAlt float64
	}{
		{
			name: "multiple events",
			events: []repository.Event{
				{Altitude: 100},
				{Altitude: 200},
				{Altitude: 150},
			},
			expectedAvgAlt: 150,
			expectedMaxAlt: 200,
			expectedMinAlt: 100,
		},
		{
			name: "single event",
			events: []repository.Event{
				{Altitude: 100},
			},
			expectedAvgAlt: 100,
			expectedMaxAlt: 100,
			expectedMinAlt: 100,
		},
		{
			name:           "no events",
			events:         []repository.Event{},
			expectedAvgAlt: 0.0,
			expectedMaxAlt: 0.0,
			expectedMinAlt: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			avgAlt, maxAlt, minAlt := CalculateAltitudes(tt.events)
			if avgAlt != tt.expectedAvgAlt {
				t.Errorf("CalculateAltitudes() averageAltitude = %f; expected %f", avgAlt, tt.expectedAvgAlt)
			}
			if maxAlt != tt.expectedMaxAlt {
				t.Errorf("CalculateAltitudes() maxAltitude = %f; expected %f", maxAlt, tt.expectedMaxAlt)
			}
			if minAlt != tt.expectedMinAlt {
				t.Errorf("CalculateAltitudes() minAltitude = %f; expected %f", minAlt, tt.expectedMinAlt)
			}
		})
	}
}

func TestCalculateElevationGainAndLoss(t *testing.T) {
	tests := []struct {
		name         string
		events       []repository.Event
		expectedGain int64
		expectedLoss int64
	}{
		{
			name: "multiple events with gain and loss",
			events: []repository.Event{
				{Altitude: 100},
				{Altitude: 200},
				{Altitude: 150},
				{Altitude: 250},
			},
			expectedGain: 200,
			expectedLoss: 50,
		},
		{
			name: "single event",
			events: []repository.Event{
				{Altitude: 100},
			},
			expectedGain: 0,
			expectedLoss: 0,
		},
		{
			name:         "no events",
			events:       []repository.Event{},
			expectedGain: 0,
			expectedLoss: 0,
		},
		{
			name: "all gain",
			events: []repository.Event{
				{Altitude: 100},
				{Altitude: 200},
				{Altitude: 300},
			},
			expectedGain: 200,
			expectedLoss: 0,
		},
		{
			name: "all loss",
			events: []repository.Event{
				{Altitude: 300},
				{Altitude: 200},
				{Altitude: 100},
			},
			expectedGain: 0,
			expectedLoss: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gain, loss := CalculateElevationGainAndLoss(tt.events)
			if gain != tt.expectedGain {
				t.Errorf("CalculateElevationGainAndLoss() elevationGain = %d; expected %d", gain, tt.expectedGain)
			}
			if loss != tt.expectedLoss {
				t.Errorf("CalculateElevationGainAndLoss() elevationLoss = %d; expected %d", loss, tt.expectedLoss)
			}
		})
	}
}
