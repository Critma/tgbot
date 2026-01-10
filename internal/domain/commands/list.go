package commands

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func (c *CommandDeps) List(userID int64) {
	reminders, err := c.app.Store.Reminders.GetByUserID(context.Background(), userID)
	if err != nil {
		log.Error().Str("message", "failed to get reminders").Err(err).Int64("userID", userID).Send()
		c.bot.Send(tgbotapi.NewMessage(userID, "Ошибка получения напоминаний"))
		return
	}
	log.Info().Str("message", "reminders listed").Any("reminders", reminders).Int64("userID", userID).Send()
	var sb strings.Builder
	for _, reminder := range reminders {
		// localTime := reminder.SheduledTime.In(loc)
		fmt.Fprintf(&sb, "🔸 %s  👉 %s\n", reminder.SheduledTime.Format("02.01.2006 15:04"), reminder.Message)
	}
	message := tgbotapi.NewMessage(userID, sb.String())
	c.bot.Send(message)
}
