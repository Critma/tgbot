package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/critma/tgsheduler/internal/config"
	"github.com/critma/tgsheduler/internal/domain"
	"github.com/critma/tgsheduler/internal/metrics"
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
		Db:          db,
		Broker:      &config.Broker{Client: client, Inspector: inspector},
		RateLimiter: rateLimiter,
		Bot:         bot,
	}

	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, syscall.SIGTERM, syscall.SIGINT)

	// asynq workers
	go tasks.StartAsynqWorkers(bot, app)

	// metrics
	go func() {
		log.Info().Msg("Start metrics server on " + cfg.MetricsAddr)
		err := metrics.Listen(cfg.MetricsAddr)
		if err != nil {
			log.Error().Str("event", "start metrics server").Str("message", "server not started").Err(err).Send()
		}
	}()

	// start bot
	domain.StartPoling(updates, app, quitChan)
}
