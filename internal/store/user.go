package store

import "time"

type User struct {
	TelegramID int64
	CreatedAt  time.Time
	UTC        int8
}
