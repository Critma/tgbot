package commands

import (
	"github.com/critma/tgsheduler/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandDeps struct {
	Bot     *tgbotapi.BotAPI
	App     *config.Application
	Updates tgbotapi.UpdatesChannel
}

func NewCommands(bot *tgbotapi.BotAPI, app *config.Application, updates tgbotapi.UpdatesChannel) *CommandDeps {
	return &CommandDeps{
		Bot:     bot,
		App:     app,
		Updates: updates,
	}
}

type Command string

const (
	Start    Command = "start"
	Add      Command = "add"
	List     Command = "list"
	Edit     Command = "edit"
	Delete   Command = "delete"
	Menu     Command = "menu"
	Timezone Command = "tz"
	Help     Command = "help"

	Menu_ru Command = "🟢Меню🟢"

	AddEveryday  Command = "addd"
	AddEveryWeek Command = "addw"
)

type Callback string

const (
	AddCallback      Callback = "add"
	ListCallback     Callback = "list"
	EditCallback     Callback = "edit"
	DeleteCallback   Callback = "delete"
	MenuCallback     Callback = "menu"
	TimezoneCallback Callback = "tz"

	DeleteItemCallback Callback = "delete_item"
)
