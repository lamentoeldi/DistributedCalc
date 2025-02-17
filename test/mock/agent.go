package mock

import (
	"DistributedCalc/pkg/models"
	"context"
	"fmt"
)

type OrchestratorMock struct {
	Err error
}

func (o *OrchestratorMock) GetTask(_ context.Context) (*models.Task, error) {
	if o.Err != nil {
		return nil, o.Err
	}

	return &models.Task{
		Id:        10,
		Operation: "+",
		Arg1:      2.0,
		Arg2:      3.0,
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

func (c *CalculatorMock) Evaluate(task *models.Task) (*models.TaskResult, error) {
	fmt.Println("evaluate called: ", task)

	if c.Err != nil {
		return nil, c.Err
	}

	return &models.TaskResult{
		Id:     10,
		Result: 5.0,
	}, nil
}
