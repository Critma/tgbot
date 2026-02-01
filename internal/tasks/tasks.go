package tasks

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/critma/tgsheduler/internal/config"
	"github.com/critma/tgsheduler/internal/domain/helpers"
	"github.com/critma/tgsheduler/internal/store"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

// Types

const (
	TypeReminderDelivery = "reminder:deliver"
)

type ReminderDeliveryPayload struct {
	ReminderID int
	UserID     int64
}

// Tasks creation

func NewReminderDeliveryTask(reminderID int, userID int64) (*asynq.Task, error) {
	payload, err := json.Marshal(ReminderDeliveryPayload{ReminderID: reminderID, UserID: userID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeReminderDelivery, payload), nil
}

// Tasks handlers

type ReminderProcessor struct {
	bot *tgbotapi.BotAPI
	app *config.Application
}

func NewReminderProcessor(bot *tgbotapi.BotAPI, app *config.Application) *ReminderProcessor {
	return &ReminderProcessor{
		bot: bot,
		app: app,
	}
}

func (p *ReminderProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var payload ReminderDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		log.Error().Str("event", "task handling").Str("message", "json.Unmarshal failed").Err(err).Send()
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	reminder, err := p.app.Store.Reminders.GetByID(ctx, payload.ReminderID)
	if err != nil {
		log.Error().Str("event", "task handling").Str("message", "failed to get reminder").Err(err).Send()
		return err
	}

	user, err := p.app.Store.Users.GetByTelegramID(context.Background(), payload.UserID)
	if err != nil {
		log.Error().Str("message", "failed to get user timezone").Err(err).Int64("userID", payload.UserID).Send()
		return err
	}

	reminder.SheduledTime = helpers.TimeToUserTZ(user, reminder.SheduledTime)
	var msg string
	if reminder.RepeatInterval.Hours() == 24 {
		msg = fmt.Sprintf("❗Уведомление❗\n[ежедневно %s] - %s", reminder.SheduledTime.Format("15:04"), reminder.Message)
	} else if reminder.RepeatInterval.Hours() == 24*7 {
		msg = fmt.Sprintf("❗Уведомление❗\n[еженедельно %s] - %s", reminder.SheduledTime.Format("15:04"), reminder.Message)
	} else {
		msg = fmt.Sprintf("❗Уведомление❗\n[%s] - %s", reminder.SheduledTime.Format("02.01.2006 15:04"), reminder.Message)
	}
	p.bot.Send(tgbotapi.NewMessage(payload.UserID, msg))
	log.Info().Str("event", "send event to tg").Int("reminderID", reminder.ID).Int64("userID", payload.UserID).Send()

	if reminder.RepeatInterval.Hours() >= 24 {
		err := p.reEnqueueTask(reminder, payload.UserID)
		return err
	} else {
		err = p.app.Store.Reminders.UpdateIsActive(context.Background(), payload.ReminderID, false)
		if err != nil {
			log.Warn().Str("message", "failed to set is_active = false").Int("reminderID", payload.ReminderID).Err(err).Int64("userID", payload.UserID).Send()
		}
	}

	return nil
}

func (p *ReminderProcessor) reEnqueueTask(reminder *store.Reminder, userID int64) error {
	task, err := NewReminderDeliveryTask(reminder.ID, userID)
	if err != nil {
		log.Error().Str("event", "re send to broker").Str("message", "could not create a task").Err(err).Send()
		return err
	}

	return store.WithTx(p.app.Db, context.Background(), func(tx *sql.Tx) error {
		reminder.SheduledTime = getRightTime(reminder.SheduledTime, reminder.RepeatInterval)
		if err := p.app.Store.Reminders.Update(context.Background(), tx, reminder); err != nil {
			log.Error().Str("event", "update reminder with new sheduled time").Str("message", "could not update the reminder").Err(err).Send()
			return err
		}

		info, err := p.app.Broker.Client.Enqueue(task, asynq.MaxRetry(1), asynq.ProcessAt(reminder.SheduledTime), asynq.Queue(reminder.TaskQueue))
		if err != nil {
			log.Error().Str("event", "re send to broker").Str("message", "could not re send a task").Err(err).Send()
			return err
		}

		reminder.TaskID = info.ID
		reminder.TaskQueue = info.Queue
		if err := p.app.Store.Reminders.Update(context.Background(), tx, reminder); err != nil {
			log.Error().Str("event", "update reminder with new task_info").Str("message", "could not update the reminder").Err(err).Send()
			//TODO:delete task from redis
			return err
		}

		log.Info().Str("event", "re sheduled task").Str("taskID", info.ID).Str("queue", info.Queue).Send()
		return nil
	})
}

func getRightTime(oldTime time.Time, repeatInterval time.Duration) time.Time {
	today := time.Now()
	updatedTime := time.Date(today.Year(), today.Month(), today.Day(), oldTime.Hour(), oldTime.Minute(), oldTime.Second(), oldTime.Nanosecond(), oldTime.Location())

	for updatedTime.Before(today) {
		updatedTime = updatedTime.Add(repeatInterval)
	}
	return updatedTime
}
