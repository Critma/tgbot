package commands

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (c *Commands) HandleHelp(ChatID int64) {
	commands := `Commands:
	/add {dd.mm.yyyy} {hh:mm} {event} - add new task
	/list - show all tasks
	/edit - edit task
	/delete - delete task
	/help - show help
	`

	msg := tgbotapi.NewMessage(ChatID, commands)
	c.bot.Send(msg)
}
