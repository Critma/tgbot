package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/critma/tgsheduler/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestUsersStore_GetByTelegramID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := &UsersStore{DB: db}

	telegramID := int64(123)
	var utcOffset int8 = 3 // UTC+3

	rows := sqlmock.NewRows([]string{"telegram_id", "utc", "created_at"}).
		AddRow(telegramID, utcOffset, time.Now())

	mock.ExpectQuery(
		`^SELECT .* FROM users WHERE telegram_id = \$1$`,
	).
		WithArgs(telegramID).
		WillReturnRows(rows)

	user, err := store.GetByTelegramID(context.Background(), telegramID)
	assert.NoError(t, err)
	assert.Equal(t, telegramID, user.TelegramID)
	assert.Equal(t, utcOffset, user.UTC)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUsersStore_GetByTelegramID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := &UsersStore{DB: db}

	telegramID := int64(999)

	mock.ExpectQuery(
		`^SELECT .* FROM users WHERE telegram_id = \$1$`,
	).
		WithArgs(telegramID).
		WillReturnError(sql.ErrNoRows)

	user, err := store.GetByTelegramID(context.Background(), telegramID)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, err, sql.ErrNoRows)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUsersStore_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	UserStore := &UsersStore{DB: db}

	user := &store.User{
		TelegramID: 123,
		UTC:        3,
	}

	mock.ExpectQuery(
		`^INSERT INTO users \(telegram_id, utc\) VALUES \(\$1, \$2\)$`,
	).
		WithArgs(user.TelegramID, user.UTC).
		WillReturnRows(sqlmock.NewRows(nil))

	err = UserStore.Create(context.Background(), user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUsersStore_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	UserStore := &UsersStore{DB: db}

	telegramID := int64(123)
	var newUTC int8 = 5

	mock.ExpectQuery(
		`^UPDATE users SET utc = \$2 WHERE telegram_id = \$1;$`,
	).
		WithArgs(telegramID, newUTC).
		WillReturnRows(sqlmock.NewRows(nil))

	user := &store.User{
		TelegramID: telegramID,
		UTC:        newUTC,
	}

	err = UserStore.Update(context.Background(), user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUsersStore_CreateOrUpdate_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	UserStore := &UsersStore{DB: db}

	telegramID := int64(123)
	var utc int8 = 3

	// Сначала GetByTelegramID вернёт ошибку (пользователя нет)
	mock.ExpectQuery(
		`^SELECT telegram_id, utc, created_at FROM users WHERE telegram_id = \$1$`,
	).WithArgs(telegramID).WillReturnError(sql.ErrNoRows)

	// Затем вызовется Create
	mock.ExpectQuery(
		`^INSERT INTO users \(telegram_id, utc\) VALUES \(\$1, \$2\)$`,
	).WithArgs(telegramID, utc).WillReturnRows(sqlmock.NewRows(nil))

	user := &store.User{
		TelegramID: telegramID,
		UTC:        utc,
	}

	err = UserStore.CreateOrUpdate(context.Background(), user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUsersStore_CreateOrUpdate_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	Usersstore := &UsersStore{DB: db}

	telegramID := int64(123)
	oldUTC := 2
	var newUTC int8 = 3

	// GetByTelegramID найдёт пользователя
	rows := sqlmock.NewRows([]string{"telegram_id", "utc", "created_at"}).
		AddRow(telegramID, oldUTC, time.Now())

	mock.ExpectQuery(
		`^SELECT telegram_id, utc, created_at FROM users WHERE telegram_id = \$1$`,
	).WithArgs(telegramID).WillReturnRows(rows)

	// Затем вызовется Update
	mock.ExpectQuery(
		`^UPDATE users SET utc = \$2 WHERE telegram_id = \$1;$`,
	).WithArgs(telegramID, newUTC).WillReturnRows(sqlmock.NewRows(nil))

	user := &store.User{
		TelegramID: telegramID,
		UTC:        newUTC,
	}

	err = Usersstore.CreateOrUpdate(context.Background(), user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUsersStore_DeleteByTelegramID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := &UsersStore{DB: db}

	telegramID := int64(123)

	mock.ExpectQuery(
		`^DELETE FROM users WHERE telegram_id = \$1$`,
	).
		WithArgs(telegramID).
		WillReturnRows(sqlmock.NewRows(nil))

	err = store.DeleteByTelegramID(context.Background(), telegramID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
