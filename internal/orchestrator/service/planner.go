package service

import (
	"DistributedCalc/internal/orchestrator/config"
	"DistributedCalc/pkg/models"
	"context"
	"fmt"
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
	q        *Queue[models.Task]
	cfg      *config.Config
}

func NewPlannerChan(cfg *config.Config, queue *Queue[models.Task]) *PlannerChan {
	return &PlannerChan{
		cfg:      cfg,
		channels: make(map[int]chan *models.TaskResult),
		q:        queue,
	}
}

func (t *PlannerChan) PlanTask(_ context.Context, task *models.Task) (TaskPromise, error) {
	id := rand.Int()
	ch := make(chan *models.TaskResult)

	t.mu.Lock()
	t.channels[id] = ch
	t.mu.Unlock()

	task.Id = id

	t.q.Enqueue(task)

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

	return fmt.Errorf("some error")
}
