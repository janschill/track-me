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
			Altitude:  int64(point.Altitude),
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
	currentDate := startDate.Add(time.Hour * 24 * time.Duration(len(data.Points)/1500))
	timeStamp := currentDate.Unix()
	// every 1500 entrys increment day
	rowCount := 0
	for _, point := range data.Points {
		_, err = Db.Exec("INSERT INTO events(tripId, imei, messageCode, timeStamp, latitude, longitude, altitude, gpsFix, course, speed, autonomous, lowBattery, intervalChange, resetDetected) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
			1, "fake-imei", 0, timeStamp, point.Latitude, point.Longitude, int(point.Altitude), 0, 0, 0, 0, 0, 0, 0)
		if err != nil {
			log.Fatal("Failed to insert into events table:", err)
		}

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

	events, err := GetAllEvents(Db)
	// events, err := GetAllEventsByDay(Db, day)
	if err != nil {
		log.Fatal("Failed to get events by day:", err)
		return
	}

	log.Printf("Reducing events from %v", len(events))
	reduced_events := Rdp(events, 0.0002) // roughly 1500 -> 321
	log.Printf(" to %v", len(reduced_events))

	distance := CalculateDistance(events)
	log.Printf("Distance all events %v", distance)
	distance = CalculateDistance(reduced_events)
	log.Printf("Distance reduced events %v", distance)

	reduced_events_json, err := json.Marshal(reduced_events)
	if err != nil {
		log.Fatal("Failed to marshal points data:", err)
		return
	}

	return
	d := Day{
		Day:                    day,
		Points:                 string(reduced_events_json),
		TripID:                 reduced_events[0].TripID,
		AverageSpeed:           0,
		MaxSpeed:               0,
		MinSpeed:               0,
		TotalDistance:          0,
		ElevationGain:          0,
		ElevationLoss:          0,
		AverageAltitude:        0,
		MaxAltitude:            0,
		MinAltitude:            0,
		MovingTimeInSeconds:    0,
		NumberOfStops:          0,
		TotalStopTimeInSeconds: 0,
	}

	err = d.Save(Db)
	if err != nil {
		log.Fatal("Failed to save day")
	}
}
