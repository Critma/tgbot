package postgres

import (
	"context"
	"database/sql"

	"github.com/critma/tgsheduler/internal/store"
)

type UsersStore struct {
	DB *sql.DB
}

func (u *UsersStore) Create(ctx context.Context, user *store.User) error {
	query := `
		INSERT INTO users (id, telegram_id, created_at) VALUES ($1, $2, $3)
	`

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

	return u.DB.QueryRowContext(ctx, query, user.ID, user.TelegramID, user.CreatedAt).Err()
}

func (u *UsersStore) DeleteByID(ctx context.Context, id int) error {
	query := `
		DELETE FROM users WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

	return u.DB.QueryRowContext(ctx, query, id).Err()
}
