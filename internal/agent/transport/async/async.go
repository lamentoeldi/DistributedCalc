package async

import (
	"DistributedCalc/internal/agent/config"
	"DistributedCalc/pkg/models"
	"context"
	"errors"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Orchestrator interface {
	GetTask(ctx context.Context) (*models.Task, error)
	PostResult(ctx context.Context, result *models.TaskResult) error
}

type Calculator interface {
	Evaluate(task *models.Task) (*models.TaskResult, error)
}

type TransportAsync struct {
	o        Orchestrator
	c        Calculator
	in       chan *models.Task
	out      chan *models.TaskResult
	cfg      *config.Config
	log      *zap.Logger
	onceRun  sync.Once
	onceStop sync.Once
	wg       sync.WaitGroup
}

func NewTransportAsync(cfg *config.Config, log *zap.Logger, o Orchestrator, c Calculator, in chan *models.Task, out chan *models.TaskResult) *TransportAsync {
	return &TransportAsync{
		o:   o,
		in:  in,
		out: out,
		c:   c,
		log: log,
		cfg: cfg,
	}
}

// StartWorkers starts worker goroutines which take *models.Task from in, evaluate term and pass *models.TaskResult to out
func (t *TransportAsync) StartWorkers(ctx context.Context) {
	t.wg.Add(t.cfg.WorkersLimit)
	for range t.cfg.WorkersLimit {
		go func() {
			defer t.wg.Done()

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
}

// produce takes *models.TaskResult from out channel and sends it back to the orchestrator
func (t *TransportAsync) produce() {
	for r := range t.out {
		err := t.o.PostResult(context.TODO(), r)
		if err != nil {
			t.log.Error(err.Error())
			continue
		}
	}
}

// consume uses long polling to receive new tasks from server and send them to in channel
func (t *TransportAsync) consume(ctx context.Context) {
	t.wg.Add(1)

	ticker := time.NewTicker(t.cfg.PollTimeout)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			t.wg.Done()
			return
		case <-ticker.C:
			task, err := t.o.GetTask(ctx)
			if err != nil {
				switch {
				case errors.Is(err, models.ErrNoTasks):
				default:
					t.log.Error(err.Error())
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
	})
}

// Shutdown closes channels which causes all workers to stop
func (t *TransportAsync) Shutdown() {
	t.onceStop.Do(func() {
		t.wg.Wait()
		close(t.out)
		close(t.in)
	})
}
