package repository

import (
	"database/sql"
	"log"
)

type Event struct {
	ID          int64
	TripID      int64
	Imei        string
	MessageCode int
	FreeText    string
	TimeStamp   int64
	Addresses   []Address
	Status      Status
	Latitude    float64
	Longitude   float64
	Altitude    float64
	GpsFix      int
	Course      float64
	Speed       float64
}

type Address struct {
	Address string
}

type Status struct {
	Autonomous     int
	LowBattery     int
	IntervalChange int
	ResetDetected  int
}

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Create(e Event) error {
	tx, err := r.db.Begin()
	if err != nil {
		log.Fatal("Couldn't begin save transaction for Event")
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO events(tripId, imei, messageCode, freeText, timeStamp, latitude, longitude, altitude, gpsFix, course, speed, autonomous, lowBattery, intervalChange, resetDetected) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(e.TripID, e.Imei, e.MessageCode, e.FreeText, e.TimeStamp, e.Latitude, e.Longitude, e.Altitude, e.GpsFix, e.Course, e.Speed, e.Status.Autonomous, e.Status.LowBattery, e.Status.IntervalChange, e.Status.ResetDetected)
	if err != nil {
		return err
	}

	eventID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	for _, addr := range e.Addresses {
		_, err := tx.Exec("INSERT INTO addresses(eventId, address) VALUES(?,?)", eventID, addr.Address)
		if err != nil {
			return err
		}
	}

	log.Printf("Saving new event to database")
	return tx.Commit()
}

func (r *EventRepository) All() ([]Event, error) {
	rows, err := r.db.Query(`SELECT id, latitude, longitude, altitude, speed, course, gpsFix, timeStamp FROM events ORDER BY timeStamp`)
	if err != nil {
		log.Printf("Error querying events: %v", err)
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event

		err := rows.Scan(&e.ID, &e.Latitude, &e.Longitude, &e.Altitude, &e.Speed, &e.Course, &e.GpsFix, &e.TimeStamp)
		if err != nil {
			log.Printf("Error scanning event row: %v", err)
		}
		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating event rows: %v", err)
	}

	return events, nil
}

func (r *EventRepository) Last() (Event, error) {
	var e Event
	row := r.db.QueryRow(`SELECT id, latitude, longitude, altitude, speed, course, gpsFix, timeStamp FROM events ORDER BY timeStamp DESC LIMIT 1`)
	err := row.Scan(&e.ID, &e.Latitude, &e.Longitude, &e.Altitude, &e.Speed, &e.Course, &e.GpsFix, &e.TimeStamp)
	if err != nil {
		log.Fatal(err)
		return Event{}, err
	}
	return e, nil
}

func (r *EventRepository) AllByDay(day string) ([]Event, error) {
	query := "SELECT id, latitude, longitude, altitude, timeStamp FROM events WHERE DATE(timeStamp, 'unixepoch') = ?"
	rows, err := r.db.Query(query, day)
	if err != nil {
		log.Printf("Error querying events: %v", err)
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event

		err := rows.Scan(&e.ID, &e.Latitude, &e.Longitude, &e.Altitude, &e.TimeStamp)
		if err != nil {
			log.Printf("Error scanning event row: %v", err)
		}
		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating event rows: %v", err)
	}

	return events, nil
}
