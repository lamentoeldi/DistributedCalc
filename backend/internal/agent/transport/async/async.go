package async

import (
	"context"
	"errors"
	"github.com/distributed-calc/v1/internal/agent/config"
	e "github.com/distributed-calc/v1/internal/agent/errors"
	"github.com/distributed-calc/v1/internal/agent/models"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Orchestrator interface {
	GetTask(ctx context.Context) (*models.AgentTask, error)
	PostResult(ctx context.Context, result *models.TaskResult) error
}

type Calculator interface {
	Evaluate(task *models.AgentTask) (*models.TaskResult, error)
}

type TransportAsync struct {
	o        Orchestrator
	c        Calculator
	in       chan *models.AgentTask
	out      chan *models.TaskResult
	cfg      *config.Config
	log      *zap.Logger
	onceRun  sync.Once
	onceStop sync.Once
}

func NewTransportAsync(cfg *config.Config, log *zap.Logger, o Orchestrator, c Calculator) *TransportAsync {
	return &TransportAsync{
		o:   o,
		in:  make(chan *models.AgentTask, cfg.BufferSize),
		out: make(chan *models.TaskResult, cfg.BufferSize),
		c:   c,
		log: log,
		cfg: cfg,
	}
}

// StartWorkers starts worker goroutines which take *models.Task from in, evaluate term and pass *models.TaskResult to out
func (t *TransportAsync) StartWorkers(ctx context.Context) {
	wg := &sync.WaitGroup{}

	wg.Add(t.cfg.WorkersLimit)
	for range t.cfg.WorkersLimit {
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case r, ok := <-t.in:
					if !ok {
						continue
					}

					res, err := t.c.Evaluate(r)
					if err != nil {
						t.log.Error(err.Error())
						continue
					}
					t.out <- res
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(t.out)
	}()
}

// produce takes *models.TaskResult from out channel and sends it back to the orchestrator
func (t *TransportAsync) produce() {
	for r := range t.out {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err := t.o.PostResult(ctx, r)
			if err != nil {
				t.log.Error("failed to post result", zap.Error(err))
			}
		}()
	}
}

// consume uses long polling to receive new tasks from server and send them to in channel
func (t *TransportAsync) consume(ctx context.Context) {
	ticker := time.NewTicker(t.cfg.PollTimeout)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			close(t.in)
			return
		case <-ticker.C:
			task, err := t.o.GetTask(ctx)
			if err != nil {
				switch {
				case errors.Is(err, e.ErrNoTasks):
				default:
					t.log.Error("failed to get task", zap.Error(err))
				}
				break
			}
			t.in <- task
		}
	}
}

// Run starts consuming tasks from server
func (t *TransportAsync) Run(ctx context.Context) {
	t.onceRun.Do(func() {
		t.StartWorkers(ctx)
		go t.consume(ctx)
		go t.produce()

		t.log.Info("started polling orchestrator")
	})
}

// Shutdown closes channels which causes all workers to stop
func (t *TransportAsync) Shutdown() {
	t.onceStop.Do(func() {
		t.log.Info("shutting down...")
		t.log.Sync()
	})
}
