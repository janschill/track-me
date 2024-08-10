package repository

import "database/sql"

type Repository struct {
	Messages *MessageRepository
	Events   *EventRepository
	Kudos    *KudosRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Messages: NewMessageRepository(db),
		Events:   NewEventRepository(db),
		Kudos:    NewKudosRepository(db),
	}
}
