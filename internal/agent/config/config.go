package config

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	PollTimeout  time.Duration
	WorkersLimit int
	MaxRetries   int
	Url          string
}

func NewConfig() *Config {
	viper.AutomaticEnv()

	url := viper.GetString("master_url")
	pollTimeout := viper.GetInt("poll_timeout")
	workersLimit := viper.GetInt("computing_power")
	maxRetries := viper.GetInt("max_retries")

	if url == "" {
		url = "http://localhost:8080"
	}

	if pollTimeout == 0 {
		pollTimeout = 50
	}

	if workersLimit == 0 {
		workersLimit = 10
	}

	if maxRetries == 0 {
		maxRetries = 3
	}

	cfg := &Config{
		PollTimeout:  time.Duration(pollTimeout) * time.Millisecond,
		WorkersLimit: workersLimit,
		MaxRetries:   maxRetries,
		Url:          url,
	}

	return cfg
}
