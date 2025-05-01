package mongo

import (
	"context"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

type Config struct {
	User     string `env:"MONGO_USER"`
	Password string `env:"MONGO_PASSWORD"`
	Host     string `env:"MONGO_HOST"`
	Port     int    `env:"MONGO_PORT"`
	DBName   string `env:"MONGO_NAME"`
}

func NewMongoConfig() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (mc *Config) GetDSN() string {
	var dsn strings.Builder
	dsn.WriteString("mongodb://")

	if mc.User != "" && mc.Password != "" {
		dsn.WriteString(
			fmt.Sprintf(
				"%s:%s",
				mc.User,
				mc.Password,
			),
		)
	}

	if mc.Host != "" && mc.Port != 0 {
		dsn.WriteString(
			fmt.Sprintf(
				"%s:%d/%s",
				mc.Host,
				mc.Port,
				mc.DBName,
			),
		)
	}

	return dsn.String()
}

func NewMongoClient(ctx context.Context) (*mongo.Client, error) {
	cfg, err := NewMongoConfig()
	if err != nil {
		return nil, err
	}

	dsn := cfg.GetDSN()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn))
	if err != nil {
		return nil, err
	}

	return client, nil
}
