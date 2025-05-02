package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/distributed-calc/v1/internal/agent/config"
	"github.com/distributed-calc/v1/internal/agent/models"
	pb "github.com/distributed-calc/v1/pkg/proto/orchestrator"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"io"
	"sync"
	"time"
)

type Service interface {
	Evaluate(task *models.AgentTask) (*models.TaskResult, error)
}

type Server struct {
	cfg     *config.Config
	client  pb.OrchestratorClient
	in      chan *models.AgentTask
	out     chan *models.TaskResult
	service Service
}

func NewServer(cfg *config.Config, client *grpc.ClientConn, service Service) *Server {
	return &Server{
		cfg:     cfg,
		client:  pb.NewOrchestratorClient(client),
		in:      make(chan *models.AgentTask, cfg.BufferSize),
		out:     make(chan *models.TaskResult, cfg.BufferSize),
		service: service,
	}
}

func (s *Server) Run(ctx context.Context) error {
	stream, err := s.client.ProcessTasks(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to orchestrator: %w", err)
	}

	ctx = stream.Context()

	eg, ctx := errgroup.WithContext(ctx)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	eg.Go(func() error {
		return s.getTasks(ctx, stream)
	})

	eg.Go(func() error {
		return s.sendTaskResults(ctx, stream)
	})

	go s.runWorkers(ctx)

	err = eg.Wait()
	if err != nil {
		return fmt.Errorf("task processing finished with error: %w", err)
	}

	return nil
}

func (s *Server) getTasks(ctx context.Context, stream grpc.BidiStreamingClient[pb.TaskResult, pb.Task]) error {
	defer close(s.in)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				return nil
			}

			if err != nil {
				return fmt.Errorf("failed to receive task: %w", err)
			}

			s.in <- &models.AgentTask{
				Id:            msg.GetId(),
				LeftArg:       msg.GetLeftArg(),
				RightArg:      msg.GetRightArg(),
				Op:            msg.GetOp(),
				OperationTime: msg.GetOperationTime(),
				Final:         msg.GetFinal(),
			}
		}
	}
}

func (s *Server) sendTaskResults(ctx context.Context, stream grpc.BidiStreamingClient[pb.TaskResult, pb.Task]) error {
	ticker := time.NewTicker(s.cfg.PollTimeout)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			task, ok := <-s.out
			if !ok {
				return nil
			}

			err := stream.Send(&pb.TaskResult{
				Id:     task.Id,
				Result: task.Result,
				Status: task.Status,
				Final:  task.Final,
			})
			if err != nil {
				return fmt.Errorf("failed to send task result: %w", err)
			}
		}
	}
}

func (s *Server) runWorkers(ctx context.Context) {
	defer close(s.out)

	wg := &sync.WaitGroup{}

	wg.Add(s.cfg.WorkersLimit)
	for range s.cfg.WorkersLimit {
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case task, ok := <-s.in:
					if !ok {
						return
					}

					result, _ := s.service.Evaluate(task)
					s.out <- result
				}
			}
		}()
	}

	wg.Wait()
}
