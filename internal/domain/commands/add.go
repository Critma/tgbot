package commands

import (
	"context"
	"strings"
	"time"

	"github.com/critma/tgsheduler/internal/store"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func (c *Commands) ShowAddTooltip(update *tgbotapi.Update) {
	message := "Введите напоминание в формате:\n\n/add дата время событие\nНапример: /add 31.12.2026 18:00 Купить билеты"
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, message)
	c.bot.Send(msg)
}

func (c *Commands) AddTask(update *tgbotapi.Update) {
	fields := strings.Fields(update.Message.Text)
	chatID := update.Message.Chat.ID
	if len(fields) < 4 {
		log.Error().Str("message", "failed to parse command").Str("command", update.Message.Text).Send()
		message := tgbotapi.NewMessage(chatID, "Неверный формат команды")
		c.bot.Send(message)
		return
	}

	parseLayout := `02.01.2006T15:04`
	toParse := fields[1] + "T" + fields[2]
	//TODO: timezone
	t, err := time.Parse(parseLayout, toParse)
	if err != nil {
		log.Error().Str("message", "failed to parse time").Str("strToParse", toParse).Err(err).Send()
		message := tgbotapi.NewMessage(chatID, "Ошибка формата даты/времени")
		c.bot.Send(message)
		return
	}

	userID := update.Message.From.ID
	update.Message.Time()
	_ = c.app.Store.Users.Create(context.Background(), userID)

	reminder := &store.Reminder{UserTelegramID: userID, Message: strings.Join(fields[3:], " "), SheduledTime: t}
	err = c.app.Store.Reminders.Create(context.Background(), reminder)
	if err != nil {
		log.Error().Str("message", "failed to create reminder").Err(err).Any("object", reminder).Send()
		message := tgbotapi.NewMessage(chatID, "Ошибка создания напоминания")
		c.bot.Send(message)
		return
	}
	log.Info().Str("message", "reminder created").Any("reminder", reminder).Send()
	message := tgbotapi.NewMessage(chatID, "Напоминание создано")
	c.bot.Send(message)
}
