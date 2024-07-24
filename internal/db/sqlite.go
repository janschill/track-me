package db

import (
	"bufio"
	"database/sql"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// var Db *sql.DB

// func closeDB() {
// 	if err := Db.Close(); err != nil {
// 		log.Fatal("Failed to close database connection:", err)
// 	}
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

// func ResetDB(filePath string) {
// 	InitializeDB(filePath)
// 	tables := []string{"events", "trips", "addresses"}
// 	for _, table := range tables {
// 		dropSQL := "DROP TABLE IF EXISTS " + table + ";"
// 		_, err := Db.Exec(dropSQL)
// 		if err != nil {
// 			log.Fatalf("Failed to drop %s table: %v", table, err)
// 		}
// 		log.Printf("%s table dropped", table)
// 	}

// 	log.Println("Database reset successfully.")
// 	closeDB()
// }

func CreateTables(filePath string) {
	Db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		return
	}
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
      date DATE PRIMARY KEY,
			points_data TEXT NOT NULL
    );`

	createMessagesTableSQL := `CREATE TABLE IF NOT EXISTS messages (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      timeStamp INTEGER NOT NULL,
      name TEXT,
      message TEXT,
      sentToGarmin INTEGER,
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
	if err := Db.Close(); err != nil {
		log.Fatal("Failed to close database connection:", err)
	}
}

func Seed(filePath string) {
	Db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		return
	}
	defer Db.Close()

	file, err := os.Open("./data/route_points.txt")
	if err != nil {
		log.Fatal("Failed to open file:", err)
	}
	defer file.Close()

	startDate := time.Date(2024, time.September, 8, 0, 0, 0, 0, time.UTC)
	// every 1500 entrys increment day
	rowCount := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ", ")
		latitudeStr := strings.Split(parts[0], ": ")[1]
		longitudeStr := strings.Split(parts[1], ": ")[1]
		elevationStr := strings.Split(parts[2], ": ")[1]
		latitude, err := strconv.ParseFloat(latitudeStr, 64)
		if err != nil {
			log.Fatal("Failed to parse latitude:", err)
		}
		longitude, err := strconv.ParseFloat(longitudeStr, 64)
		if err != nil {
			log.Fatal("Failed to parse latitude:", err)
		}
		elevation, err := strconv.ParseFloat(elevationStr, 64)
		if err != nil {
			log.Fatal("Failed to parse latitude:", err)
		}

		currentDate := startDate.Add(time.Hour * 24 * time.Duration(rowCount/1500))
		timeStamp := currentDate.Unix()

		_, err = Db.Exec("INSERT INTO events(tripId, imei, messageCode, timeStamp, latitude, longitude, altitude, gpsFix, course, speed, autonomous, lowBattery, intervalChange, resetDetected) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
			1, "fake-imei", 0, timeStamp, latitude, longitude, int(elevation), 0, 0, 0, 0, 0, 0, 0)
		if err != nil {
			log.Fatal("Failed to insert into events table:", err)
		}

		rowCount++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("Error reading file:", err)
	}
}
