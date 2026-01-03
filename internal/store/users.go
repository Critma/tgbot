package store

import "time"

type User struct {
	ID         int
	TelegramID int
	CreatedAt  time.Time
}
