package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/critma/tgsheduler/internal/config"
	"github.com/critma/tgsheduler/internal/domain/helpers"
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
	msg := tgbotapi.NewMessage(payload.UserID, fmt.Sprintf("❗Уведомление❗\n[%s] - %s", reminder.SheduledTime.Format("02.01.2006 15:04"), reminder.Message))
	p.bot.Send(msg)
	log.Info().Str("event", "send event to tg").Int64("userID", payload.UserID).Send()

	return nil
}
