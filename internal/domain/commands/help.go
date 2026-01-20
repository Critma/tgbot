package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (c *CommandDeps) HandleHelp(update *tgbotapi.Update) error {
	message := `Commands:
	/` + string(Add) + ` {dd.mm.yyyy} {hh:mm} {event} - добавить уведомление
	/` + string(List) + ` - показать ваш список уведомлений
	/` + string(Edit) + ` {id} {new event} - редактировать уведомление
	/` + string(Delete) + `  - удалить уведомление
	/` + string(Help) + `  - показать помощь
	/` + string(Timezone) + `  - показать UTC
	/` + string(Timezone) + ` {value}  - изменить UTC
	`

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	_, err := c.Bot.Send(msg)
	return err
}
