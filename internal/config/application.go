package config

import (
	"database/sql"

	"github.com/rs/zerolog"
)

type Application struct {
	Config *Config
	DB     *sql.DB
	logger *zerolog.Logger
}
