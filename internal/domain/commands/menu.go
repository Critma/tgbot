package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (c *CommandDeps) ShowInlineMenu(update *tgbotapi.Update) error {
	message := `Привет, это бот для отложенных уведомлений`
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Add", string(AddCallback)),
			tgbotapi.NewInlineKeyboardButtonData("List", string(ListCallback)),
			tgbotapi.NewInlineKeyboardButtonData("Edit", string(EditCallback)),
			tgbotapi.NewInlineKeyboardButtonData("Remove", string(DeleteCallback)),
			tgbotapi.NewInlineKeyboardButtonData("Timezone", string(TimezoneCallback)),
		),
	)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	msg.ReplyMarkup = keyboard
	_, err := c.Bot.Send(msg)
	return err
}

func (c *CommandDeps) GetKeyboardMenu() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(string(Menu_ru)),
		),
	)
	return keyboard
}
