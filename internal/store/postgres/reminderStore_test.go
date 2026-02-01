package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/critma/tgsheduler/internal/store"
	"github.com/sanyokbig/pqinterval"
	"github.com/stretchr/testify/assert"
)

func TestRiminderStore_GetByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := &RemindersStore{DB: db}

	userID := int64(123)
	Expectmessage := "Test reminder"
	ExpectedTaskID := "task-1"
	ExpectedTaskQueue := "queue-1"
	rows := sqlmock.NewRows([]string{
		"id", "user_id", "message", "scheduled_time", "repeat_interval",
		"is_active", "created_at", "updated_at", "task_id", "task_queue",
	}).AddRow(1, userID, Expectmessage, time.Now().UTC(), pqinterval.New(0, 1, 0, 0, 0, 0), true, time.Now(), time.Now(), ExpectedTaskID, ExpectedTaskQueue)

	mock.ExpectQuery(`SELECT .* FROM reminders WHERE user_id = \$1`).WithArgs(userID).WillReturnRows(rows)

	reminders, err := store.GetByUserID(context.Background(), userID)
	assert.NoError(t, err)
	assert.Len(t, reminders, 1)
	assert.Equal(t, Expectmessage, reminders[0].Message)
	assert.Equal(t, userID, reminders[0].UserTelegramID)
	ExpectedDuration := time.Hour * 24
	assert.Equal(t, ExpectedDuration, reminders[0].RepeatInterval)
	assert.Equal(t, true, reminders[0].IsActive)
	assert.Equal(t, ExpectedTaskID, reminders[0].TaskID)
	assert.Equal(t, ExpectedTaskQueue, reminders[0].TaskQueue)
}

func TestRemindersStore_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	ReminderStore := &RemindersStore{DB: db}

	reminder := store.Reminder{
		UserTelegramID: 123,
		Message:        "Test reminder",
		SheduledTime:   time.Now().UTC(),
		RepeatInterval: 24 * time.Hour,
		TaskID:         "task-1",
		TaskQueue:      "queue-1",
	}

	rows := sqlmock.NewRows([]string{"id", "is_active"}).
		AddRow(1, true)

	mock.ExpectBegin()
	mock.ExpectQuery(
		`^INSERT INTO reminders \(user_id, message, scheduled_time, repeat_interval, task_id, task_queue\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\) RETURNING id, is_active$`,
	).
		WithArgs(
			reminder.UserTelegramID,
			reminder.Message,
			reminder.SheduledTime,
			pqinterval.Duration(reminder.RepeatInterval),
			reminder.TaskID,
			reminder.TaskQueue,
		).
		WillReturnRows(rows)
	mock.ExpectCommit()
	tx, _ := db.Begin()
	_, err = ReminderStore.Create(context.Background(), tx, &reminder)

	assert.NoError(t, err)
	assert.Equal(t, 1, reminder.ID)
	assert.True(t, reminder.IsActive)

	err = tx.Commit()

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRemindersStore_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	ReminderStore := &RemindersStore{DB: db}

	reminder := &store.Reminder{
		ID:             1,
		UserTelegramID: 123,
		Message:        "Updated reminder",
		SheduledTime:   time.Now().Add(time.Hour).UTC(),
		RepeatInterval: 2 * time.Hour,
		IsActive:       false,
		TaskID:         "new-task",
		TaskQueue:      "new-queue",
	}

	mock.ExpectBegin()
	mock.ExpectExec(
		`^UPDATE reminders SET message = \$1, scheduled_time = \$2, repeat_interval = \$3, is_active = \$4, task_id = \$5, task_queue = \$6 WHERE id = \$7$`,
	).
		WithArgs(
			reminder.Message,
			reminder.SheduledTime,
			pqinterval.Duration(reminder.RepeatInterval),
			reminder.IsActive,
			reminder.TaskID,
			reminder.TaskQueue,
			reminder.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	tx, _ := db.Begin()
	err = ReminderStore.Update(context.Background(), tx, reminder)
	assert.NoError(t, err)
	err = tx.Commit()
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRemindersStore_UpdateMessage(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := &RemindersStore{DB: db}

	reminderID := 1
	newMessage := "New message text"

	mock.ExpectExec(
		`^UPDATE reminders SET message = \$1 WHERE id = \$2$`,
	).
		WithArgs(newMessage, reminderID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = store.UpdateMessage(context.Background(), reminderID, newMessage)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRemindersStore_UpdateIsActive(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := &RemindersStore{DB: db}

	reminderID := 1
	isActive := false

	mock.ExpectExec(
		`^UPDATE reminders SET is_active = \$1 WHERE id = \$2$`,
	).
		WithArgs(isActive, reminderID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = store.UpdateIsActive(context.Background(), reminderID, isActive)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRemindersStore_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := &RemindersStore{DB: db}

	reminderID := 1
	userID := int64(123)
	message := "Single reminder"
	taskID := "task-1"
	taskQueue := "queue-1"

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "message", "scheduled_time", "repeat_interval",
		"is_active", "created_at", "updated_at", "task_id", "task_queue",
	}).AddRow(
		reminderID, userID, message, time.Now().UTC(),
		pqinterval.New(0, 1, 0, 0, 0, 0),
		true, time.Now(), time.Now(), taskID, taskQueue,
	)

	mock.ExpectQuery(
		`^SELECT id, user_id, message, scheduled_time, repeat_interval, is_active, created_at, updated_at, task_id, task_queue FROM reminders WHERE id = \$1$`,
	).
		WithArgs(reminderID).
		WillReturnRows(rows)

	reminder, err := store.GetByID(context.Background(), reminderID)
	assert.NoError(t, err)
	assert.Equal(t, reminderID, reminder.ID)
	assert.Equal(t, message, reminder.Message)
	assert.Equal(t, userID, reminder.UserTelegramID)
	assert.Equal(t, 24*time.Hour, reminder.RepeatInterval)
	assert.Equal(t, taskID, reminder.TaskID)
	assert.Equal(t, taskQueue, reminder.TaskQueue)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRemindersStore_DeleteByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := &RemindersStore{DB: db}

	reminderID := 1

	mock.ExpectExec(
		`^DELETE FROM reminders WHERE id = \$1$`,
	).
		WithArgs(reminderID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = store.DeleteByID(context.Background(), reminderID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRemindersStore_GetActiveByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := &RemindersStore{DB: db}

	userID := int64(123)

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "message", "scheduled_time", "repeat_interval",
		"is_active", "created_at", "updated_at", "task_id", "task_queue",
	}).
		AddRow(1, userID, "Active reminder", time.Now(), pqinterval.New(0, 1, 0, 0, 0, 0), true, time.Now(), time.Now(), "t1", "q1").
		AddRow(2, userID, "Inactive reminder", time.Now(), pqinterval.New(0, 1, 0, 0, 0, 0), false, time.Now(), time.Now(), "t2", "q2")

	mock.ExpectQuery(
		`^SELECT id, user_id, message, scheduled_time, repeat_interval, is_active, created_at, updated_at, task_id, task_queue FROM reminders WHERE user_id = \$1$`,
	).
		WithArgs(userID).
		WillReturnRows(rows)

	reminders, err := store.GetActiveByUserID(context.Background(), userID)
	assert.NoError(t, err)
	assert.Len(t, reminders, 1)
	assert.True(t, reminders[0].IsActive)
	assert.Equal(t, "Active reminder", reminders[0].Message)
	assert.NoError(t, mock.ExpectationsWereMet())
}
