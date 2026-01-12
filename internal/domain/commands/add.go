package commands

import (
	"context"
	"strings"
	"time"

	"github.com/critma/tgsheduler/internal/logger"
	"github.com/critma/tgsheduler/internal/store"
	"github.com/critma/tgsheduler/internal/tasks"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

func (c *CommandDeps) ShowAddTooltip(userID int64) {
	message := "Введите напоминание в формате команды:\n/add {дата} {время} {событие}\nНапример:\n /add 31.12.2026 18:00 Купить билеты"
	msg := tgbotapi.NewMessage(userID, message)
	c.Bot.Send(msg)
}

func (c *CommandDeps) AddTask(update *tgbotapi.Update) {
	fields := strings.Fields(update.Message.Text)
	chatID := update.Message.Chat.ID
	if len(fields) == 1 {
		c.ShowAddTooltip(update.Message.From.ID)
		return
	} else if len(fields) < 4 {
		logger.AddUserInfo(update, log.Error().Str("message", "failed to parse command").Str("command", update.Message.Text)).Send()
		message := tgbotapi.NewMessage(chatID, "Неверный формат команды")
		c.Bot.Send(message)
		return
	}

	parseLayout := `02.01.2006T15:04`
	toParse := fields[1] + "T" + fields[2]
	t, err := time.Parse(parseLayout, toParse)
	if err != nil {
		logger.AddUserInfo(update, log.Error().Str("message", "failed to parse time").Str("strToParse", toParse).Err(err)).Send()
		message := tgbotapi.NewMessage(chatID, "Ошибка формата даты/времени")
		c.Bot.Send(message)
		return
	}

	userID := update.Message.From.ID
	user, _ := c.App.Store.Users.GetByTelegramID(context.Background(), userID)
	if user == nil {
		logger.AddUserInfo(update, log.Error().Str("message", "failed to get user").Err(err)).Send()
		c.Bot.Send(tgbotapi.NewMessage(chatID, "Ошибка получения информации о пользователе"))
		return
	}
	t = t.In(time.FixedZone("Custom_time", int(time.Hour.Seconds())*int(user.UTC)))

	reminder := &store.Reminder{UserTelegramID: userID, Message: strings.Join(fields[3:], " "), SheduledTime: t}
	result, err := c.App.Store.Reminders.Create(context.Background(), reminder)
	if err != nil {
		logger.AddUserInfo(update, log.Error().Str("message", "failed to create reminder").Err(err).Any("reminder", reminder).Any("user", user)).Send()
		c.Bot.Send(tgbotapi.NewMessage(chatID, "Ошибка создания уведомления"))
		return
	}

	// broker
	//TODO time fix
	shouldReturn := sendToBroker(result, userID, update, c, chatID, t)
	if shouldReturn {
		return
	}

	logger.AddUserInfo(update, log.Info().Str("message", "reminder created").Any("reminder", reminder)).Send()
	c.Bot.Send(tgbotapi.NewMessage(chatID, "Уведомление создано!"))
}

func sendToBroker(result *store.Reminder, userID int64, update *tgbotapi.Update, c *CommandDeps, chatID int64, t time.Time) bool {
	task, err := tasks.NewReminderDeliveryTask(result.ID, userID)
	if err != nil {
		logger.AddUserInfo(update, log.Info().Str("event", "send to broker").Str("message", "could not schedule task").Err(err)).Send()
		c.App.Store.Reminders.DeleteByID(context.Background(), result.ID)
		c.Bot.Send(tgbotapi.NewMessage(chatID, "Ошибка создания уведомления"))
		return true
	}

	//TODO время не верно рассчитывается
	info, err := c.App.Broker.Enqueue(task, asynq.ProcessAt(t))
	if err != nil {
		logger.AddUserInfo(update, log.Info().Str("event", "send to broker").Str("message", "could not schedule task").Err(err)).Send()
		c.App.Store.Reminders.DeleteByID(context.Background(), result.ID)
		c.Bot.Send(tgbotapi.NewMessage(chatID, "Ошибка создания уведомления"))
		return true
	}
	logger.AddUserInfo(update, log.Info().Str("event", "sheduled task").Str("infoID", info.ID).Str("queue", info.Queue)).Send()

	// fmt.Printf("\n%s\n", t.Format("02.01.2006 15:04 -07:00"))
	return false
}
