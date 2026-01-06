package commands

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (c *Commands) ShowMenu(chatID int64) {
	message := `Hello! This is a bot for managing your tasks.`
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Add", "add"),
			tgbotapi.NewInlineKeyboardButtonData("List", "list"),
			tgbotapi.NewInlineKeyboardButtonData("Edit", "edit"),
			tgbotapi.NewInlineKeyboardButtonData("Remove", "remove"),
		),
	)
	// var numericKeyboard = tgbotapi.NewReplyKeyboard(
	// 	tgbotapi.NewKeyboardButtonRow(
	// 		tgbotapi.NewKeyboardButton("1"),
	// 		tgbotapi.NewKeyboardButton("2"),
	// 		tgbotapi.NewKeyboardButton("3"),
	// 	),
	// 	tgbotapi.NewKeyboardButtonRow(
	// 		tgbotapi.NewKeyboardButton("4"),
	// 		tgbotapi.NewKeyboardButton("5"),
	// 		tgbotapi.NewKeyboardButton("6"),
	// 	),
	// )

	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyMarkup = keyboard
	c.bot.Send(msg)
}
