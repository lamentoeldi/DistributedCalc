package memory

import (
	"context"
	"fmt"
	"github.com/distributed-calc/v1/internal/orchestrator/errors"
	"github.com/distributed-calc/v1/internal/orchestrator/models"
	"sync"
)

type RepositoryMemory struct {
	expM  map[string]*models.Expression
	expMu sync.RWMutex

	taskM  map[string]*models.Task
	taskMu sync.RWMutex
}

func NewRepositoryMemory() *RepositoryMemory {
	return &RepositoryMemory{
		expM:  make(map[string]*models.Expression),
		taskM: make(map[string]*models.Task),
	}
}

func (rm *RepositoryMemory) Add(_ context.Context, exp *models.Expression) error {
	rm.expMu.Lock()
	defer rm.expMu.Unlock()

	rm.expM[exp.Id] = exp
	return nil
}

func (rm *RepositoryMemory) Get(_ context.Context, id string) (*models.Expression, error) {
	rm.expMu.RLock()
	val, ok := rm.expM[id]
	rm.expMu.RUnlock()
	if !ok {
		return nil, errors.ErrExpressionDoesNotExist
	}

	return val, nil
}

func (rm *RepositoryMemory) GetAll(_ context.Context) ([]*models.Expression, error) {
	expressions := make([]*models.Expression, 0)

	rm.expMu.RLock()
	for _, val := range rm.expM {
		expressions = append(expressions, val)
	}
	rm.expMu.RUnlock()

	if len(expressions) < 1 {
		return nil, errors.ErrNoExpressions
	}

	return expressions, nil
}

func (rm *RepositoryMemory) AddTasks(_ context.Context, tasks []*models.Task) error {
	rm.taskMu.Lock()
	defer rm.taskMu.Unlock()
	for _, t := range tasks {
		rm.taskM[t.ID] = t
	}
	return nil
}

func (rm *RepositoryMemory) GetTask(_ context.Context) (*models.Task, error) {
	rm.taskMu.RLock()
	defer rm.taskMu.RUnlock()
	for _, task := range rm.taskM {
		if task.Status == "ready" {
			delete(rm.taskM, task.ID)
			return task, nil
		}
	}
	return nil, fmt.Errorf("%w: no ready task found", errors.ErrNoTasks)
}

func (rm *RepositoryMemory) UpdateTask(_ context.Context, task *models.Task) error {
	rm.taskMu.Lock()
	defer rm.taskMu.Unlock()
	rm.taskM[task.ID] = task

	for id, t := range rm.taskM {
		if t.LeftID != nil && *t.LeftID == task.ID {
			t.LeftArg = task.Result
			t.LeftID = nil
			if t.LeftID == nil && t.RightID == nil {
				t.Status = "ready"
			}
			rm.taskM[id] = t
		}
		if t.RightID != nil && *t.RightID == task.ID {
			t.RightArg = task.Result
			t.RightID = nil
			if t.LeftID == nil && t.RightID == nil {
				t.Status = "ready"
			}
			rm.taskM[id] = t
		}
	}

	return nil
}

func (rm *RepositoryMemory) DeleteTasks(_ context.Context, expID string) error {
	rm.taskMu.Lock()
	defer rm.taskMu.Unlock()

	for _, task := range rm.taskM {
		if task.ExpID == expID {
			delete(rm.taskM, task.ID)
		}
	}

	return nil
}

func (rm *RepositoryMemory) Update(_ context.Context, exp *models.Expression) error {
	rm.expMu.Lock()
	defer rm.expMu.Unlock()

	rm.expM[exp.Id] = exp

	return nil
}
