package repository

import (
	"database/sql"
	"log"
)

type Message struct {
	ID           int64
	TripID       int64
	Message      string
	Name         string
	TimeStamp    int64
	SentToGarmin bool
	FromGarmin bool
}

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(m Message) error {
	tx, err := r.db.Begin()
	if err != nil {
		log.Fatal("Couldn't begin save transaction for Message")
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO messages(tripId, message, name, timeStamp, sentToGarmin, fromGarmin) VALUES(?,?,?,?,?,?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(m.TripID, m.Message, m.Name, m.TimeStamp, m.SentToGarmin, m.FromGarmin)
	if err != nil {
		return err
	}

	log.Printf("Saving new message to database")
	return tx.Commit()
}

func (r *MessageRepository) All() ([]Message, error) {
	rows, err := r.db.Query(`SELECT id, tripId, message, name, timeStamp, sentToGarmin, fromGarmin FROM messages ORDER BY timeStamp DESC`)
	if err != nil {
		log.Printf("Error querying messages: %v", err)
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message

		err := rows.Scan(&m.ID, &m.TripID, &m.Message, &m.Name, &m.TimeStamp, &m.SentToGarmin, &m.FromGarmin)
		if err != nil {
			log.Printf("Error scanning message row: %v", err)
		}
		messages = append(messages, m)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating message rows: %v", err)
	}

	return messages, nil
}
