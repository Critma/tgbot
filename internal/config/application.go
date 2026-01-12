package config

import (
	"github.com/critma/tgsheduler/internal/store"
	"github.com/hibiken/asynq"
)

type Application struct {
	Config *Config
	Store  store.Storage
	Broker *asynq.Client
}
