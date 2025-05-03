package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/distributed-calc/v1/internal/orchestrator/blacklist/redis"
	"github.com/distributed-calc/v1/internal/orchestrator/config"
	"github.com/distributed-calc/v1/internal/orchestrator/repository/mongo"
	"github.com/distributed-calc/v1/internal/orchestrator/service"
	"github.com/distributed-calc/v1/internal/orchestrator/transport/grpc"
	"github.com/distributed-calc/v1/internal/orchestrator/transport/http"
	"github.com/distributed-calc/v1/pkg/authenticator"
	mongo2 "github.com/distributed-calc/v1/pkg/mongo"
	redis2 "github.com/distributed-calc/v1/pkg/redis"
	"go.uber.org/zap"
	g "google.golang.org/grpc"
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

	mongoConfig, err := mongo2.NewMongoConfig()
	mongoClient, err := mongo2.NewMongoClient(ctx)
	if err != nil {
		logger.Fatal("failed to init mongo", zap.Error(err))
	}
	repo := mongo.NewMongoRepository(mongoConfig, mongoClient)

	redisClient, err := redis2.NewRedis(nil)
	if err != nil {
		logger.Fatal("failed to init redis", zap.Error(err))
	}
	bl := redis.NewBlacklist(redisClient)

	auth := authenticator.NewAuthenticator(accessPk, refreshPk, accessTTL, refreshTTL)

	app := service.NewService(cfg, repo, repo, repo, auth, bl)

	httpServer := http.NewServer(&http.Config{
		Host: cfg.Host,
		Port: cfg.HttpPort,
	}, app, logger)

	server := g.NewServer()

	grpcServer := grpc.NewServer(&grpc.Config{
		Host:            cfg.Host,
		GRPCPort:        cfg.GrpcPort,
		SendTaskBackoff: cfg.PollDelay,
	}, server, logger, app)

	httpServer.Run()
	grpcServer.Run()

	<-ctx.Done()
	httpServer.Shutdown(ctx)
	grpcServer.Shutdown()
}
