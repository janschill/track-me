package repository

import "database/sql"

type Repository struct {
	Messages    *MessageRepository
	Events      *EventRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Messages:    NewMessageRepository(db),
		Events:      NewEventRepository(db),
	}
}
