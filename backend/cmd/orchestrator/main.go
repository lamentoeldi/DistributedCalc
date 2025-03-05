package main

import (
	"context"
	"github.com/distributed-calc/v1/internal/orchestrator/adapters/queue"
	"github.com/distributed-calc/v1/internal/orchestrator/config"
	"github.com/distributed-calc/v1/internal/orchestrator/repository/memory"
	"github.com/distributed-calc/v1/internal/orchestrator/service"
	"github.com/distributed-calc/v1/internal/orchestrator/transport/http"
	"github.com/distributed-calc/v1/pkg/models"
	"go.uber.org/zap"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	// Runtime context, cancelled on interrupt
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.NewConfig()
	if err != nil {
		lg, _ := zap.NewProduction()
		lg.Fatal(err.Error())
	}

	logConfig := zap.NewProductionConfig()
	logConfig.Level = zap.NewAtomicLevelAt(cfg.LogLevel)

	logger, err := logConfig.Build()
	if err != nil {
		log.Fatal(err)
	}

	q := queue.NewQueueChan[models.Task](64)

	rep := memory.NewRepositoryMemory()
	planner := service.NewPlannerChan(cfg, q)
	app := service.NewService(rep, planner, q)

	transportCfg := &http.TransportHttpConfig{
		Host: cfg.Host,
		Port: cfg.Port,
	}

	transport := http.NewTransportHttp(app, logger, transportCfg)

	transport.Run()

	<-ctx.Done()
	transport.Shutdown(ctx)
}
