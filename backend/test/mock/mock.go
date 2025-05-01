package mock

import (
	"context"
	"fmt"
	"github.com/distributed-calc/v1/internal/agent/models"
)

type OrchestratorMock struct {
	Err error
}

func (o *OrchestratorMock) GetTask(_ context.Context) (*models.AgentTask, error) {
	if o.Err != nil {
		return nil, o.Err
	}

	return &models.AgentTask{
		Id:       fmt.Sprint(10),
		Op:       "+",
		LeftArg:  2.0,
		RightArg: 3.0,
	}, nil
}

func (o *OrchestratorMock) PostResult(_ context.Context, _ *models.TaskResult) error {
	if o.Err != nil {
		return o.Err
	}

	return nil
}

type CalculatorMock struct {
	Err error
}

func (c *CalculatorMock) Evaluate(task *models.AgentTask) (*models.TaskResult, error) {
	fmt.Println("evaluate called: ", task)

	if c.Err != nil {
		return nil, c.Err
	}

	return &models.TaskResult{
		Id:     fmt.Sprint(10),
		Result: 5.0,
	}, nil
}

type ServiceMock struct {
	Err error
}
