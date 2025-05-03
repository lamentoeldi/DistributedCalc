package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

var (
	errInvalidTimeout      = fmt.Errorf("poll_timeout must be positive integer")
	errInvalidWorkersLimit = fmt.Errorf("computing_power must be positive integer")
	errInvalidMaxRetries   = fmt.Errorf("max_retries must be positive integer")
	errInvalidBufferSize   = fmt.Errorf("buffer_size must be positive integer")
)

type Config struct {
	PollTimeout     time.Duration `env:"POLL_TIMEOUT" env-default:"100ms"`
	WorkersLimit    int           `env:"WORKERS_LIMIT" env-default:"10"`
	MaxRetries      int           `env:"MAX_RETRIES" env-default:"3"`
	OrchestratorURL string        `env:"ORCHESTRATOR_URL" env-default:"http://localhost:8080"`
	BufferSize      int           `env:"BUFFER_SIZE" env-default:"10"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	if cfg.PollTimeout < 1 {
		return nil, errInvalidTimeout
	}

	if cfg.WorkersLimit < 1 {
		return nil, errInvalidWorkersLimit
	}

	if cfg.MaxRetries < 0 {
		return nil, errInvalidMaxRetries
	}

	if cfg.BufferSize < 1 {
		return nil, errInvalidBufferSize
	}

	return &cfg, nil
}
