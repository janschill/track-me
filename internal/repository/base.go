package repository

import "database/sql"

type Repository struct {
    Messages *MessageRepository
    Events  *EventRepository
    Days    *DayRepository
}

func NewRepository(db *sql.DB) *Repository {
    return &Repository{
        Messages: NewMessageRepository(db),
        Events:  NewEventRepository(db),
        Days:    NewDayRepository(db),
    }
}
