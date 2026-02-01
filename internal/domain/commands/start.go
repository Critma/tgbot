package commands

import (
	"context"

	"github.com/critma/tgsheduler/internal/store"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (c *CommandDeps) createUser(userID int64) error {
	user := &store.User{TelegramID: userID, UTC: 3}

	return c.App.Store.Users.Create(context.Background(), user)
}

func (c *CommandDeps) Start(update *tgbotapi.Update) error {
	if err := c.createUser(update.Message.From.ID); err != nil {
		return err
	}
	if err := c.ShowInlineMenu(update); err != nil {
		return err
	}
	if err := c.ShowTimezoneTooltip(update); err != nil {
		return err
	}
	return nil
}
