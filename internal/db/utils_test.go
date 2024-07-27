package db

import (
	"math"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

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
			relativeThreshold := 0.05

			if math.Abs(result-expectedDistance) > relativeThreshold*expectedDistance {
				t.Errorf("CalculateDistance = %f; expected %f (file: %s)", result, expectedDistance, path)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to walk through directory: %v", err)
	}
}
