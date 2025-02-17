package main

import (
	"DistributedCalc/internal/agent/adapters/orchestrator"
	"DistributedCalc/internal/agent/config"
	"DistributedCalc/internal/agent/service"
	"DistributedCalc/internal/agent/transport/async"
	"DistributedCalc/pkg/models"
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

	cfg := config.NewConfig()

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}

	api := orchestrator.NewOrchestrator(http.DefaultClient, cfg.Url)

	app := service.NewService()

	in := make(chan *models.Task, 64)
	out := make(chan *models.TaskResult, 64)

	transport := async.NewTransportAsync(cfg, logger, api, app, in, out)
	transport.Run(ctx)

	// Graceful stop
	<-ctx.Done()
	transport.Shutdown()
}
