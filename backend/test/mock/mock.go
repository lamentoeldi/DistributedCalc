package mock

import (
	"context"
	"fmt"
	ma "github.com/distributed-calc/v1/internal/agent/models"
	mo "github.com/distributed-calc/v1/internal/orchestrator/models"
	"github.com/google/uuid"
)

type OrchestratorMock struct {
	Err error
}

func (o *OrchestratorMock) GetTask(_ context.Context) (*ma.AgentTask, error) {
	if o.Err != nil {
		return nil, o.Err
	}

	return &ma.AgentTask{
		Id:       fmt.Sprint(10),
		Op:       "+",
		LeftArg:  2.0,
		RightArg: 3.0,
	}, nil
}

func (o *OrchestratorMock) PostResult(_ context.Context, _ *ma.TaskResult) error {
	if o.Err != nil {
		return o.Err
	}

	return nil
}

type CalculatorMock struct {
	Err error
}

func (c *CalculatorMock) Evaluate(task *ma.AgentTask) (*ma.TaskResult, error) {
	fmt.Println("evaluate called: ", task)

	if c.Err != nil {
		return nil, c.Err
	}

	return &ma.TaskResult{
		Id:     fmt.Sprint(10),
		Result: 5.0,
	}, nil
}

type ServiceMock struct {
	Err error
}

func (s ServiceMock) Evaluate(_ context.Context, expression string) (string, error) {
	if s.Err != nil {
		return "", s.Err
	}

	id, _ := uuid.NewV7()

	return id.String(), nil
}

func (s ServiceMock) Get(_ context.Context, _ string) (*mo.Expression, error) {
	if s.Err != nil {
		return nil, s.Err
	}

	return &mo.Expression{
		Id:     fmt.Sprint(10),
		Result: 0.0,
		Status: "pending",
	}, nil
}

func (s ServiceMock) GetAll(_ context.Context) ([]*mo.Expression, error) {
	if s.Err != nil {
		return nil, s.Err
	}

	return []*mo.Expression{
		{
			Id:     fmt.Sprint(10),
			Result: 0.0,
			Status: "pending",
		},
	}, nil
}

func (s ServiceMock) GetTask(_ context.Context) (*mo.AgentTask, error) {
	if s.Err != nil {
		return nil, s.Err
	}

	return &mo.AgentTask{
		Id:            fmt.Sprint(10),
		LeftArg:       10,
		RightArg:      10,
		Op:            "+",
		OperationTime: 0,
		Final:         true,
	}, nil
}

func (s ServiceMock) FinishTask(_ context.Context, _ *mo.TaskResult) error {
	if s.Err != nil {
		return s.Err
	}

	return nil
}

func (s ServiceMock) Finalize(_ context.Context, _ string, _ float64) error {
	if s.Err != nil {
		return s.Err
	}

	return nil
}
