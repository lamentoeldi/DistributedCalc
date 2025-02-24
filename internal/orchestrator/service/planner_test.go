package service

import (
	"DistributedCalc/internal/orchestrator/adapters/queue"
	"DistributedCalc/internal/orchestrator/config"
	"DistributedCalc/pkg/models"
	"context"
	"errors"
	"testing"
	"time"
)

func TestPlannerChan_PlanTask(t *testing.T) {
	cfg := &config.Config{
		AdditionTime:       1 * time.Millisecond,
		SubtractionTime:    1 * time.Millisecond,
		MultiplicationTime: 1 * time.Millisecond,
		DivisionTime:       1 * time.Millisecond,
	}
	q := queue.NewQueueChan[models.Task](64)
	planner := NewPlannerChan(cfg, q)

	cases := []struct {
		name    string
		task    *models.Task
		wantErr error
	}{
		{
			name:    "valid addition task",
			task:    &models.Task{Arg1: 2, Arg2: 3, Operation: "+"},
			wantErr: nil,
		},
		{
			name:    "division by zero",
			task:    &models.Task{Arg1: 2, Arg2: 0, Operation: "/"},
			wantErr: models.ErrDivisionByZero,
		},
		{
			name:    "unsupported operation",
			task:    &models.Task{Arg1: 2, Arg2: 3, Operation: "^"},
			wantErr: models.ErrUnsupportedOperation,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := planner.PlanTask(context.Background(), tc.task)
			if tc.wantErr != nil {
				if err == nil || !errors.Is(err, tc.wantErr) {
					t.Errorf("expected error %v, got %v", tc.wantErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestPlannerChan_FinishTask(t *testing.T) {
	// It is impossible to unit test FinishTask() due to its specifics
}
