package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/critma/tgsheduler/internal/domain/helpers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func (c *CommandDeps) List(userID int64) {
	reminders, err := c.App.Store.Reminders.GetActiveByUserID(context.Background(), userID)
	if err != nil {
		log.Error().Str("message", "failed to get reminders").Err(err).Int64("userID", userID).Send()
		c.Bot.Send(tgbotapi.NewMessage(userID, "Ошибка получения напоминаний"))
		return
	}
	log.Info().Str("message", "reminders listed").Any("reminders", reminders).Int64("userID", userID).Send()
	var sb strings.Builder
	if len(reminders) == 0 {
		sb.WriteString("У вас нет запланированных уведомлений!")
	} else {
		user, err := c.App.Store.Users.GetByTelegramID(context.Background(), userID)
		if err != nil {
			log.Error().Str("message", "failed to get user timezone").Err(err).Int64("userID", userID).Send()
			c.Bot.Send(tgbotapi.NewMessage(userID, "Ошибка получения информации о часовом поясе пользователя"))
		}

		sb.WriteString("Ваши запланированные уведомления (id дата время описание):\n")
		for _, reminder := range reminders {
			reminder.SheduledTime = helpers.TimeToUserTZ(user, reminder.SheduledTime)
			if reminder.RepeatInterval.Hours() == 24 {
				fmt.Fprintf(&sb, "%v🔸 ежедневно %s  👉 %s\n", reminder.ID, reminder.SheduledTime.Format("15:04"), reminder.Message)
			} else if reminder.RepeatInterval.Hours() == 24*7 {
				fmt.Fprintf(&sb, "%v🔸 еженедельно %s  👉 %s\n", reminder.ID, reminder.SheduledTime.Format("15:04"), reminder.Message)
			} else {
				fmt.Fprintf(&sb, "%v🔸 %s  👉 %s\n", reminder.ID, reminder.SheduledTime.Format("02.01.2006 15:04"), reminder.Message)
			}
		}
	}
	message := tgbotapi.NewMessage(userID, sb.String())
	c.Bot.Send(message)
}
