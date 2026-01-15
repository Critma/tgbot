package config

import (
	"database/sql"

	"github.com/critma/tgsheduler/internal/ratelimiter"
	"github.com/critma/tgsheduler/internal/store"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
)

type Application struct {
	Config      *Config
	Store       store.Storage
	Db          *sql.DB
	Broker      *Broker
	RateLimiter ratelimiter.Limiter
	Bot         *tgbotapi.BotAPI
}

type Broker struct {
	Client    *asynq.Client
	Inspector *asynq.Inspector
}
