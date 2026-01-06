package commands

import (
	"github.com/critma/tgsheduler/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Commands struct {
	bot *tgbotapi.BotAPI
	app *config.Application
}

func NewCommands(bot *tgbotapi.BotAPI, app *config.Application) *Commands {
	return &Commands{
		bot: bot,
		app: app,
	}
}
