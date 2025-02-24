package service

import (
	"DistributedCalc/pkg/models"
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
		task     *models.Task
		expected *models.TaskResult
		wantErr  bool
	}{
		{
			name: "addition",
			task: &models.Task{
				Id:            1,
				Operation:     "+",
				Arg1:          2,
				Arg2:          3,
				OperationTime: 100,
			},
			expected: &models.TaskResult{
				Id:     1,
				Result: 5,
			},
			wantErr: false,
		},
		{
			name: "subtraction",
			task: &models.Task{
				Id:            2,
				Operation:     "-",
				Arg1:          5,
				Arg2:          3,
				OperationTime: 200,
			},
			expected: &models.TaskResult{
				Id:     2,
				Result: 2,
			},
			wantErr: false,
		},
		{
			name: "multiplication",
			task: &models.Task{
				Id:            3,
				Operation:     "*",
				Arg1:          2,
				Arg2:          3,
				OperationTime: 300,
			},
			expected: &models.TaskResult{
				Id:     3,
				Result: 6,
			},
			wantErr: false,
		},
		{
			name: "division",
			task: &models.Task{
				Id:            4,
				Operation:     "/",
				Arg1:          10,
				Arg2:          2,
				OperationTime: 400,
			},
			expected: &models.TaskResult{
				Id:     4,
				Result: 5,
			},
			wantErr: false,
		},
		{
			name: "division by zero",
			task: &models.Task{
				Id:            5,
				Operation:     "/",
				Arg1:          5,
				Arg2:          0,
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
