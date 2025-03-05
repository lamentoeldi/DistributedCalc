package service

import (
	"fmt"
	"github.com/distributed-calc/v1/pkg/models"
	"time"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Evaluate(t *models.Task) (*models.TaskResult, error) {
	time.Sleep(time.Duration(t.OperationTime) * time.Millisecond)

	switch t.Operation {
	case "+":
		return &models.TaskResult{
			Id:     t.Id,
			Result: t.Arg1 + t.Arg2,
		}, nil
	case "-":
		return &models.TaskResult{
			Id:     t.Id,
			Result: t.Arg1 - t.Arg2,
		}, nil
	case "*":
		return &models.TaskResult{
			Id:     t.Id,
			Result: t.Arg1 * t.Arg2,
		}, nil
	case "/":
		if t.Arg2 == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		return &models.TaskResult{
			Id:     t.Id,
			Result: t.Arg1 / t.Arg2,
		}, nil
	}

	return nil, fmt.Errorf("unknown operation %s", t.Operation)
}
