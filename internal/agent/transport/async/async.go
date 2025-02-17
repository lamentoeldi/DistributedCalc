package async

import (
	"DistributedCalc/internal/agent/adapters/orchestrator"
	"DistributedCalc/internal/agent/config"
	"DistributedCalc/pkg/models"
	"context"
	"go.uber.org/zap"
	"time"
)

type Calculator interface {
	Evaluate(task *models.Task) (*models.TaskResult, error)
}

type TransportAsync struct {
	o   *orchestrator.Orchestrator
	c   Calculator
	in  chan *models.Task
	out chan *models.TaskResult
	cfg *config.Config
	log *zap.Logger
}

func NewTransportAsync(cfg *config.Config, log *zap.Logger, o *orchestrator.Orchestrator, c Calculator, in chan *models.Task, out chan *models.TaskResult) *TransportAsync {
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
func (t *TransportAsync) StartWorkers() {
	for range t.cfg.WorkersLimit {
		go func() {
			for r := range t.in {
				res, err := t.c.Evaluate(r)
				if err != nil {
					t.log.Error(err.Error())
					continue
				}
				t.out <- res
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
	ticker := time.NewTicker(t.cfg.PollTimeout)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			task, err := t.o.GetTask(nil)
			if err != nil {
				t.log.Error(err.Error())
				continue
			}
			t.in <- task
		}
	}
}

// Run starts consuming tasks from server
func (t *TransportAsync) Run(ctx context.Context) {
	t.StartWorkers()
	go t.consume(ctx)
	go t.produce()
}

// Shutdown closes channels which causes all workers to stop
func (t *TransportAsync) Shutdown() {
	close(t.out)
	close(t.in)
}
