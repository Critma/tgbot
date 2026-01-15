package store

import (
	"context"
	"database/sql"
	"time"
)

const (
	QueryTimeoutDuration = 1 * time.Second
)

type Storage struct {
	Users interface {
		Create(context.Context, *User) error
		CreateOrUpdate(context.Context, *User) error
		GetByTelegramID(context.Context, int64) (*User, error)
		DeleteByTelegramID(context.Context, int64) error
	}
	Reminders interface {
		Create(context.Context, *sql.Tx, *Reminder) (*Reminder, error)
		Update(context.Context, *Reminder) error
		UpdateMessage(ctx context.Context, reminderID int, message string) error
		GetByID(context.Context, int) (*Reminder, error)
		GetByUserID(context.Context, int64) ([]*Reminder, error)
		DeleteByID(context.Context, int) error
	}
}

func WithTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
