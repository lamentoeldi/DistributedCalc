package service

import (
	"DistributedCalc/pkg/models"
	"context"
)

type Orchestrator interface {
	GetTask(ctx context.Context) (*models.Task, error)
	PostResult(ctx context.Context, result *models.TaskResult) error
}

type Queue[T any] struct {
	queue chan *T
}

func NewQueue[T any](queueSize int) *Queue[T] {
	return &Queue[T]{
		queue: make(chan *T, queueSize),
	}
}

func (q *Queue[T]) Enqueue(obj *T) {
	q.queue <- obj
}

func (q *Queue[T]) Dequeue() (*T, error) {
	select {
	case obj := <-q.queue:
		return obj, nil
	default:
		return nil, models.ErrNoTasks
	}
}

type Service struct {
	o Orchestrator
}

func NewService(o Orchestrator) *Service {
	return &Service{
		o: o,
	}
}

func (s *Service) Evaluate(_ context.Context, expression string) (float64, error) {
	panic("unimplemented")
}
