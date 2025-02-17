package service

import (
	"fmt"
	"testing"
)

func TestTokenize(t *testing.T) {
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
			tokens, err := tokenize(tc.exp)
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

func TestBuildAST(t *testing.T) {
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
			tokens, err := tokenize(tc.expression)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			ast, err := buildAST(tokens)
			fmt.Println(*ast)
			if tc.wantErr == false && err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if tc.wantErr == true && err == nil {
				t.Error("expected error, got none")
			}

			res, err := evaluateAST(ast)
			if int(tc.expected*1000) != int(res*1000) {
				t.Errorf("expected %v, got %v", tc.expected, res)
			}

			if tc.wantErr == false && err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if tc.wantErr == true && err == nil {
				t.Error("expected error, got none")
			}
		})
	}
}
