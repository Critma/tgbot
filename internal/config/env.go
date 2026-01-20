package config

import (
	"log"

	"github.com/spf13/viper"
)

const (
	ENV_PATH = "configs"
)

type Config struct {
	DBHost         string `mapstructure:"POSTGRES_HOST"`
	DBUserName     string `mapstructure:"POSTGRES_USER"`
	DBUserPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName         string `mapstructure:"POSTGRES_DB"`
	DBPort         string `mapstructure:"POSTGRES_PORT"`
	DB_URL         string `mapstructure:"POSTGRES_URL"`

	REDIS_URL string `mapstructure:"REDIS_URL"`

	RatelimiterEnabled          bool `mapstructure:"RATELIMITER_ENABLED"`
	RatelimiterRequests         int  `mapstructure:"RATELIMITER_REQUESTS"`
	RatelimiterTimeFrameSeconds int  `mapstructure:"RATELIMITER_TIMEFRAME_SECONDS"`

	MetricsAddr string `mapstructure:"METRICS_ADDR"`

	TGBOT_TOKEN string `mapstructure:"BOT_TOKEN"`

	TGWorkersNum int `mapstructure:"NUM_WORKER_POOL"`
}

func LoadConfig() (c *Config, err error) {
	viper.AddConfigPath(ENV_PATH)
	viper.AddConfigPath("../../configs")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		return
	}

	log.Printf("Config loaded: %#v\n", c)
	return
}
