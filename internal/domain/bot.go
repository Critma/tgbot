package domain

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/critma/tgsheduler/internal/config"
	"github.com/critma/tgsheduler/internal/domain/commands"
	"github.com/critma/tgsheduler/internal/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func StartPoling(updates tgbotapi.UpdatesChannel, app *config.Application) {
	c := commands.NewCommands(app.Bot, app, updates)
	wg := sync.WaitGroup{}
	for i := range app.Config.TGWorkersNum {
		wg.Go(func() {
			worker(c, i)
		})
	}
	wg.Wait()
}

func worker(c *commands.CommandDeps, workerID int) {
	for update := range c.Updates {
		shouldSkip := handleRateLimier(update.FromChat().ID, c.App)
		if shouldSkip {
			continue
		}

		if update.Message != nil {
			logger.AddUserInfo(&update, log.Info().Str("event", "receive message")).Send()
			if update.Message.IsCommand() {
				handleCommands(&update, c)
			} else if update.Message.Text != "" {
				handleText(&update, c)
			}

		} else if update.CallbackQuery != nil {
			handleCallbacks(&update, c)
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
		c.AddTask(update, 0)
	case commands.AddEveryday:
		c.AddTask(update, time.Hour*24)
	case commands.AddEveryWeek:
		c.AddTask(update, time.Hour*24*7)
	case commands.Timezone:
		c.SaveTimezone(update)
	case commands.List:
		c.List(update.Message.From.ID)
	case commands.Edit:
		c.EditTask(update)
	case commands.Delete:
		c.ShowDeleteList(update.Message.From.ID)
	case commands.Help:
		c.HandleHelp(update)
	}
}

func handleText(update *tgbotapi.Update, c *commands.CommandDeps) {
	switch commands.Command(update.Message.Text) {
	case commands.Menu_ru:
		c.ShowInlineMenu(update)
	}
}

func handleCallbacks(update *tgbotapi.Update, c *commands.CommandDeps) {
	callback := update.CallbackQuery
	log.Info().Str("event", "receive callback").Str("user", callback.From.UserName).Int64("userID", callback.From.ID).Str("callback", callback.Data).Send()

	clMessage := ""
	switch commands.Callback(callback.Data) {
	case commands.AddCallback:
		c.ShowAddTooltip(callback.From.ID)
		clMessage = "Добавить уведомление"
	case commands.EditCallback:
		c.ShowEditTooltip(callback.From.ID)
		clMessage = "Редактировать уведомление"
	case commands.DeleteCallback:
		c.ShowDeleteList(update.CallbackQuery.From.ID)
		clMessage = "Показать список для удаления"
	case commands.TimezoneCallback:
		c.ShowUserTimezone(update.CallbackQuery.From.ID)
		clMessage = "Показать UTC"
	case commands.ListCallback:
		c.List(callback.From.ID)
		clMessage = "Показать список уведомлений"
	}
	if strings.HasPrefix(callback.Data, string(commands.DeleteItemCallback)) {
		c.DeleteReminder(update)
		clMessage = "Удалить уведомление"
	}

	cl := tgbotapi.NewCallback(callback.ID, clMessage)
	c.Bot.Request(cl)
}

func handleRateLimier(userID int64, app *config.Application) bool {
	if !app.Config.RatelimiterEnabled {
		return false
	}

	allow, toUnlock := app.RateLimiter.Allow(userID)
	if !allow {
		log.Info().Str("event", "ratelimiter").Int64("userID", userID).Msg("rate limit exceeded")
		app.Bot.Send(tgbotapi.NewMessage(userID, fmt.Sprintf("Достигнут лимит запросов, попробуйте через %.0f секунд.", toUnlock.Seconds())))
		return true
	}

	return false
}
