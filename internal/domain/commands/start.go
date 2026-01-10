package commands

import (
	"context"

	"github.com/critma/tgsheduler/internal/store"
)

func (c *CommandDeps) CreateUser(userID int64) {
	user := &store.User{TelegramID: userID, UTC: 3}

	_ = c.app.Store.Users.Create(context.Background(), user)
}
