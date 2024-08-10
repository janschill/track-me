package repository

import (
	"database/sql"
	"log"
)

type Kudos struct {
	Day  string
	Count int
}

type KudosRepository struct {
	db *sql.DB
}

func NewKudosRepository(db *sql.DB) *KudosRepository {
	return &KudosRepository{db: db}
}

func (r *KudosRepository) Increment(day string) error {
	tx, err := r.db.Begin()
	if err != nil {
		log.Fatal("Couldn't begin save transaction for Kudos")
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO kudos(day, count) VALUES (?, 1) ON CONFLICT(day) DO UPDATE SET count = count + 1")
	if err != nil {
		log.Print("stmt", err)
		return err
	}

	_, err = stmt.Exec(day)
	if err != nil {
		log.Print("exec", err)
		return err
	}

	log.Printf("Saving new kudos to database")
	return tx.Commit()
}

func (r *KudosRepository) All() ([]Kudos, error) {
	rows, err := r.db.Query(`SELECT day, count FROM kudos ORDER BY day DESC`)
	if err != nil {
		log.Printf("Error querying kudos: %v", err)
		return nil, err
	}
	defer rows.Close()

	var kudos []Kudos
	for rows.Next() {
		var k Kudos

		err := rows.Scan(&k.Day, &k.Count)
		if err != nil {
			log.Printf("Error scanning kudos row: %v", err)
		}
		kudos = append(kudos, k)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating kudos rows: %v", err)
	}

	return kudos, nil
}
