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
		INSERT INTO users (telegram_id, utc) VALUES ($1, $2)
	`

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

	return u.DB.QueryRowContext(ctx, query, user.TelegramID, user.UTC).Err()
}

func (u *UsersStore) CreateOrUpdate(ctx context.Context, user *store.User) error {
	existed, _ := u.GetByTelegramID(ctx, user.TelegramID)
	if existed == nil {
		return u.Create(ctx, user)
	} else {
		return u.Update(ctx, user)
	}
}

func (u *UsersStore) Update(ctx context.Context, user *store.User) error {
	query := `
			UPDATE users SET utc = $2 WHERE telegram_id = $1;
		`

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

	return u.DB.QueryRowContext(ctx, query, user.TelegramID, user.UTC).Err()
}

func (u *UsersStore) GetByTelegramID(ctx context.Context, userTelegramID int64) (*store.User, error) {
	query := `
			SELECT telegram_id, utc, created_at FROM users WHERE telegram_id = $1
		`

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

	user := &store.User{}
	if err := u.DB.QueryRowContext(ctx, query, userTelegramID).Scan(&user.TelegramID, &user.UTC, &user.CreatedAt); err != nil {
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
