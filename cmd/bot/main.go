package main

import (
	"flag"

	"github.com/critma/tgsheduler/internal/config"
	"github.com/critma/tgsheduler/internal/domain"
	"github.com/critma/tgsheduler/internal/store/postgres"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	debug := flag.Bool("debug", false, "sets log level to debug")
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("debug level")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Error to load config")
	}
	db, err := postgres.New(cfg.DB_URL, 10, 10, "1s")
	if err != nil {
		log.Fatal().Err(err).Msg("error no open connection")
	}
	app := &config.Application{
		Config: cfg,
		DB:     db,
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TGBOT_TOKEN)
	if err != nil {
		log.Fatal().Stack().Err(err)
	}
	log.Info().Msgf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	domain.Receiver(updates, bot, app)
}
