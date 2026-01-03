package postgres

import (
	"database/sql"

	"github.com/critma/tgsheduler/internal/store"
)

func NewStorage(db *sql.DB) *store.Storage {
	return &store.Storage{
		Users:     &UsersStore{db},
		Reminders: &RemindersStore{db},
	}
}
