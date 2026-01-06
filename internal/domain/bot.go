package domain

import (
	"github.com/critma/tgsheduler/internal/config"
	"github.com/critma/tgsheduler/internal/domain/commands"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func Receiver(updates tgbotapi.UpdatesChannel, bot *tgbotapi.BotAPI, app *config.Application) {
	c := commands.NewCommands(bot, app)

	for update := range updates {
		if update.Message != nil {
			log.Info().Str("event", "receive message").Str("user", update.Message.From.UserName).Str("text", update.Message.Text).Send()
			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "start":
					c.ShowMenu(update.Message.Chat.ID)
				case "add":
					c.AddTask(&update)
				// case "list":
				// 	//TODO
				// case "edit":
				// 	//TODO
				// case "delete":
				// 	//TODO
				case "help":
					c.HandleHelp(update.Message.Chat.ID)
				}
			}

		} else if update.CallbackQuery != nil {
			callback := update.CallbackQuery
			log.Info().Str("event", "receive callback").Str("user", callback.From.UserName).Str("callback", callback.Data).Send()

			switch callback.Data {
			case "add":
				c.ShowAddTooltip(&update)
			}
		}
	}
}
