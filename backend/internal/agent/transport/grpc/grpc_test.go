package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/distributed-calc/v1/internal/agent/config"
	"github.com/distributed-calc/v1/internal/agent/models"
	"github.com/distributed-calc/v1/internal/agent/service"
	pb "github.com/distributed-calc/v1/pkg/proto/orchestrator"
	"github.com/distributed-calc/v1/test/mock"
	"io"
	"testing"
)

func TestGetTasks(t *testing.T) {
	svc := service.NewService()

	server := &Server{
		cfg: &config.Config{
			PollTimeout:  100,
			WorkersLimit: 1,
		},
		in:      make(chan *models.AgentTask),
		out:     make(chan *models.TaskResult),
		service: svc,
	}

	stream := mock.NewMockBidiClientStream[pb.TaskResult, pb.Task]()

	go func() {
		defer func() {
			stream.RecvClosed = true
			close(stream.RecvCh)
		}()

		for i := range 3 {
			stream.RecvCh <- pb.Task{
				Id:    fmt.Sprintf("test:recv:%d", i),
				Final: false,
			}

			stream.SetRecvErr(io.EOF)
		}
	}()

	go func() {
		for task := range server.in {
			fmt.Printf("task: %v", task)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := server.getTasks(ctx, stream)
	if err != nil && !errors.Is(err, context.Canceled) {
		t.Errorf("error getting tasks: %v", err)
	}
}

func TestSendTaskResults(t *testing.T) {
	svc := service.NewService()

	server := &Server{
		cfg: &config.Config{
			PollTimeout:  100,
			WorkersLimit: 1,
		},
		in:      make(chan *models.AgentTask),
		out:     make(chan *models.TaskResult),
		service: svc,
	}

	stream := mock.NewMockBidiClientStream[pb.TaskResult, pb.Task]()

	go func() {
		defer func() {
			close(server.out)
		}()

		for i := range 3 {
			server.out <- &models.TaskResult{
				Id:    fmt.Sprintf("test:recv:%d", i),
				Final: false,
			}

			stream.SetRecvErr(io.EOF)
		}
	}()

	go func() {
		for task := range stream.SendCh {
			fmt.Printf("task: %v", task)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := server.sendTaskResults(ctx, stream)
	if err != nil && !errors.Is(err, context.Canceled) {
		t.Errorf("error getting tasks: %v", err)
	}
}
