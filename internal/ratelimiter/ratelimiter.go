package ratelimiter

import "time"

type Limiter interface {
	Allow(userID int64) (bool, time.Duration)
}

type Config struct {
	RequestPerTimeFrame int
	TimeFrame           time.Duration
	Enabled             bool
}
