package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
	"time"
)

type Config struct {
	PollTimeout  time.Duration
	WorkersLimit int
	MaxRetries   int
	Url          string

	LogLevel zapcore.Level
}

func NewConfig() *Config {
	viper.AutomaticEnv()

	url := viper.GetString("master_url")
	pollTimeout := viper.GetInt("poll_timeout")
	workersLimit := viper.GetInt("computing_power")
	maxRetries := viper.GetInt("max_retries")

	level := viper.GetString("log_level")

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

	l, err := zapcore.ParseLevel(level)
	if err != nil {
		l = zapcore.InfoLevel
	}

	cfg := &Config{
		PollTimeout:  time.Duration(pollTimeout) * time.Millisecond,
		WorkersLimit: workersLimit,
		MaxRetries:   maxRetries,
		Url:          url,
		LogLevel:     l,
	}

	return cfg
}
