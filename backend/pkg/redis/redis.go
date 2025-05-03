package redis

import (
	"crypto/tls"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host     string `env:"REDIS_HOST"`
	Port     int    `env:"REDIS_PORT"`
	User     string `env:"REDIS_USER"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB" env-default:"0"`
	MinConns int    `env:"REDIS_MIN_CONNS" env-default:"5"`
	MaxConns int    `env:"REDIS_MAX_CONNS" env-default:"10"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read redis config: %w", err)
	}

	return &cfg, nil
}

func NewRedis(tlsConfig *tls.Config) (redis.UniversalClient, error) {
	cfg, err := NewConfig()
	if err != nil {
		return nil, fmt.Errorf("error initializing redis: %w", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:           fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Username:       cfg.User,
		Password:       cfg.Password,
		DB:             cfg.DB,
		MinIdleConns:   cfg.MinConns,
		MaxActiveConns: cfg.MaxConns,
		TLSConfig:      tlsConfig,
	})

	return client, nil
}
