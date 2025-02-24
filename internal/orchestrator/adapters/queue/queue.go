package queue

import (
	"DistributedCalc/pkg/models"
	"context"
)

type QueueChan[T any] struct {
	queue chan *T
}

func NewQueueChan[T any](queueSize int) *QueueChan[T] {
	return &QueueChan[T]{
		queue: make(chan *T, queueSize),
	}
}

func (q *QueueChan[T]) Enqueue(_ context.Context, obj *T) error {
	q.queue <- obj
	return nil
}

func (q *QueueChan[T]) Dequeue(_ context.Context) (*T, error) {
	select {
	case obj := <-q.queue:
		return obj, nil
	default:
		return nil, models.ErrNoTasks
	}
}
