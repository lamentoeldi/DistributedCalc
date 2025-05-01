package main

import (
	"context"
	"github.com/distributed-calc/v1/internal/orchestrator/config"
	"github.com/distributed-calc/v1/internal/orchestrator/repository/memory"
	"github.com/distributed-calc/v1/internal/orchestrator/service"
	"github.com/distributed-calc/v1/internal/orchestrator/transport/http"
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

	rep := memory.NewRepositoryMemory()
	app := service.NewService(rep, rep)

	transportCfg := &http.TransportHttpConfig{
		Host: cfg.Host,
		Port: cfg.Port,
	}

	transport := http.NewTransportHttp(app, logger, transportCfg)

	transport.Run()

	<-ctx.Done()
	transport.Shutdown(ctx)
}
