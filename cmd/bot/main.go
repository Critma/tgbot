package main

import (
	"flag"
	"time"

	"github.com/critma/tgsheduler/internal/config"
	"github.com/critma/tgsheduler/internal/domain"
	"github.com/critma/tgsheduler/internal/ratelimiter"
	"github.com/critma/tgsheduler/internal/store/postgres"
	"github.com/critma/tgsheduler/internal/tasks"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	debug := flag.Bool("debug", false, "sets log level to debug")
	flag.Parse()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("debug level")
	}

	// cfg
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Error to load config")
	}
	db, err := postgres.New(cfg.DB_URL, 10, 10, "1s")
	if err != nil {
		log.Fatal().Err(err).Msg("error no open connection")
	}

	// asynq client
	clientOpt := asynq.RedisClientOpt{Addr: cfg.REDIS_URL, DB: 0}
	client := asynq.NewClient(clientOpt)
	defer client.Close()

	inspector := asynq.NewInspector(clientOpt)
	defer inspector.Close()

	//ratelimiter
	rateLimiter := ratelimiter.NewFixedWindowRateLimiter(cfg.RatelimiterRequests, time.Second*time.Duration(cfg.RatelimiterTimeFrameSeconds))

	// tgbot
	bot, err := tgbotapi.NewBotAPI(cfg.TGBOT_TOKEN)
	if err != nil {
		log.Fatal().Stack().Err(err)
	}
	log.Info().Str("start", "authorized on account "+bot.Self.UserName).Send()
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	app := &config.Application{
		Config:      cfg,
		Store:       *postgres.NewStorage(db),
		Broker:      &config.Broker{Client: client, Inspector: inspector},
		RateLimiter: rateLimiter,
		Bot:         bot,
	}

	//workers
	go startWorkers(cfg, bot, app)

	//start bot
	domain.Receiver(updates, app)
}

func startWorkers(cfg *config.Config, bot *tgbotapi.BotAPI, app *config.Application) {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.REDIS_URL},
		asynq.Config{
			Concurrency: 5,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)
	mux := asynq.NewServeMux()
	mux.Handle(tasks.TypeReminderDelivery, tasks.NewReminderProcessor(bot, app))
	if err := srv.Run(mux); err != nil {
		log.Fatal().Msgf("could not run workers-server: %v", err)
	}
}
