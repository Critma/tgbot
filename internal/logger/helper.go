package logger

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

func AddUserInfo(update *tgbotapi.Update, ev *zerolog.Event) *zerolog.Event {
	return ev.Str("user", update.FromChat().UserName).Str("text", update.Message.Text).Int64("userID", update.Message.From.ID)
}
