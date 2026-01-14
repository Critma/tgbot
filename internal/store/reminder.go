package store

import "time"

type Reminder struct {
	ID             int
	UserTelegramID int64
	Message        string
	SheduledTime   time.Time
	RepeatInterval time.Duration
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	TaskID         string
	TaskQueue      string
}
