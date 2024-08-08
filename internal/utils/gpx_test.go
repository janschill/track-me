package utils

import (
	"os"
	"reflect"
	"testing"

	"github.com/janschill/track-me/internal/repository"
)

func TestReadGPXDataFromFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "example.*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	exampleData := `{"distance":"10","movingTime":"1h","averageSpeed":"10 km/h","elevationGain":"100m","elevationLoss":"50m","points":[{"longitude":10.0,"latitude":20.0,"altitude":30.0,"timestamp":1234567890}]}`
	if _, err := tmpFile.WriteString(exampleData); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	data, err := ReadGPXDataFromFile(tmpFile.Name())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if data.Distance != "10" {
		t.Errorf("Expected distance to be '10', got '%s'", data.Distance)
	}

	_, err = ReadGPXDataFromFile("nonexistent.json")
	if err == nil {
		t.Errorf("Expected an error for nonexistent file, got nil")
	}
}

func TestConvertPointsToEvents(t *testing.T) {
	points := []TestPoint{
		{Longitude: 10.0, Latitude: 20.0, Altitude: 30.0, TimeStamp: 1234567890},
	}
	expectedEvents := []repository.Event{
		{Longitude: 10.0, Latitude: 20.0, Altitude: 30.0, TimeStamp: 1234567890},
	}

	events := ConvertPointsToEvents(points)

	if !reflect.DeepEqual(events, expectedEvents) {
		t.Errorf("Expected events to be %+v, got %+v", expectedEvents, events)
	}
}
