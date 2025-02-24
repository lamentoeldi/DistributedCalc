package config

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
	"time"
)

var (
	errInvalidTimeout      = fmt.Errorf("poll_timeout must be positive integer")
	errInvalidWorkersLimit = fmt.Errorf("computing_power must be positive integer")
	errInvalidMaxRetries   = fmt.Errorf("max_retries must be positive integer")
	errInvalidBufferSize   = fmt.Errorf("buffer_size must be positive integer")
)

type Config struct {
	PollTimeout  time.Duration
	WorkersLimit int
	MaxRetries   int
	Url          string
	BufferSize   int

	LogLevel zapcore.Level
}

func NewConfig() (*Config, error) {
	viper.AutomaticEnv()

	url := viper.GetString("master_url")
	pollTimeout := viper.GetInt("poll_timeout")
	workersLimit := viper.GetInt("computing_power")
	maxRetries := viper.GetInt("max_retries")
	bufferSize := viper.GetInt("buffer_size")

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

	if bufferSize == 0 {
		bufferSize = 128
	}

	l, err := zapcore.ParseLevel(level)
	if err != nil {
		l = zapcore.InfoLevel
	}

	if pollTimeout < 1 {
		return nil, errInvalidTimeout
	}

	if workersLimit < 1 {
		return nil, errInvalidWorkersLimit
	}

	if maxRetries < 0 {
		return nil, errInvalidMaxRetries
	}

	if bufferSize < 1 {
		return nil, errInvalidBufferSize
	}

	cfg := &Config{
		PollTimeout:  time.Duration(pollTimeout) * time.Millisecond,
		WorkersLimit: workersLimit,
		MaxRetries:   maxRetries,
		Url:          url,
		LogLevel:     l,
		BufferSize:   bufferSize,
	}

	return cfg, nil
}
