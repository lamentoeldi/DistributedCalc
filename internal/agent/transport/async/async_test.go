package async

import (
	"DistributedCalc/internal/agent/config"
	"DistributedCalc/pkg/models"
	"DistributedCalc/test/mock"
	"context"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestNewTransportAsync(t *testing.T) {
	cfg := &config.Config{
		PollTimeout:  100,
		WorkersLimit: 5,
		Url:          "",
	}
	log, _ := zap.NewDevelopment()

	o := &mock.OrchestratorMock{Err: nil}
	c := &mock.CalculatorMock{Err: nil}

	in := make(chan *models.Task, 64)
	out := make(chan *models.TaskResult, 64)

	transport := NewTransportAsync(cfg, log, o, c, in, out)
	if transport == nil {
		t.Error("Failed to create transport")
	}
}

func TestTransportAsync_Run(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := &config.Config{
		PollTimeout:  1000,
		WorkersLimit: 5,
		Url:          "",
	}
	log, _ := zap.NewDevelopment()

	o := &mock.OrchestratorMock{Err: nil}
	c := &mock.CalculatorMock{Err: nil}

	in := make(chan *models.Task, 64)
	out := make(chan *models.TaskResult, 64)

	transport := NewTransportAsync(cfg, log, o, c, in, out)
	if transport == nil {
		t.Fatal("Failed to create transport")
	}

	go func() {
		time.Sleep(2 * time.Second)
		cancel()
	}()

	transport.Run(ctx)
	<-ctx.Done()
	transport.Shutdown()
}

func TestTransportAsync_Shutdown(t *testing.T) {
	cfg := &config.Config{
		PollTimeout:  100,
		WorkersLimit: 5,
		Url:          "",
	}
	log, _ := zap.NewDevelopment()

	o := &mock.OrchestratorMock{Err: nil}
	c := &mock.CalculatorMock{Err: nil}

	in := make(chan *models.Task, 64)
	out := make(chan *models.TaskResult, 64)

	transport := NewTransportAsync(cfg, log, o, c, in, out)
	if transport == nil {
		t.Fatal("Failed to create transport")
	}

	transport.Shutdown()
}
