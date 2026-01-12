package commands

import (
	"context"
	"strconv"
	"strings"

	"github.com/critma/tgsheduler/internal/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func (c *CommandDeps) ShowEditTooltip(userID int64) {
	message := "Введите изменения в формате команды:\n/edit {идентификатор} {новое название}\nНапример:\n /edit 2 Собрание"
	msg := tgbotapi.NewMessage(userID, message)
	c.Bot.Send(msg)
}

func (c *CommandDeps) EditTask(update *tgbotapi.Update) {
	fields := strings.Fields(update.Message.Text)
	chatID := update.Message.Chat.ID
	if len(fields) < 3 {
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

	newMessage := strings.Join(fields[2:], " ")
	err = c.App.Store.Reminders.UpdateMessage(context.Background(), int(reminderID), newMessage)
	if err != nil {
		logger.AddUserInfo(update, log.Error().Str("message", "failed to update message in reminder").Err(err)).Send()
		c.Bot.Send(tgbotapi.NewMessage(chatID, "Ошибка изменения уведомления"))
		return
	}

	logger.AddUserInfo(update, log.Info().Str("message", "reminder updated").Str("newMessage", newMessage)).Send()
	c.Bot.Send(tgbotapi.NewMessage(chatID, "Изменения сохранены!"))
}
