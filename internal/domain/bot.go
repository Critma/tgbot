package domain

import (
	"github.com/critma/tgsheduler/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func Receiver(updates tgbotapi.UpdatesChannel, bot *tgbotapi.BotAPI, app *config.Application) {
	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Info().Str("event", "receive message").Str("user", update.Message.From.UserName).Str("text", update.Message.Text).Send()

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello!")
				bot.Send(msg)
			case "add":
				//TODO
			case "list":
				//TODO
			case "edit":
				//TODO
			case "delete":
				//TODO
			case "help":
				//TODO
			}
		}

		// if update.Message != nil {
		// 	log.Info().Msgf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		// 	msg.ReplyToMessageID = update.Message.MessageID

		// 	bot.Send(msg)
		// }
	}
}
