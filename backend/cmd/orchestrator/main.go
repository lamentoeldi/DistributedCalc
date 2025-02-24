package main

import (
	"DistributedCalc/internal/orchestrator/adapters/queue"
	"DistributedCalc/internal/orchestrator/config"
	"DistributedCalc/internal/orchestrator/repository/memory"
	"DistributedCalc/internal/orchestrator/service"
	"DistributedCalc/internal/orchestrator/transport/http"
	"DistributedCalc/pkg/models"
	"context"
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
