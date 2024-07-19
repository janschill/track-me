package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

func closeDB() {
	if err := Db.Close(); err != nil {
		log.Fatal("Failed to close database connection:", err)
	}
}

func initializeDB(filePath string) {
	var err error
	Db, err = sql.Open("sqlite3", filePath)
	if err != nil {
		log.Fatal(err)
	}

	if err = Db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to the SQLite database successfully.")
}

func DestroyDB(filePath string) {
  err := os.Remove(filePath)
  if err != nil {
      log.Fatalf("Failed to delete database file: %v", err)
  }

  log.Println("Database file deleted successfully.")
}

func ResetDB(filePath string) {
	initializeDB(filePath)
	tables := []string{"events", "trips", "addresses"}
	for _, table := range tables {
		dropSQL := "DROP TABLE IF EXISTS " + table + ";"
		_, err := Db.Exec(dropSQL)
		if err != nil {
			log.Fatalf("Failed to drop %s table: %v", table, err)
		}
		log.Printf("%s table dropped", table)
	}

	log.Println("Database reset successfully.")
  closeDB()
}


func CreateTables(filePath string) {
	initializeDB(filePath)
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
        "altitude" REAL,
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

	_, err := Db.Exec(createTripTableSQL)
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

	log.Println("Tables created successfully.")
  closeDB()
}

func Seed(filePath string) {

}
