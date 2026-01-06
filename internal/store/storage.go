package store

import (
	"context"
	"time"
)

const (
	QueryTimeoutDuration = 1 * time.Second
)

type Storage struct {
	Users interface {
		Create(context.Context, int64) error
		GetByTelegramID(context.Context, int64) (*User, error)
		DeleteByTelegramID(context.Context, int64) error
	}
	Reminders interface {
		Create(context.Context, *Reminder) error
		Update(context.Context, *Reminder) error
		GetByUserID(context.Context, int) ([]*Reminder, error)
		DeleteByID(context.Context, int) error
	}
}
