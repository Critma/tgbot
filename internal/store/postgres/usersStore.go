package postgres

import (
	"context"
	"database/sql"

	"github.com/critma/tgsheduler/internal/store"
)

type UsersStore struct {
	DB *sql.DB
}

func (u *UsersStore) Create(ctx context.Context, userTelegramID int64) error {
	query := `
		INSERT INTO users (telegram_id) VALUES ($1)
	`

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

	return u.DB.QueryRowContext(ctx, query, userTelegramID).Err()
}

func (u *UsersStore) GetByTelegramID(ctx context.Context, userTelegramID int64) (*store.User, error) {
	query := `
			SELECT * FROM users WHERE telegram_id = $1
		`

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

	user := &store.User{}
	if err := u.DB.QueryRowContext(ctx, query, userTelegramID).Scan(&user.TelegramID, &user.CreatedAt); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UsersStore) DeleteByTelegramID(ctx context.Context, userTelegramID int64) error {
	query := `
		DELETE FROM users WHERE telegram_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

	return u.DB.QueryRowContext(ctx, query, userTelegramID).Err()
}
