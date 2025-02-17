package config

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	PollTimeout  time.Duration
	WorkersLimit int
	Url          string
}

func NewConfig() *Config {
	viper.AutomaticEnv()

	pollTimeout := viper.GetInt("poll_timeout")
	workersLimit := viper.GetInt("computing_power")

	cfg := &Config{
		PollTimeout:  time.Duration(pollTimeout) * time.Millisecond,
		WorkersLimit: workersLimit,
	}

	return cfg
}
