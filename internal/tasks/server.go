package tasks

import (
	"github.com/critma/tgsheduler/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

func StartAsynqWorkers(bot *tgbotapi.BotAPI, app *config.Application) {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: app.Config.REDIS_URL},
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
	mux.Handle(TypeReminderDelivery, NewReminderProcessor(bot, app))
	if err := srv.Run(mux); err != nil {
		log.Fatal().Msgf("could not run workers-server: %v", err)
	}
}
