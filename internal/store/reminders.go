package store

import "time"

type Reminder struct {
	ID             int
	UserID         int
	Message        string
	SheduledTime   time.Time
	RepeatInterval time.Duration
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
