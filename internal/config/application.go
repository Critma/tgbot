package config

import (
	"github.com/critma/tgsheduler/internal/store"
	"github.com/hibiken/asynq"
)

type Application struct {
	Config *Config
	Store  store.Storage
	Broker *Broker
}

type Broker struct {
	Client    *asynq.Client
	Inspector *asynq.Inspector
}
