package mongo

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ilyakaznacheev/cleanenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

type Config struct {
	User           string `env:"MONGO_USER"`
	Password       string `env:"MONGO_PASSWORD"`
	Host           string `env:"MONGO_HOST"`
	Port           int    `env:"MONGO_PORT"`
	DBName         string `env:"MONGO_NAME"`
	MigrationsPath string `env:"MONGO_MIGRATIONS_PATH"`
	AuthSource     string `env:"MONGO_AUTH_SOURCE" env-default:"admin"`
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
				"@%s:%d/%s",
				mc.Host,
				mc.Port,
				mc.DBName,
			),
		)
	}

	if mc.AuthSource != "" {
		authSource := fmt.Sprintf("?authSource=%s", mc.AuthSource)
		dsn.WriteString(authSource)
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

	if cfg.MigrationsPath != "" {
		migrationsPath := fmt.Sprintf("file://%s", cfg.MigrationsPath)
		migr, err := migrate.New(migrationsPath, dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to init migrations: %w", err)
		}

		err = migr.Up()
		if err != nil {
			return nil, fmt.Errorf("failed to run migrations: %w", err)
		}
	}

	return client, nil
}
