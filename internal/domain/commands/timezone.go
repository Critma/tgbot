package commands

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/critma/tgsheduler/internal/logger"
	"github.com/critma/tgsheduler/internal/store"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func (c *CommandDeps) ShowTimezoneTooltip(update *tgbotapi.Update) error {
	message := fmt.Sprintf("Выберите своё UTC время, по умолчанию +3 (Московское время)\nКоманда: /%s +5", string(Timezone))
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	msg.ReplyMarkup = c.GetKeyboardMenu()
	_, err := c.Bot.Send(msg)
	return err
}

func (c *CommandDeps) SaveTimezone(update *tgbotapi.Update) error {
	fields := strings.Fields(update.Message.Text)
	if len(fields) == 1 {
		c.ShowUserTimezone(update.Message.From.ID)
		return nil
	} else if len(fields) != 2 {
		logger.AddUserInfo(update, log.Error().Str("message", "failed to parse timezone").Any("fields", fields).Str("error", "fields len != 2")).Send()
		c.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка формата команды"))
		return errors.New("Ошибка формата команды")
	}

	utc, err := strconv.ParseInt(fields[1], 10, 8)
	if err != nil {
		logger.AddUserInfo(update, log.Error().Str("message", "failed to parse timezone").Str("strToParse", fields[1]).Err(err)).Send()
		c.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка формата utc"))
		return errors.New("Ошибка формата utc")
	}
	if !c.IsCorrectUTC(int8(utc)) {
		logger.AddUserInfo(update, log.Error().Str("message", "failed to save timezone").Str("timezone", fields[1]).Str("error", "timezone should be -12 - +14")).Send()
		c.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка формата часового пояса, UTC должен быть от -12 до +14"))
		return errors.New("Ошибка формата часового пояса, UTC должен быть от -12 до +14")
	}

	userID := update.Message.From.ID
	user := &store.User{TelegramID: userID, UTC: int8(utc)}

	err = c.App.Store.Users.CreateOrUpdate(context.Background(), user)
	var tgMessage string
	if err != nil {
		logger.AddUserInfo(update, log.Error().Str("message", "failed to save timezone").Str("timezone", fields[1]).Err(err)).Send()
		tgMessage = "Ошибка сохранения часового пояса"
	} else {
		logger.AddUserInfo(update, log.Info().Str("message", "timezone saved").Str("timezone", fields[1])).Send()
		tgMessage = "Часовой пояс сохранён!"
	}
	_, err = c.Bot.Send(tgbotapi.NewMessage(userID, tgMessage))
	return err
}

func (c *CommandDeps) ShowUserTimezone(userID int64) error {
	user, err := c.App.Store.Users.GetByTelegramID(context.Background(), userID)
	if err != nil {
		log.Error().Str("message", "get-user error").Err(err).Str("userID", "userID").Send()
		c.Bot.Send(tgbotapi.NewMessage(userID, "Ошибка получения информации о пользователе"))
		return errors.New("ошибка получения информации о пользователе")
	}
	sign := ""
	if user.UTC < 0 {
		sign = "-"
	} else if user.UTC > 0 {
		sign = "+"
	}
	msg := fmt.Sprintf("🕛Ваш UTC : %s%v🕛\nДля смены: /%s {значение}\nНапример: /%s +5", sign, user.UTC, string(Timezone), string(Timezone))
	_, err = c.Bot.Send(tgbotapi.NewMessage(userID, msg))
	return err
}

func (c *CommandDeps) IsCorrectUTC(utc int8) bool {
	return -12 <= utc && utc <= 14
}
