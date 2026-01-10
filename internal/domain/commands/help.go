package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (c *CommandDeps) HandleHelp(update *tgbotapi.Update) {
	message := `Commands:
	/` + string(Add) + ` {dd.mm.yyyy} {hh:mm} {event} - add new task
	/` + string(List) + ` - show all tasks
	/` + string(Edit) + `  - edit task
	/` + string(Delete) + `  - delete task
	/` + string(Help) + `  - show help
	`

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	c.bot.Send(msg)
}
