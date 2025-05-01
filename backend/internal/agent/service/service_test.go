package service

import (
	"fmt"
	"github.com/distributed-calc/v1/internal/agent/models"
	"testing"
)

func TestNewService(t *testing.T) {
	service := NewService()
	if service == nil {
		t.Fatal("NewService() returned nil")
	}
}

func TestService_Evaluate(t *testing.T) {
	cases := []struct {
		name     string
		task     *models.AgentTask
		expected *models.TaskResult
		wantErr  bool
	}{
		{
			name: "addition",
			task: &models.AgentTask{
				Id:            fmt.Sprint(1),
				Op:            "+",
				LeftArg:       2,
				RightArg:      3,
				OperationTime: 100,
			},
			expected: &models.TaskResult{
				Id:     fmt.Sprint(1),
				Result: 5,
			},
			wantErr: false,
		},
		{
			name: "subtraction",
			task: &models.AgentTask{
				Id:            fmt.Sprint(2),
				Op:            "-",
				LeftArg:       5,
				RightArg:      3,
				OperationTime: 200,
			},
			expected: &models.TaskResult{
				Id:     fmt.Sprint(2),
				Result: 2,
			},
			wantErr: false,
		},
		{
			name: "multiplication",
			task: &models.AgentTask{
				Id:            fmt.Sprint(3),
				Op:            "*",
				LeftArg:       2,
				RightArg:      3,
				OperationTime: 300,
			},
			expected: &models.TaskResult{
				Id:     fmt.Sprint(3),
				Result: 6,
			},
			wantErr: false,
		},
		{
			name: "division",
			task: &models.AgentTask{
				Id:            fmt.Sprint(4),
				Op:            "/",
				LeftArg:       10,
				RightArg:      2,
				OperationTime: 400,
			},
			expected: &models.TaskResult{
				Id:     fmt.Sprint(4),
				Result: 5,
			},
			wantErr: false,
		},
		{
			name: "division by zero",
			task: &models.AgentTask{
				Id:            fmt.Sprint(5),
				Op:            "/",
				LeftArg:       5,
				RightArg:      0,
				OperationTime: 500,
			},
			wantErr: true,
		},
	}

	service := NewService()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r, err := service.Evaluate(tc.task)
			if tc.wantErr && err == nil {
				t.Error("expected error, got none")
			}

			if tc.wantErr == false && err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if err != nil {
				return
			}

			if r.Id != tc.expected.Id || r.Result != tc.expected.Result {
				t.Errorf("expected %v, got %v", tc.expected, r)
			}
		})
	}
}
