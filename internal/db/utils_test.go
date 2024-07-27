package db

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

type TestPoint struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Altitude  float64 `json:"altitude"`
}

type GPXData struct {
	Distance string      `json:"distance"`
	Points   []TestPoint `json:"points"`
}

func readGPXDataFromFile(path string) (GPXData, error) {
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

func TestCalculateDistance(t *testing.T) {
	dirPath := "./test_data"
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			data, err := readGPXDataFromFile(path)
			if err != nil {
				t.Fatalf("Failed to read GPX data from file %s: %v", path, err)
			}

			expectedDistance, err := strconv.ParseFloat(data.Distance, 64)
			if err != nil {
				t.Fatalf("Failed to parse expected distance from file %s: %v", path, err)
			}

			events := convertPointsToEvents(data.Points)
			result := CalculateDistance(events)
			threshold := 1.5 // Allowable error margin in km

			if math.Abs(result-expectedDistance) > threshold {
				t.Errorf("CalculateDistance = %f; expected %f (file: %s)", result, expectedDistance, path)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to walk through directory: %v", err)
	}
}

func convertPointsToEvents(points []TestPoint) []Event {
	events := make([]Event, len(points))
	for i, point := range points {
		events[i] = Event{
			Longitude: point.Longitude,
			Latitude:  point.Latitude,
			Altitude:  int(point.Altitude),
		}
	}
	return events
}
