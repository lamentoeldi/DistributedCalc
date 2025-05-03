package http

import (
	"bytes"
	"context"
	"fmt"
	"github.com/distributed-calc/v1/internal/orchestrator/errors"
	"github.com/distributed-calc/v1/test/mock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

var th *Server
var s *mock.ServiceMock

func init() {
	log, _ := zap.NewDevelopment()
	s = &mock.ServiceMock{Err: nil}
	cfg := &TransportHttpConfig{
		Host: "localhost",
		Port: 8080,
	}

	th = NewServer(s, log, cfg)
}

func TestTransportHttp_Run(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			t.Fatalf("unexpected panic: %v", err)
		}
	}()

	th.Run()
}

func TestTransportHttp_Shutdown(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			t.Fatalf("unexpected panic: %v", err)
		}
	}()

	th.Run()
	th.Shutdown(context.Background())
}

func TestTransportHttp_handlePing(t *testing.T) {
	req := httptest.NewRequest("GET", "/ping", nil)
	r := httptest.NewRecorder()

	th.handlePing(r, req)

	if r.Code != http.StatusOK {
		t.Errorf("expected status code 200, got %d", r.Code)
	}
}

func TestTransportHttp_handleCalculate(t *testing.T) {
	defer func() {
		s.Err = nil
	}()

	cases := []struct {
		name           string
		expression     string
		method         string
		err            error
		expectedStatus int
	}{
		{
			name:           "valid expression",
			expression:     `{"expression": "2+2"}`,
			method:         "POST",
			err:            nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "bad request",
			expression:     "invalid expression",
			method:         "POST",
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "unprocessable entity",
			expression:     `{"expression": "2+("}`,
			method:         "POST",
			err:            errors.ErrInvalidExpression,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "method not allowed",
			expression:     "2+2",
			method:         "GET",
			err:            nil,
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s.Err = tc.err

			req := httptest.NewRequest(tc.method, "/api/v1/calculate", bytes.NewReader([]byte(tc.expression)))
			r := httptest.NewRecorder()

			th.handleCalculate(r, req)

			if r.Code != tc.expectedStatus {
				t.Errorf("expected status code %d, got %d", tc.expectedStatus, r.Code)
			}
		})
	}
}

func TestTransportHttp_handleExpressions(t *testing.T) {
	defer func() {
		s.Err = nil
	}()

	cases := []struct {
		name           string
		method         string
		err            error
		expectedStatus int
	}{
		{
			name:           "success",
			method:         "GET",
			err:            nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "not found",
			method:         "GET",
			err:            errors.ErrNoExpressions,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "method not allowed",
			method:         "POST",
			err:            nil,
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s.Err = tc.err

			req := httptest.NewRequest(tc.method, "/api/v1/expressions", nil)
			r := httptest.NewRecorder()

			th.handleExpressions(r, req)

			if r.Code != tc.expectedStatus {
				t.Errorf("expected status code %d, got %d", tc.expectedStatus, r.Code)
			}
		})
	}
}

func TestTransportHttp_handleExpression(t *testing.T) {
	defer func() {
		s.Err = nil
	}()

	cases := []struct {
		name           string
		method         string
		err            error
		expectedStatus int
	}{
		{
			name:           "success",
			method:         "GET",
			err:            nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "not found",
			method:         "GET",
			err:            errors.ErrExpressionDoesNotExist,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "method not allowed",
			method:         "POST",
			err:            nil,
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s.Err = tc.err

			req := httptest.NewRequest(tc.method, "/api/v1/expressions/d8241c51-8782-42fb-9cb7-61ca519064d9", nil)
			r := httptest.NewRecorder()

			th.handleExpression(r, req)

			if r.Code != tc.expectedStatus {
				t.Errorf("expected status code %d, got %d", tc.expectedStatus, r.Code)
			}
		})
	}
}

func TestTransportHttp_handleGetTask(t *testing.T) {
	defer func() {
		s.Err = nil
	}()

	cases := []struct {
		name           string
		method         string
		err            error
		expectedStatus int
	}{
		{
			name:           "success",
			method:         "GET",
			err:            nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "not found",
			method:         "GET",
			err:            errors.ErrNoTasks,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s.Err = tc.err
			fmt.Println(s.Err, tc.err)

			req := httptest.NewRequest(tc.method, "/internal/tasks", nil)
			r := httptest.NewRecorder()

			th.handleGetTask(r, req)

			if r.Code != tc.expectedStatus {
				t.Errorf("expected status code %d, got %d", tc.expectedStatus, r.Code)
			}
		})
	}
}

func TestTransportHttp_handlePostTask(t *testing.T) {
	defer func() {
		s.Err = nil
	}()

	cases := []struct {
		name           string
		method         string
		body           string
		err            error
		expectedStatus int
	}{
		{
			name:           "success",
			method:         "POST",
			body:           `{"expression": "2+2"}`,
			err:            nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "bad request",
			method:         "POST",
			body:           "invalid expression",
			err:            errors.ErrInvalidExpression,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "not found",
			method:         "POST",
			body:           `{"expression": "2+2"}`,
			err:            errors.ErrExpressionDoesNotExist,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s.Err = tc.err

			req := httptest.NewRequest(tc.method, "/internal/tasks", bytes.NewReader([]byte(tc.body)))
			r := httptest.NewRecorder()

			th.handlePostResult(r, req)

			if r.Code != tc.expectedStatus {
				t.Errorf("expected status code %d, got %d", tc.expectedStatus, r.Code)
			}
		})
	}
}
