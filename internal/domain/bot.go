package domain

import (
	"github.com/critma/tgsheduler/internal/config"
	"github.com/critma/tgsheduler/internal/domain/commands"
	"github.com/critma/tgsheduler/internal/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func Receiver(updates tgbotapi.UpdatesChannel, bot *tgbotapi.BotAPI, app *config.Application) {
	c := commands.NewCommands(bot, app)

	for update := range updates {
		if update.Message != nil {
			logger.AddUserInfo(&update, log.Info().Str("event", "receive message")).Send()
			if update.Message.IsCommand() {
				handleCommands(&update, c)
			} else if update.Message.Text != "" {
				// text messages
				HandleText(&update, c)
			}

		} else if update.CallbackQuery != nil {
			HandleCallbacks(&update, c)
		}
	}
}

func handleCommands(update *tgbotapi.Update, c *commands.CommandDeps) {
	switch commands.Command(update.Message.Command()) {
	case commands.Start:
		c.CreateUser(update.Message.From.ID)
		c.ShowInlineMenu(update)
		c.ShowTimezoneTooltip(update)
	case commands.Menu:
		c.ShowInlineMenu(update)
	case commands.Add:
		c.AddTask(update)
	case commands.Timezone:
		c.SaveTimezone(update)
	case "list":
		c.List(update.Message.From.ID)
	// case "edit":
	// 	//TODO
	// case "delete":
	// 	//TODO
	case "help":
		c.HandleHelp(update)
	}
}

func HandleText(update *tgbotapi.Update, c *commands.CommandDeps) {
	switch commands.Command(update.Message.Text) {
	case commands.Menu_ru:
		c.ShowInlineMenu(update)
	}
}

func HandleCallbacks(update *tgbotapi.Update, c *commands.CommandDeps) {
	callback := update.CallbackQuery
	log.Info().Str("event", "receive callback").Str("user", callback.From.UserName).Int64("userID", callback.From.ID).Str("callback", callback.Data).Send()

	switch commands.Callback(callback.Data) {
	case commands.AddCallback:
		c.ShowAddTooltip(update)
	case commands.TimezoneCallback:
		c.ShowUserTimezone(update.CallbackQuery.From.ID)
	case commands.ListCallback:
		c.List(callback.From.ID)
	}
}
