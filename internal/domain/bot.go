package domain

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/critma/tgsheduler/internal/config"
	"github.com/critma/tgsheduler/internal/domain/commands"
	"github.com/critma/tgsheduler/internal/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func StartPoling(updates tgbotapi.UpdatesChannel, app *config.Application, quitChan <-chan os.Signal) {
	workerChan := make(chan tgbotapi.Update)

	c := commands.NewCommands(app.Bot, app)
	wg := sync.WaitGroup{}
	for i := range app.Config.TGWorkersNum {
		wg.Go(func() {
			worker(c, workerChan, i)
		})
	}

	go func() {
		for update := range updates {
			workerChan <- update
		}
		shutdown(workerChan)
	}()

	go func() {
		<-quitChan
		log.Info().Msg("Shutdown bot...")
		time.Sleep(2 * time.Second)
		shutdown(workerChan)
	}()

	wg.Wait()
}

func shutdown(workerChan chan tgbotapi.Update) {
	close(workerChan)
}

func worker(c *commands.CommandDeps, workerChan <-chan tgbotapi.Update, workerID int) {
	for update := range workerChan {
		start := time.Now()

		shouldSkip := handleRateLimier(update.FromChat().ID, c.App)
		if shouldSkip {
			observeRequest(time.Since(start), Skip, "", Skipped)
			continue
		}

		if update.Message != nil {
			logger.AddUserInfo(&update, log.Info().Str("event", "receive message")).Send()
			if update.Message.IsCommand() {
				if err := handleCommands(&update, c); err != nil {
					observeRequest(time.Since(start), Command, update.Message.Command(), Failed)
				} else {
					observeRequest(time.Since(start), Command, update.Message.Command(), Success)
				}
			} else if update.Message.Text != "" {
				if err := handleText(&update, c); err != nil {
					observeRequest(time.Since(start), Text, update.Message.Text, Failed)
				} else {
					observeRequest(time.Since(start), Text, update.Message.Text, Success)
				}
			}

		} else if update.CallbackQuery != nil {
			if err := handleCallbacks(&update, c); err != nil {
				observeRequest(time.Since(start), Callback, update.CallbackQuery.Data, Failed)
			} else {
				observeRequest(time.Since(start), Callback, update.CallbackQuery.Data, Success)
			}
		}
	}
}

func handleCommands(update *tgbotapi.Update, c *commands.CommandDeps) (err error) {
	switch commands.Command(update.Message.Command()) {
	case commands.Start:
		err = c.Start(update)
	case commands.Menu:
		err = c.ShowInlineMenu(update)
	case commands.Add:
		err = c.AddTask(update, 0)
	case commands.AddEveryday:
		err = c.AddTask(update, time.Hour*24)
	case commands.AddEveryWeek:
		err = c.AddTask(update, time.Hour*24*7)
	case commands.Timezone:
		err = c.SaveTimezone(update)
	case commands.List:
		err = c.List(update.Message.From.ID)
	case commands.Edit:
		err = c.EditTask(update)
	case commands.Delete:
		err = c.ShowDeleteList(update.Message.From.ID)
	case commands.Help:
		err = c.HandleHelp(update)
	}
	return
}

func handleText(update *tgbotapi.Update, c *commands.CommandDeps) (err error) {
	switch commands.Command(update.Message.Text) {
	case commands.Menu_ru:
		err = c.ShowInlineMenu(update)
	}
	return
}

func handleCallbacks(update *tgbotapi.Update, c *commands.CommandDeps) (err error) {
	callback := update.CallbackQuery
	log.Info().Str("event", "receive callback").Str("user", callback.From.UserName).Int64("userID", callback.From.ID).Str("callback", callback.Data).Send()

	clMessage := ""
	switch commands.Callback(callback.Data) {
	case commands.AddCallback:
		err = c.ShowAddTooltip(callback.From.ID)
		clMessage = "Добавить уведомление"
	case commands.EditCallback:
		err = c.ShowEditTooltip(callback.From.ID)
		clMessage = "Редактировать уведомление"
	case commands.DeleteCallback:
		err = c.ShowDeleteList(update.CallbackQuery.From.ID)
		clMessage = "Показать список для удаления"
	case commands.TimezoneCallback:
		err = c.ShowUserTimezone(update.CallbackQuery.From.ID)
		clMessage = "Показать UTC"
	case commands.ListCallback:
		err = c.List(callback.From.ID)
		clMessage = "Показать список уведомлений"
	}
	if strings.HasPrefix(callback.Data, string(commands.DeleteItemCallback)) {
		err = c.DeleteReminder(update)
		clMessage = "Удалить уведомление"
	}

	cl := tgbotapi.NewCallback(callback.ID, clMessage)
	_, err = c.Bot.Request(cl)
	return err
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
