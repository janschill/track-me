package repository

import (
	"database/sql"
	"log"
)

type Day struct {
	ID                     int64
	TimeStamp              int64
	Points                 string
	TripID                 int64
	AverageSpeed           float64
	MaxSpeed               float64
	MinSpeed               float64
	TotalDistance          float64
	ElevationGain          int64
	ElevationLoss          int64
	AverageAltitude        float64
	MaxAltitude            float64
	MinAltitude            float64
	MovingTimeInSeconds    int64
	NumberOfStops          int64
	TotalStopTimeInSeconds int64
}

type DayRepository struct {
	db *sql.DB
}

func NewDayRepository(db *sql.DB) *DayRepository {
	return &DayRepository{db: db}
}

func (r *DayRepository) Create(d Day) error {
	tx, err := r.db.Begin()
	if err != nil {
		log.Fatal("Couldn't begin save transaction for Day")
		return err
	}
	stmt, err := tx.Prepare(`INSERT INTO events_cache (points, tripId, averageSpeed, maxSpeed, minSpeed, totalDistanceInMeters, elevationGain, elevationLoss, averageAltitude, maxAltitude, minAltitude, movingTimeInSeconds, numberOfStops, totalStopTimeInSeconds, timeStamp) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Fatal("Failed to prepare INSERT ", err)
		return err
	}
	_, err = stmt.Exec(d.Points, d.TripID, d.AverageSpeed, d.MaxSpeed, d.MinSpeed, d.TotalDistance, d.ElevationGain, d.ElevationLoss, d.AverageAltitude, d.MaxAltitude, d.MinAltitude, d.MovingTimeInSeconds, d.NumberOfStops, d.TotalStopTimeInSeconds, d.TimeStamp)
	if err != nil {
		log.Fatal("Failed to exec INSERT ", err)
		return err
	}

	log.Printf("Saving new day to database")
	return tx.Commit()
}

func (r *DayRepository) All() ([]Day, error) {
	rows, err := r.db.Query(`SELECT * FROM events_cache ORDER BY timeStamp`)
	if err != nil {
		log.Printf("Error querying days: %v", err)
		return nil, err
	}
	defer rows.Close()

	var days []Day
	for rows.Next() {
		var d Day

		err := rows.Scan(&d.ID, &d.Points, &d.TripID, &d.AverageSpeed, &d.MaxSpeed, &d.MinSpeed, &d.TotalDistance, &d.ElevationGain, &d.ElevationLoss, &d.AverageAltitude, &d.MaxAltitude, &d.MinAltitude, &d.MovingTimeInSeconds, &d.NumberOfStops, &d.TotalStopTimeInSeconds, &d.TimeStamp)
		if err != nil {
			log.Printf("Error scanning day row: %v", err)
		}
		days = append(days, d)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating day rows: %v", err)
	}

	return days, nil
}
