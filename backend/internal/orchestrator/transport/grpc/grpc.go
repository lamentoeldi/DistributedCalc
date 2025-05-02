package grpc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/distributed-calc/v1/internal/orchestrator/config"
	"github.com/distributed-calc/v1/internal/orchestrator/models"
	pb "github.com/distributed-calc/v1/pkg/proto/orchestrator"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"io"
	"net"
	"time"
)

type Service interface {
	GetTask(ctx context.Context) (*models.AgentTask, error)
	FinishTask(ctx context.Context, result *models.TaskResult) error
}

type Server struct {
	pb.UnimplementedOrchestratorServer
	cfg     *config.Config
	server  *grpc.Server
	log     *zap.Logger
	service Service
}

func NewServer(cfg *config.Config, server *grpc.Server, log *zap.Logger, service Service) *Server {
	app := &Server{
		cfg:     cfg,
		server:  server,
		log:     log,
		service: service,
	}

	pb.RegisterOrchestratorServer(server, app)

	return app
}

func (s *Server) ProcessTasks(stream grpc.BidiStreamingServer[pb.TaskResult, pb.Task]) error {
	ctx := stream.Context()

	eg, ctx := errgroup.WithContext(ctx)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	eg.Go(func() error {
		return s.sendTasks(ctx, stream)
	})

	eg.Go(func() error {
		defer cancel()
		return s.getTaskResults(ctx, stream)
	})

	err := eg.Wait()
	if err != nil {
		return fmt.Errorf("failed to process tasks: %w", err)
	}

	return nil
}

func (s *Server) sendTasks(ctx context.Context, stream grpc.BidiStreamingServer[pb.TaskResult, pb.Task]) error {
	ticker := time.NewTicker(s.cfg.PollDelay)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			task, err := s.service.GetTask(ctx)
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			if err != nil {
				s.log.Error("failed to get task", zap.Error(err))
				return fmt.Errorf("failed to get task: %w", err)
			}

			err = stream.Send(&pb.Task{
				Id:            task.Id,
				LeftArg:       task.LeftArg,
				RightArg:      task.RightArg,
				Op:            task.Op,
				OperationTime: task.OperationTime,
				Final:         task.Final,
			})
			if err != nil {
				s.log.Error("failed to send task", zap.Error(err))
				return fmt.Errorf("failed to send task: %w", err)
			}
		}
	}
}

func (s *Server) getTaskResults(ctx context.Context, stream grpc.BidiStreamingServer[pb.TaskResult, pb.Task]) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				s.log.Info("client closed the stream")
				return nil
			}

			if err != nil {
				s.log.Error("failed to receive task result", zap.Error(err))
				return fmt.Errorf("failed to receive task result: %w", err)
			}

			err = s.service.FinishTask(ctx, &models.TaskResult{
				Id:     msg.GetId(),
				Result: msg.GetResult(),
				Status: msg.GetStatus(),
				Final:  msg.GetFinal(),
			})
			if err != nil {
				s.log.Error("failed to finish task", zap.Error(err))
				return fmt.Errorf("failed to finish task: %w", err)
			}
		}
	}
}

func (s *Server) Run() {
	go func() {
		addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.GrpcPort)
		l, err := net.Listen("tcp", addr)
		if err != nil {
			s.log.Fatal("failed to listen", zap.Error(err))
		}

		err = s.server.Serve(l)
		if err != nil {
			s.log.Fatal("failed to listen", zap.Error(err))
		}
	}()
}

func (s *Server) Shutdown() {
	s.server.GracefulStop()
}
