package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func (c *CommandDeps) ShowDeleteList(userID int64) {
	reminders, err := c.App.Store.Reminders.GetByUserID(context.Background(), userID)
	if err != nil {
		log.Error().Str("message", "failed to get reminders").Err(err).Int64("userID", userID).Send()
		c.Bot.Send(tgbotapi.NewMessage(userID, "Ошибка получения напоминаний"))
		return
	}

	buttons := make([][]tgbotapi.InlineKeyboardButton, len(reminders))
	for i, rem := range reminders {
		msg := fmt.Sprintf("%s %d", DeleteItemCallback, rem.ID)
		btn := tgbotapi.InlineKeyboardButton{
			Text:         fmt.Sprintf("❌ Удалить '%s'", rem.Message),
			CallbackData: &msg,
		}
		buttons[i] = []tgbotapi.InlineKeyboardButton{btn}
	}
	keyboard := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}

	sb := &strings.Builder{}
	if len(reminders) == 0 {
		sb.WriteString(`У вас нет напоминаний!`)
	} else {
		sb.WriteString("Список ваших напоминаний:")
		for i, rem := range reminders {
			fmt.Fprintf(sb, "\n%d. %s (%v)", i+1, rem.Message, rem.SheduledTime.Format("02.01.2006 15:04"))
		}
	}

	msg := tgbotapi.NewMessage(userID, sb.String())
	msg.ReplyMarkup = keyboard
	_, err = c.Bot.Send(msg)
	if err != nil {
		log.Error().Str("message", "failed to send listToDelete").Err(err).Int64("userID", userID).Send()
	}
}

func (c *CommandDeps) DeleteReminder(update *tgbotapi.Update) {
	fields := strings.Fields(update.CallbackQuery.Data)

	if len(fields) != 2 {
		log.Error().Str("message", "failed format delete callback").Any("fields", fields).Int64("userID", update.CallbackQuery.From.ID).Send()
		c.Bot.Send(tgbotapi.NewMessage(update.CallbackQuery.From.ID, "Ошибка удаления напоминания"))
		return
	}

	reminderID, err := strconv.ParseInt(fields[1], 10, 32)
	if err != nil {
		log.Error().Str("message", "failed to parse format delete callback").Any("Parsenumber", fields[1]).Err(err).Int64("userID", update.CallbackQuery.From.ID).Send()
		c.Bot.Send(tgbotapi.NewMessage(update.CallbackQuery.From.ID, "Ошибка форматы команды"))
		return
	}

	err = c.App.Store.Reminders.DeleteByID(context.Background(), int(reminderID))
	if err != nil {
		log.Error().Str("message", "failed to delete reminder").Int64("reminderID", reminderID).Err(err).Int64("userID", update.CallbackQuery.From.ID).Send()
		c.Bot.Send(tgbotapi.NewMessage(update.CallbackQuery.From.ID, "Ошибка удаления напоминания"))
		return
	}
	log.Info().Str("message", "delete reminder").Int64("reminderId", reminderID).Send()
	c.Bot.Send(tgbotapi.NewMessage(update.CallbackQuery.From.ID, "Удаление успешно!"))
}
