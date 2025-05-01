package service

import (
	"fmt"
	"github.com/distributed-calc/v1/internal/agent/models"
	"time"
)

const (
	statusSuccess = "completed"
	statusFailure = "failure"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Evaluate(t *models.AgentTask) (*models.TaskResult, error) {
	time.Sleep(time.Duration(t.OperationTime) * time.Millisecond)

	switch t.Op {
	case "+":
		return &models.TaskResult{
			Id:     t.Id,
			Result: t.LeftArg + t.RightArg,
			Status: statusSuccess,
			Final:  t.Final,
		}, nil
	case "-":
		return &models.TaskResult{
			Id:     t.Id,
			Result: t.LeftArg - t.RightArg,
			Status: statusSuccess,
			Final:  t.Final,
		}, nil
	case "*":
		return &models.TaskResult{
			Id:     t.Id,
			Result: t.LeftArg * t.RightArg,
			Status: statusSuccess,
			Final:  t.Final,
		}, nil
	case "/":
		if t.RightArg == 0 {
			return &models.TaskResult{
				Id:     t.Id,
				Status: statusFailure,
				Final:  t.Final,
			}, fmt.Errorf("division by zero")
		}

		return &models.TaskResult{
			Id:     t.Id,
			Result: t.LeftArg / t.RightArg,
			Status: statusSuccess,
			Final:  t.Final,
		}, nil
	case "":
		return &models.TaskResult{
			Id:     t.Id,
			Result: t.LeftArg,
			Status: statusSuccess,
			Final:  t.Final,
		}, nil
	}

	return nil, fmt.Errorf("unknown operation %s", t.Op)
}
