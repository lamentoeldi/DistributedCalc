package service

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"github.com/distributed-calc/v1/internal/orchestrator/adapters/queue"
	memory2 "github.com/distributed-calc/v1/internal/orchestrator/blacklist/memory"
	"github.com/distributed-calc/v1/internal/orchestrator/config"
	"github.com/distributed-calc/v1/internal/orchestrator/repository/memory"
	"github.com/distributed-calc/v1/pkg/authenticator"
	"github.com/distributed-calc/v1/pkg/models"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestService_tokenize(t *testing.T) {
	s := NewService(nil, nil, nil, nil, nil)

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
	s := NewService(nil, nil, nil, nil, nil)

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
	s := NewService(r, p, q, nil, nil)

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
	s := NewService(r, p, q, nil, nil)

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

	err := s.repo.Add(context.Background(), &models.Expression{
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
	s := NewService(r, p, q, nil, nil)

	exp := &models.Expression{
		Id:     uuid.NewString(),
		Status: "testing",
		Result: 0,
	}

	err := s.repo.Add(context.Background(), exp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = s.GetAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTaskAST(t *testing.T) {
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

	s := NewService(nil, nil, nil, nil, nil)

	for _, tc := range cases {
		tokens, err := s.tokenize(tc.expression)
		if err != nil {
			t.Error(err)
		}

		ast, err := s.buildAST(tokens)
		if err != nil {
			t.Error(err)
		}

		tasks, err := s.taskAST(ast)
		if err != nil {
			t.Error(err)
		}

		_ = ast
		if err != nil {
			t.Error("failed to build AST")
		}

		for _, ta := range tasks {
			toPrint := fmt.Sprintf("%s %s %s %s %v", ta.id, ta.leftID, ta.rightID, ta.operator, ta.topLevel)
			if ta.arg1 != nil {
				toPrint += fmt.Sprintf(" %.2f", *ta.arg1)
			}

			if ta.arg2 != nil {
				toPrint += fmt.Sprintf(" %.2f", *ta.arg2)
			}

			if ta.result != nil {
				toPrint += fmt.Sprintf(" %.2f", *ta.result)
			}

			fmt.Println(toPrint)
		}
	}
}

func TestService_VerifyJWT(t *testing.T) {
	accessPk, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	refreshPk, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	s := NewService(
		memory.NewRepositoryMemory(),
		nil, nil,
		authenticator.NewAuthenticator(accessPk, refreshPk, 10*time.Second, 60*time.Second),
		memory2.NewBlacklist(),
	)

	cases := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "success",
			wantErr: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, refresh, err := s.auth.SignTokens(s.auth.IssueTokens("dqd"))
			if err != nil {
				t.Fatal(err)
			}

			err = s.VerifyJWT(context.Background(), refresh)
			if tc.wantErr == false && err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if tc.wantErr == true && err == nil {
				t.Error("expected error, got none")
			}
		})
	}
}

func TestService_RefreshTokens(t *testing.T) {
	accessPk, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	refreshPk, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	s := NewService(
		memory.NewRepositoryMemory(),
		nil, nil,
		authenticator.NewAuthenticator(accessPk, refreshPk, 10*time.Second, 60*time.Second),
		memory2.NewBlacklist(),
	)

	cases := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "success",
			wantErr: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, refresh, err := s.auth.SignTokens(s.auth.IssueTokens("dqd"))
			if err != nil {
				t.Fatal(err)
			}

			_, refresh, err = s.RefreshTokens(context.Background(), refresh)
			if tc.wantErr == false && err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if tc.wantErr == true && err == nil {
				t.Error("expected error, got none")
			}

			err = s.VerifyJWT(context.Background(), refresh)
			if tc.wantErr == true && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}
