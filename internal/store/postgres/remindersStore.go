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

func (r *RemindersStore) Create(ctx context.Context, tx *sql.Tx, reminder *store.Reminder) (*store.Reminder, error) {
	query := `
		INSERT INTO reminders (user_id, message, scheduled_time, repeat_interval, task_id, task_queue)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, is_active
	`

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

	row := tx.QueryRowContext(ctx, query, reminder.UserTelegramID,
		reminder.Message, reminder.SheduledTime, reminder.RepeatInterval,
		reminder.TaskID, reminder.TaskQueue)
	if row.Err() != nil {
		return nil, row.Err()
	}
	if err := row.Scan(&reminder.ID, &reminder.IsActive); err != nil {
		return nil, err
	}

	return reminder, nil
}

func (r *RemindersStore) Update(ctx context.Context, tx *sql.Tx, reminder *store.Reminder) error {
	query := `
		UPDATE reminders
		 SET message = $1, scheduled_time = $2, repeat_interval = $3, is_active = $4, task_id = $5, task_queue = $6
		 WHERE id = $7
	`

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

	result, err := tx.ExecContext(ctx, query, reminder.Message, reminder.SheduledTime, reminder.RepeatInterval, reminder.IsActive, reminder.TaskID, reminder.TaskQueue, reminder.ID)
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

func (r *RemindersStore) UpdateMessage(ctx context.Context, reminderID int, message string) error {
	query := `
		UPDATE reminders
		 SET message = $1
		 WHERE id = $2
	`

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

	result, err := r.DB.ExecContext(ctx, query, message, reminderID)
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

func (r *RemindersStore) UpdateIsActive(ctx context.Context, reminderID int, isActive bool) error {
	query := `
		UPDATE reminders
		 SET is_active = $1
		 WHERE id = $2
	`

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

	result, err := r.DB.ExecContext(ctx, query, isActive, reminderID)
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

// return time in UTC+0
func (r *RemindersStore) GetByUserID(ctx context.Context, userID int64) ([]*store.Reminder, error) {
	query := `
		SELECT id, user_id, message, scheduled_time, repeat_interval, is_active, created_at, updated_at, task_id, task_queue
		FROM reminders WHERE user_id = $1
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
		err := rows.Scan(&reminder.ID, &reminder.UserTelegramID, &reminder.Message, &reminder.SheduledTime, &ival, &reminder.IsActive, &reminder.CreatedAt, &reminder.UpdatedAt, &reminder.TaskID, &reminder.TaskQueue)
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

// return time in UTC+0
func (r *RemindersStore) GetByID(ctx context.Context, id int) (*store.Reminder, error) {
	query := `
		SELECT id, user_id, message, scheduled_time, repeat_interval, is_active, created_at, updated_at, task_id, task_queue
		FROM reminders
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

	row := r.DB.QueryRowContext(ctx, query, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	reminder := &store.Reminder{}
	var ival pqinterval.Interval
	err := row.Scan(&reminder.ID, &reminder.UserTelegramID, &reminder.Message, &reminder.SheduledTime, &ival, &reminder.IsActive, &reminder.CreatedAt, &reminder.UpdatedAt, &reminder.TaskID, &reminder.TaskQueue)
	if err != nil {
		return nil, err
	}

	reminder.RepeatInterval, err = ival.Duration()
	if err != nil {
		return nil, err
	}

	return reminder, nil
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

func (r *RemindersStore) GetActiveByUserID(ctx context.Context, userID int64) ([]*store.Reminder, error) {
	reminders, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]*store.Reminder, 0, len(reminders))
	for _, reminder := range reminders {
		if reminder.IsActive {
			result = append(result, reminder)
		}
	}
	return result, nil
}
