package main

import (
	"DistributedCalc/internal/orchestrator/config"
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

	cfg := config.NewConfig()

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}

	queue := service.NewQueue[models.Task](64)

	transportCfg := &http.TransportHttpConfig{
		Host: cfg.Host,
		Port: cfg.Port,
	}

	transport := http.NewTransportHttp(nil, logger, transportCfg, queue)

	transport.Run()

	<-ctx.Done()
	transport.Shutdown(ctx)
}
