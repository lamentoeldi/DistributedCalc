package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/distributed-calc/v1/internal/orchestrator/adapters/queue"
	"github.com/distributed-calc/v1/internal/orchestrator/config"
	"github.com/distributed-calc/v1/internal/orchestrator/repository/memory"
	"github.com/distributed-calc/v1/internal/orchestrator/service"
	"github.com/distributed-calc/v1/internal/orchestrator/transport/http"
	"github.com/distributed-calc/v1/pkg/authenticator"
	"github.com/distributed-calc/v1/pkg/models"
	"go.uber.org/zap"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Runtime context, cancelled on interrupt
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger, _ := zap.NewDevelopment()

	cfg, err := config.NewConfig()
	if err != nil {
		lg, _ := zap.NewProduction()
		lg.Fatal(err.Error())
	}

	q := queue.NewQueueChan[models.Task](64)

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

	rep := memory.NewRepositoryMemory()
	planner := service.NewPlannerChan(cfg, q)
	auth := authenticator.NewAuthenticator(accessPk, refreshPk, accessTTL, refreshTTL)
	app := service.NewService(rep, planner, q, auth)

	transportCfg := &http.TransportHttpConfig{
		Host: cfg.Host,
		Port: cfg.Port,
	}

	transport := http.NewTransportHttp(app, logger, transportCfg, auth)

	transport.Run()

	<-ctx.Done()
	transport.Shutdown(ctx)
}
