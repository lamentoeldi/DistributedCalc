package main

import (
	"context"
	"fmt"
	"github.com/distributed-calc/v1/internal/agent/config"
	"github.com/distributed-calc/v1/internal/agent/service"
	"github.com/distributed-calc/v1/internal/agent/transport/grpc"
	"go.uber.org/zap"
	g "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os/signal"
	"syscall"
)

func main() {
	// Runtime context, cancelled on interrupt
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatal("error loading config", zap.Error(err))
	}

	app := service.NewService()

	creds := g.WithTransportCredentials(insecure.NewCredentials())

	addr := fmt.Sprintf("%s:%d", cfg.OrchestratorHost, cfg.OrchestratorPort)
	client, err := g.NewClient(addr, creds)
	if err != nil {
		logger.Fatal("error creating gRPC client", zap.Error(err))
	}

	server := grpc.NewServer(cfg, client, app)

	err = server.Run(ctx)
	if err != nil {
		logger.Fatal("error starting server", zap.Error(err))
	}

	// Graceful stop
	<-ctx.Done()
}
