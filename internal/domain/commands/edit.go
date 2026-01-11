package commands

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/critma/tgsheduler/internal/logger"
	"github.com/critma/tgsheduler/internal/store"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func (c *CommandDeps) ShowEditTooltip(userID int64) {
	message := "Введите изменения в формате команды:\n/edit {идентификатор} {дата} {время} {событие}\nНапример:\n /edit 2 30.12.2026 20:00 Собрание"
	msg := tgbotapi.NewMessage(userID, message)
	c.Bot.Send(msg)
}

func (c *CommandDeps) EditTask(update *tgbotapi.Update) {
	fields := strings.Fields(update.Message.Text)
	chatID := update.Message.Chat.ID
	if len(fields) < 5 {
		logger.AddUserInfo(update, log.Error().Str("message", "failed to parse command").Str("command", update.Message.Text)).Send()
		message := tgbotapi.NewMessage(chatID, "Неверный формат команды")
		c.Bot.Send(message)
		return
	}

	reminderID, err := strconv.ParseInt(fields[1], 10, 32)
	if err != nil {
		logger.AddUserInfo(update, log.Error().Str("message", "failed to parse reminderID").Str("reminderToParse", fields[1]).Err(err)).Send()
		message := tgbotapi.NewMessage(chatID, "Ошибка формата команды")
		c.Bot.Send(message)
	}

	parseLayout := `02.01.2006T15:04`
	toParse := fields[2] + "T" + fields[3]
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

	reminder := &store.Reminder{ID: int(reminderID), UserTelegramID: userID, Message: strings.Join(fields[4:], " "), SheduledTime: t, IsActive: true}
	err = c.App.Store.Reminders.Update(context.Background(), reminder)
	if err != nil {
		logger.AddUserInfo(update, log.Error().Str("message", "failed to update reminder").Err(err).Any("reminder", reminder).Any("user", user)).Send()
		c.Bot.Send(tgbotapi.NewMessage(chatID, "Ошибка изменения уведомления"))
		return
	}
	logger.AddUserInfo(update, log.Info().Str("message", "reminder updated").Any("reminder", reminder)).Send()
	c.Bot.Send(tgbotapi.NewMessage(chatID, "Изменения сохранены!"))
}
