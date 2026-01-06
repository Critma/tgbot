package config

import (
	"github.com/critma/tgsheduler/internal/store"
)

type Application struct {
	Config *Config
	Store  store.Storage
}
