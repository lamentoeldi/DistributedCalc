package service

import (
	"context"
	"github.com/distributed-calc/v1/internal/orchestrator/models"
	"github.com/distributed-calc/v1/internal/orchestrator/repository/memory"
	"github.com/google/uuid"
	"testing"
)

func TestValidate(t *testing.T) {
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
			err := validate(tc.exp)
			if tc.wantErr == false && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if tc.wantErr == true && err == nil {
				t.Error("expected error, got none")
			}
		})
	}
}

func TestParseExpression(t *testing.T) {
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
			_, err := parseExpression(tc.expression, tc.name)
			if tc.wantErr == false && err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if tc.wantErr == true && err == nil {
				t.Error("expected error, got none")
			}
		})
	}
}

func TestService_Evaluate(t *testing.T) {
	repo := memory.NewRepositoryMemory()
	s := NewService(nil, repo, repo, nil, nil, nil)

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
			_, err := s.Evaluate(context.Background(), tc.expression, "")
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
	repo := memory.NewRepositoryMemory()
	s := NewService(nil, repo, repo, nil, nil, nil)

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

	err := s.expRepo.Add(context.Background(), &models.Expression{
		Id:     found,
		Status: "testing",
		Result: 0,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.Get(context.Background(), tc.id, "")
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
	repo := memory.NewRepositoryMemory()
	s := NewService(nil, repo, repo, nil, nil, nil)

	exp := &models.Expression{
		Id:     uuid.NewString(),
		Status: "testing",
		Result: 0,
	}

	err := s.expRepo.Add(context.Background(), exp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = s.GetAll(context.Background(), "", "", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
