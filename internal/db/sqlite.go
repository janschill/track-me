package db

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/janschill/track-me/internal/utils"
	_ "github.com/mattn/go-sqlite3"
)

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

func Seed(filePath string) {
	Db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		return
	}
	defer Db.Close()
	path := "./internal/db/test_data/Great_Divide_2024.json"

	data, err := utils.ReadGPXDataFromFile(path)
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
