package mock

import (
	"DistributedCalc/pkg/models"
	"context"
	"fmt"
	"math/rand"
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

type ServiceMock struct {
	Err error
}

func (s *ServiceMock) StartEvaluation(_ context.Context, _ string) (int, error) {
	return rand.Int(), s.Err
}

func (s *ServiceMock) Get(_ context.Context, id int) (*models.Expression, error) {
	return &models.Expression{
		Id:     id,
		Status: "completed",
		Result: 0,
	}, s.Err
}

func (s *ServiceMock) GetAll(_ context.Context) ([]*models.Expression, error) {
	return []*models.Expression{
		{
			Id:     rand.Int(),
			Status: "completed",
			Result: 0,
		},
	}, s.Err
}

func (s *ServiceMock) FinishTask(_ context.Context, _ *models.TaskResult) error {
	return s.Err
}

func (s *ServiceMock) Enqueue(_ context.Context, _ *models.Task) error {
	return s.Err
}

func (s *ServiceMock) Dequeue(_ context.Context) (*models.Task, error) {
	return &models.Task{
		Id:            rand.Int(),
		Arg1:          2,
		Arg2:          3,
		Operation:     "+",
		OperationTime: 5,
	}, s.Err
}
