package service

import (
	"DistributedCalc/internal/orchestrator/config"
	"DistributedCalc/pkg/models"
	"context"
	"errors"
	"math"
	"math/rand"
	"sync"
)

const (
	number      = "NUMBER"
	operator    = "OPERATOR"
	parenthesis = "PARENTHESIS"
)

type PromiseChan struct {
	ch chan *models.TaskResult
}

// WaitForResult receives *models.TaskResult from channel and returns float64 result
// It closes the channel right after the value is returned
func (p *PromiseChan) WaitForResult(_ context.Context) (float64, error) {
	res := <-p.ch
	close(p.ch)
	return res.Result, nil
}

type PlannerChan struct {
	channels map[int]chan *models.TaskResult
	mu       sync.Mutex
	q        Queue[models.Task]
	cfg      *config.Config
}

func NewPlannerChan(cfg *config.Config, queue Queue[models.Task]) *PlannerChan {
	return &PlannerChan{
		cfg:      cfg,
		channels: make(map[int]chan *models.TaskResult),
		q:        queue,
	}
}

func (t *PlannerChan) PlanTask(ctx context.Context, task *models.Task) (TaskPromise, error) {
	id := rand.Int()
	ch := make(chan *models.TaskResult)

	t.mu.Lock()
	t.channels[id] = ch
	t.mu.Unlock()

	task.Id = id

	switch task.Operation {
	case "+":
		task.OperationTime = t.cfg.AdditionTime.Milliseconds()
	case "-":
		task.OperationTime = t.cfg.SubtractionTime.Milliseconds()
	case "*":
		task.OperationTime = t.cfg.MultiplicationTime.Milliseconds()
	case "/":
		if math.Abs(task.Arg2) < 1e-9 {
			return nil, models.ErrDivisionByZero
		}
		task.OperationTime = t.cfg.DivisionTime.Milliseconds()
	default:
		return nil, errors.Join(models.ErrUnsupportedOperation, errors.New(task.Operation))
	}

	err := t.q.Enqueue(ctx, task)
	if err != nil {
		return nil, err
	}

	return &PromiseChan{
		ch: ch,
	}, nil
}

func (t *PlannerChan) FinishTask(_ context.Context, res *models.TaskResult) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	ch, ok := t.channels[res.Id]
	if ok {
		ch <- res
		delete(t.channels, res.Id)
		return nil
	}

	return models.ErrTaskDoesNotExist
}
