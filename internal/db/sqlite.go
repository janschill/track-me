package db

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
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

// var Db *sql.DB

// func closeDB() {
// 	if err := Db.Close(); err != nil {
// 		log.Fatal("Failed to close database connection:", err)
// 	}
// }

// type Point struct {
// 	Latitude  float64 `json:"latitude"`
// 	Longitude float64 `json:"longitude"`
// 	TimeStamp int64   `json:"timestamp"`
// 	EventID   int64   `json:eventID`
// }

func InitializeDB(filePath string) (*sql.DB, error) {
	Db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		return nil, err
	}

	if err = Db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Database connection established")

	return Db, nil
}

func DestroyDB(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		log.Fatalf("Failed to delete database file: %v", err)
	}

	log.Println("Database file deleted successfully.")
}

func CreateTables(filePath string) {
	Db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		return
	}
	defer Db.Close()
	for _, table := range schema.Tables {
		_, err := Db.Exec(table.Definition)
		if err != nil {
			log.Printf("Failed to create table %s: %v", table.Name, err)
			return
		}
	}
	log.Println("All tables created successfully.")
}

func Clear(filePath string) {
	Db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		return
	}
	defer Db.Close()
	for _, table := range schema.Tables {
		_, err := Db.Exec("DELETE FROM " + table.Name)
		if err != nil {
			log.Printf("Failed to clear table %s: %v", table.Name, err)
			return
		}
	}
	log.Println("All data cleared from the database.")
}

func convertPointsToEvents(points []TestPoint) []Event {
	events := make([]Event, len(points))
	for i, point := range points {
		events[i] = Event{
			Longitude: point.Longitude,
			Latitude:  point.Latitude,
			Altitude:  point.Altitude,
			TimeStamp: point.TimeStamp,
		}
	}
	return events
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

func Seed(filePath string) {
	Db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		return
	}
	defer Db.Close()
	path := "./internal/db/test_data/Great_Divide_2024.json"

	data, err := readGPXDataFromFile(path)
	if err != nil {
		log.Fatalf("Failed to read GPX data from file %s: %v", path, err)
	}
	startDate := time.Date(2024, time.September, 8, 0, 0, 0, 0, time.UTC)
	totalPoints := len(data.Points)
	totalDuration := time.Hour * 24 * time.Duration(totalPoints/1500) // total duration in days
	timeIncrement := totalDuration / time.Duration(totalPoints)       // time increment per point
	currentTime := startDate
	rowCount := 0
	// every 1500 entrys increment day
	for _, point := range data.Points {
		timeStamp := currentTime.Unix()
		_, err = Db.Exec("INSERT INTO events(tripId, imei, messageCode, timeStamp, latitude, longitude, altitude, gpsFix, course, speed, autonomous, lowBattery, intervalChange, resetDetected) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
			1, "fake-imei", 0, timeStamp, point.Latitude, point.Longitude, point.Altitude, 0, 0, 0, 0, 0, 0, 0)
		if err != nil {
			log.Fatal("Failed to insert into events table:", err)
		}

		currentTime = currentTime.Add(timeIncrement)
		rowCount++
	}
}

func Aggregate(filePath string, day string) {
	Db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		log.Fatal("Failed to open database connection:", err)
		return
	}
	defer Db.Close()

	// events, err := GetAllEvents(Db)
	events, err := GetAllEventsByDay(Db, day)
	if err != nil {
		log.Fatal("Failed to get events by day:", err)
		return
	}
	movingTime, averageSpeed := CalculateMovingTimeAndAverageSpeed(events, 0.001)
	gain, loss := CalculateElevationGainAndLoss(events)
	// maxSpeed := CalculateMaxSpeed(events)
	averageAltitude, maxAltitude, minAltitude := CalculateAltitudes(events)
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

	log.Printf("Reducing events from %v", len(events))
	reduced_events := Rdp(events, 0.0002) // roughly 1500 -> 321
	log.Printf("Number of reduced_events after Rdp %v", len(reduced_events))

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

	reduced_events_json, err := json.Marshal(reduced_events)
	if err != nil {
		log.Fatal("Failed to marshal points data:", err)
		return
	}

	d := Day{
		TimeStamp:              events[0].TimeStamp,
		Points:                 string(reduced_events_json),
		TripID:                 reduced_events[0].TripID,
		AverageSpeed:           averageSpeed,
		MaxSpeed:               0,
		MinSpeed:               0,
		TotalDistance:          CalculateDistance(reduced_events),
		ElevationGain:          gain,
		ElevationLoss:          loss,
		AverageAltitude:        float64(averageAltitude),
		MaxAltitude:            maxAltitude,
		MinAltitude:            minAltitude,
		MovingTimeInSeconds:    int64(movingTime),
		NumberOfStops:          0,
		TotalStopTimeInSeconds: 0,
	}

	err = d.Save(Db)
	if err != nil {
		log.Fatal("Failed to save day")
	}
}
