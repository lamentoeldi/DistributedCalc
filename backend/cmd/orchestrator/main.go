package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	memory2 "github.com/distributed-calc/v1/internal/orchestrator/blacklist/memory"
	"github.com/distributed-calc/v1/internal/orchestrator/config"
	"github.com/distributed-calc/v1/internal/orchestrator/repository/memory"
	"github.com/distributed-calc/v1/internal/orchestrator/service"
	"github.com/distributed-calc/v1/internal/orchestrator/transport/http"
	"github.com/distributed-calc/v1/pkg/authenticator"
	"go.uber.org/zap"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Runtime context, cancelled on interrupt
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger, _ := zap.NewProduction()

	cfg, err := config.NewConfig()
	if err != nil {
	}

	// TODO: read from config
	accessPk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		logger.Fatal("failed to start", zap.Error(err))
	}

	refreshPk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		logger.Fatal("failed to start", zap.Error(err))
	}

	accessTTL := 10 * time.Minute
	refreshTTL := 7 * 24 * time.Hour

	repo := memory.NewRepositoryMemory()
	auth := authenticator.NewAuthenticator(accessPk, refreshPk, accessTTL, refreshTTL)
	bl := memory2.NewBlacklist()

	app := service.NewService(repo, repo, repo, auth, bl)

	transportCfg := &http.TransportHttpConfig{
		Host: cfg.Host,
		Port: cfg.Port,
	}

	transport := http.NewTransportHttp(app, logger, transportCfg)

	transport.Run()

	<-ctx.Done()
	transport.Shutdown(ctx)
}
