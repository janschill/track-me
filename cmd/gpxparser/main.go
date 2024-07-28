package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"time"
)

type GPX struct {
	XMLName xml.Name `xml:"gpx"`
	Trk     Trk      `xml:"trk"`
}

type Trk struct {
	TrkSeg TrkSeg `xml:"trkseg"`
}

type TrkSeg struct {
	TrkPts []TrkPt `xml:"trkpt"`
}

type TrkPt struct {
	Lat  float64 `xml:"lat,attr"`
	Lon  float64 `xml:"lon,attr"`
	Ele  float64 `xml:"ele"`
	Time string  `xml:"time"`
}

type Point struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Altitude  float64 `json:"altitude"`
	TimeStamp int64   `json:"timestamp"`
}

type Output struct {
	Distance      string  `json:"distance"`
	MovingTime    string  `json:"movingTime"`
	AverageSpeed  string  `json:"averageSpeed"`
	ElevationGain string  `json:"elevationGain"`
	ElevationLoss string  `json:"elevationLoss"`
	Points        []Point `json:"points"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <input.gpx> <output.json>")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Error opening GPX file: %v\n", err)
		return
	}
	defer file.Close()

	gpxData, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Error reading GPX file: %v\n", err)
		return
	}

	var gpx GPX
	err = xml.Unmarshal(gpxData, &gpx)
	if err != nil {
		fmt.Printf("Error parsing GPX file: %v\n", err)
		return
	}

	var points []Point
	for _, trkPt := range gpx.Trk.TrkSeg.TrkPts {
		var timestamp time.Time
		if trkPt.Time != "" {
			parsedTimestamp, err := time.Parse(time.RFC3339, trkPt.Time)
			if err != nil {
				fmt.Printf("Error parsing time: %v\n", err)
				return
			} else {
				timestamp = parsedTimestamp
			}

		}
		points = append(points, Point{
			Longitude: trkPt.Lon,
			Latitude:  trkPt.Lat,
			Altitude:  trkPt.Ele,
			TimeStamp: timestamp.Unix(),
		})
	}

	output := Output{
		Distance:      "TODO",
		MovingTime:    "TODO",
		AverageSpeed:  "TODO",
		ElevationGain: "TODO",
		ElevationLoss: "TODO",
		Points:        points,
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling JSON: %v\n", err)
		return
	}

	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing JSON file: %v\n", err)
		return
	}

	fmt.Println("Successfully converted GPX to JSON")
}
