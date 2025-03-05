package service

import (
	"context"
	"fmt"
	"github.com/distributed-calc/v1/internal/orchestrator/adapters/queue"
	"github.com/distributed-calc/v1/internal/orchestrator/config"
	"github.com/distributed-calc/v1/internal/orchestrator/repository/memory"
	"github.com/distributed-calc/v1/pkg/models"
	"github.com/google/uuid"
	"log"
	"testing"
	"time"
)

func TestService_tokenize(t *testing.T) {
	s := NewService(nil, nil, nil)

	cases := []struct {
		name    string
		exp     string
		wantErr bool
	}{
		{
			name:    "simple expression",
			exp:     "2+2",
			wantErr: false,
		},
		{
			name:    "expression with parenthesis",
			exp:     "(2+2)*3",
			wantErr: false,
		},
		{
			name:    "expression with operator priority",
			exp:     "2*2+3",
			wantErr: false,
		},
		{
			name:    "expression with invalid character",
			exp:     "2+2a",
			wantErr: true,
		},
		{
			name:    "expression with float numbers",
			exp:     "3.14+2",
			wantErr: false,
		},
		{
			name:    "expression with nested parenthesis",
			exp:     "((2+3)*2)+1",
			wantErr: false,
		},
		{
			name:    "expressions with explicit signs",
			exp:     "(+2-3)",
			wantErr: false,
		},
		{
			name:    "expression with invalid parenthesis",
			exp:     "2+(3",
			wantErr: true,
		},
		{
			name:    "expression with invalid nested parenthesis",
			exp:     "(2+(3)",
			wantErr: true,
		},
		{
			name:    "expression with invalid float numbers",
			exp:     "3.14.2+2",
			wantErr: true,
		},
		{
			name:    "expression with operators in the beginning",
			exp:     "*2+3",
			wantErr: true,
		},
		{
			name:    "expression with operators in the end",
			exp:     "2+3*",
			wantErr: true,
		},
		{
			name:    "expression with empty parentheses",
			exp:     "2+()-1",
			wantErr: true,
		},
		{
			name:    "expression with operators mismatch",
			exp:     "2++3",
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tokens, err := s.tokenize(tc.exp)
			fmt.Println("tokens: ", tokens)
			if tc.wantErr == false && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if tc.wantErr == true && err == nil {
				t.Error("expected error, got none")
			}
		})
	}
}

func TestService_buildAST(t *testing.T) {
	s := NewService(nil, nil, nil)

	cases := []struct {
		name       string
		expression string
		expected   float64
		wantErr    bool
	}{
		{
			name:       "simple expression",
			expression: "2+2",
			expected:   4,
			wantErr:    false,
		},
		{
			name:       "expression with parenthesis",
			expression: "(2+2)*3",
			expected:   12,
			wantErr:    false,
		},
		{
			name:       "expression with operator priority",
			expression: "2*2+3",
			expected:   7,
			wantErr:    false,
		},
		{
			name:       "expression with float numbers",
			expression: "3.14+2",
			expected:   5.14,
			wantErr:    false,
		},
		{
			name:       "expression with nested parenthesis",
			expression: "((2+3)*2)+1",
			expected:   11,
			wantErr:    false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tokens, err := s.tokenize(tc.expression)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			ast, err := s.buildAST(tokens)
			fmt.Println(*ast)
			if tc.wantErr == false && err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if tc.wantErr == true && err == nil {
				t.Error("expected error, got none")
			}
		})
	}
}

func TestService_evaluateAST(t *testing.T) {
	cfg := &config.Config{
		Host: "",
		Port: 0,

		AdditionTime:       100 * time.Millisecond,
		SubtractionTime:    500 * time.Millisecond,
		MultiplicationTime: 2000 * time.Millisecond,
		DivisionTime:       1000 * time.Millisecond,
	}

	r := memory.NewRepositoryMemory()
	q := queue.NewQueueChan[models.Task](64)
	p := NewPlannerChan(cfg, q)
	s := NewService(r, p, q)

	cases := []struct {
		name       string
		expression string
		expected   float64
		wantErr    bool
	}{
		{
			name:       "simple expression",
			expression: "2+2",
			expected:   4,
			wantErr:    false,
		},
		{
			name:       "expression with parenthesis",
			expression: "(2+2)*3",
			expected:   12,
			wantErr:    false,
		},
		{
			name:       "expression with operator priority",
			expression: "2*2+3",
			expected:   7,
			wantErr:    false,
		},
		{
			name:       "expression with float numbers",
			expression: "3.14+2",
			expected:   5.14,
			wantErr:    false,
		},
		{
			name:       "expression with nested parenthesis",
			expression: "((2+3)*2)+1",
			expected:   11,
			wantErr:    false,
		},
		{
			name:       "expression with paralleled tasks",
			expression: "(2+3)*(4+1)",
			expected:   25,
			wantErr:    false,
		},
		{
			name:       "expression with more paralleled tasks",
			expression: "(2+3)*(4+1)+(5+1)*(5+5)",
			expected:   85,
			wantErr:    false,
		},
		{
			name:       "negative unary expression",
			expression: "-2",
			expected:   -2,
			wantErr:    false,
		},
		{
			name:       "positive unary expression",
			expression: "2",
			expected:   2,
			wantErr:    false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// This is agent implementation for tests
			go func() {
				for {
					task, err := q.Dequeue(context.Background())
					if err != nil {
						continue
					}

					time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)

					res := &models.TaskResult{
						Id: task.Id,
					}
					switch task.Operation {
					case "+":
						res.Result = task.Arg1 + task.Arg2
					case "-":
						res.Result = task.Arg1 - task.Arg2
					case "*":
						res.Result = task.Arg1 * task.Arg2
					case "/":
						res.Result = task.Arg1 / task.Arg2
					}

					err = s.p.FinishTask(context.TODO(), res)
					if err != nil {
						log.Printf("Error finishing task %d: %v", task.Id, err)
						continue
					}
				}
			}()

			tokens, err := s.tokenize(tc.expression)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			ast, err := s.buildAST(tokens)
			if tc.wantErr == false && err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			res, err := s.evaluateAST(ast)
			if tc.wantErr == false && err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if tc.wantErr == true && err == nil {
				t.Error("expected error, got none")
			}

			if int(tc.expected*1000) != int(res*1000) {
				t.Errorf("expected %f, got %f", tc.expected, res)
			}

			log.Printf("calculated value: %.2f for %s", res, tc.name)
		})
	}
}

func TestService_StartEvaluation(t *testing.T) {
	cfg := &config.Config{
		Host: "",
		Port: 0,

		AdditionTime:       100 * time.Millisecond,
		SubtractionTime:    500 * time.Millisecond,
		MultiplicationTime: 2000 * time.Millisecond,
		DivisionTime:       1000 * time.Millisecond,
	}

	r := memory.NewRepositoryMemory()
	q := queue.NewQueueChan[models.Task](64)
	p := NewPlannerChan(cfg, q)
	s := NewService(r, p, q)

	cases := []struct {
		name       string
		expression string
		expected   float64
		wantErr    bool
	}{
		{
			name:       "valid expression",
			expression: "2+2",
			wantErr:    false,
		},
		{
			name:       "invalid expression",
			expression: "2+2/",
			wantErr:    true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.StartEvaluation(context.Background(), tc.expression)
			if tc.wantErr == false && err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if tc.wantErr == true && err == nil {
				t.Error("expected error, got none")
			}
		})
	}
}

func TestService_Get(t *testing.T) {
	cfg := &config.Config{
		Host: "",
		Port: 0,

		AdditionTime:       100 * time.Millisecond,
		SubtractionTime:    500 * time.Millisecond,
		MultiplicationTime: 2000 * time.Millisecond,
		DivisionTime:       1000 * time.Millisecond,
	}

	r := memory.NewRepositoryMemory()
	q := queue.NewQueueChan[models.Task](64)
	p := NewPlannerChan(cfg, q)
	s := NewService(r, p, q)

	found := uuid.NewString()

	cases := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "success",
			id:      found,
			wantErr: false,
		},
		{
			name:    "expression not found",
			id:      uuid.NewString(),
			wantErr: true,
		},
	}

	err := s.r.Add(context.Background(), &models.Expression{
		Id:     found,
		Status: "testing",
		Result: 0,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.Get(context.Background(), tc.id)
			if tc.wantErr == false && err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if tc.wantErr == true && err == nil {
				t.Error("expected error, got none")
			}
		})
	}
}

func TestService_GetAll(t *testing.T) {
	cfg := &config.Config{
		Host: "",
		Port: 0,

		AdditionTime:       100 * time.Millisecond,
		SubtractionTime:    500 * time.Millisecond,
		MultiplicationTime: 2000 * time.Millisecond,
		DivisionTime:       1000 * time.Millisecond,
	}

	r := memory.NewRepositoryMemory()
	q := queue.NewQueueChan[models.Task](64)
	p := NewPlannerChan(cfg, q)
	s := NewService(r, p, q)

	exp := &models.Expression{
		Id:     uuid.NewString(),
		Status: "testing",
		Result: 0,
	}

	err := s.r.Add(context.Background(), exp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = s.GetAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
