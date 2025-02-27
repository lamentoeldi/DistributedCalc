package main

import (
	"DistributedCalc/internal/agent/adapters/orchestrator"
	"DistributedCalc/internal/agent/config"
	"DistributedCalc/internal/agent/service"
	"DistributedCalc/internal/agent/transport/async"
	"context"
	"go.uber.org/zap"
	"log"
	"net/http"
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
	defer logger.Sync()

	api := orchestrator.NewOrchestrator(http.DefaultClient, cfg.Url, cfg.MaxRetries)
	err = api.Ping()
	if err != nil {
		logger.Fatal("Failed to connect to orchestrator", zap.Error(err))
	}

	app := service.NewService()

	transport := async.NewTransportAsync(cfg, logger, api, app)
	transport.Run(ctx)

	// Graceful stop
	<-ctx.Done()
	transport.Shutdown()
}
