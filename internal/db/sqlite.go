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
}

type GPXData struct {
	Distance string      `json:"distance"`
	Points   []TestPoint `json:"points"`
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

	createTripTableSQL := `CREATE TABLE IF NOT EXISTS trips (
        "id" INTEGER PRIMARY KEY AUTOINCREMENT,
        "startTime" DATETIME,
        "endTime" DATETIME,
        "description" TEXT
    );`

	createEventTableSQL := `CREATE TABLE IF NOT EXISTS events (
        "id" INTEGER PRIMARY KEY AUTOINCREMENT,
        "tripId" INTEGER NOT NULL,
        "imei" TEXT NOT NULL,
        "messageCode" INTEGER NOT NULL,
        "freeText" TEXT,
        "timeStamp" INTEGER NOT NULL,
        "latitude" REAL,
        "longitude" REAL,
        "altitude" INTEGER,
        "gpsFix" INTEGER,
        "course" REAL,
        "speed" REAL,
        "autonomous" INTEGER,
        "lowBattery" INTEGER,
        "intervalChange" INTEGER,
        "resetDetected" INTEGER,
        FOREIGN KEY(tripId) REFERENCES trips(id)
    );`

	createAddressTableSQL := `CREATE TABLE IF NOT EXISTS addresses (
      "id" INTEGER PRIMARY KEY AUTOINCREMENT,
      "eventId" INTEGER NOT NULL,
      "address" TEXT NOT NULL,
      FOREIGN KEY (eventId) REFERENCES Event(id)
    );`

	createEventsCacheTableSQL := `CREATE TABLE IF NOT EXISTS events_cache (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
			points TEXT NOT NULL,
	    tripId INTEGER NOT NULL,
			averageSpeed REAL,
			maxSpeed REAL,
			minSpeed REAL,
			totalDistanceInMeters INTEGER,
			elevationGain INTEGER,
			elevationLoss INTEGER,
			averageAltitude REAL,
			maxAltitude INTEGER,
			minAltitude INTEGER,
			movingTimeInSeconds INTEGER,
			numberOfStops INTEGER,
			totalStopTimeInSeconds INTEGER,
	    date DATE NOT NULL
    );`

	createMessagesTableSQL := `CREATE TABLE IF NOT EXISTS messages (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      tripId INTEGER NOT NULL,
      timeStamp INTEGER NOT NULL,
      name TEXT,
      message TEXT,
      sentToGarmin INTEGER,
      FOREIGN KEY(tripId) REFERENCES trips(id)
    );`

	_, err = Db.Exec(createTripTableSQL)
	if err != nil {
		log.Fatal("Failed to create trips table:", err)
	}

	_, err = Db.Exec(createEventTableSQL)
	if err != nil {
		log.Fatal("Failed to create events table:", err)
	}

	_, err = Db.Exec(createAddressTableSQL)
	if err != nil {
		log.Fatal("Failed to create addresses table:", err)
	}

	_, err = Db.Exec(createEventsCacheTableSQL)
	if err != nil {
		log.Fatal("Failed to create events_cache table:", err)
	}

	_, err = Db.Exec(createMessagesTableSQL)
	if err != nil {
		log.Fatal("Failed to create messages table:", err)
	}

	log.Println("Tables created successfully.")
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
