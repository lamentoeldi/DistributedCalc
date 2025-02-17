package service

import (
	"DistributedCalc/pkg/models"
	"context"
)

type Orchestrator interface {
	GetTask(ctx context.Context) (*models.Task, error)
	PostResult(ctx context.Context, result *models.TaskResult) error
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
