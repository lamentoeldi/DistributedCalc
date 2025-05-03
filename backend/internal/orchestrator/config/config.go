package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

var (
	errInvalidPort      = fmt.Errorf("port must be number between 1 and 65535")
	errInvalidSleepTime = fmt.Errorf("sleep time must be positive")
)

type Config struct {
	Host     string `env:"HOST" env-default:"0.0.0.0"`
	HttpPort int    `env:"HTTP_PORT" env-default:"8080"`
	GrpcPort int    `env:"GRPC_PORT" env-default:"50051"`

	AdditionTime       time.Duration `env:"ADDITION_TIME" env-default:"1ms"`
	SubtractionTime    time.Duration `env:"SUBTRACTION_TIME" env-default:"1ms"`
	MultiplicationTime time.Duration `env:"MULTIPLICATION_TIME" env-default:"1ms"`
	DivisionTime       time.Duration `env:"DIVISION_TIME" env-default:"1ms"`

	PollDelay time.Duration `env:"POLL_DELAY" env-default:"500ms"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	if cfg.HttpPort < 1 || cfg.HttpPort > 65536 {
		return nil, errInvalidPort
	}

	if cfg.GrpcPort < 1 || cfg.GrpcPort > 65536 {
		return nil, errInvalidPort
	}

	if cfg.AdditionTime < 0 || cfg.SubtractionTime < 0 || cfg.MultiplicationTime < 0 || cfg.DivisionTime < 0 {
		return nil, errInvalidSleepTime
	}

	return &cfg, nil
}
