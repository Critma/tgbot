package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/critma/tgsheduler/internal/store"
	"github.com/sanyokbig/pqinterval"
)

type RemindersStore struct {
	DB *sql.DB
}

func (r *RemindersStore) Create(ctx context.Context, reminder *store.Reminder) error {
	query := `
		INSERT INTO reminders (user_id, message, scheduled_time, repeat_interval)
		 VALUES ($1, $2, $3, $4)
	`

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

	_, err := r.DB.ExecContext(ctx, query, reminder.UserTelegramID, reminder.Message, reminder.SheduledTime, reminder.RepeatInterval)
	return err
}

func (r *RemindersStore) Update(ctx context.Context, reminder *store.Reminder) error {
	query := `
		UPDATE reminders
		 SET message = $1, scheduled_time = $2, repeat_interval = $3, is_active = $4
		 WHERE id = $5
	`

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

	result, err := r.DB.ExecContext(ctx, query, reminder.Message, reminder.SheduledTime, reminder.RepeatInterval, reminder.IsActive, reminder.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("zero affected rows")
	}
	return nil
}

func (r *RemindersStore) GetByUserID(ctx context.Context, userID int64) ([]*store.Reminder, error) {
	query := `
		SELECT id, user_id, message, scheduled_time, repeat_interval, is_active, created_at, updated_at from reminders WHERE user_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	reminders := []*store.Reminder{}
	for rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}

		reminder := &store.Reminder{}
		var ival pqinterval.Interval
		//TODO tz not scanned
		err := rows.Scan(&reminder.ID, &reminder.UserTelegramID, &reminder.Message, &reminder.SheduledTime, &ival, &reminder.IsActive, &reminder.CreatedAt, &reminder.UpdatedAt)
		if err != nil {
			return nil, err
		}

		reminder.RepeatInterval, err = ival.Duration()
		if err != nil {
			return nil, err
		}

		reminders = append(reminders, reminder)
	}

	return reminders, nil
}

func (r *RemindersStore) DeleteByID(ctx context.Context, id int) error {
	query := `
		DELETE FROM reminders
		 WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

	_, err := r.DB.ExecContext(ctx, query, id)
	return err
}
