package memory

import (
	"DistributedCalc/pkg/models"
	"context"
	"sync"
)

type RepositoryMemory struct {
	m  map[string]*models.Expression
	mu sync.Mutex
}

func NewRepositoryMemory() *RepositoryMemory {
	return &RepositoryMemory{
		m: make(map[string]*models.Expression),
	}
}

func (r *RepositoryMemory) Add(_ context.Context, exp *models.Expression) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.m[exp.Id] = exp
	return nil
}

func (r *RepositoryMemory) Get(_ context.Context, id string) (*models.Expression, error) {
	r.mu.Lock()
	val, ok := r.m[id]
	r.mu.Unlock()
	if !ok {
		return nil, models.ErrExpressionDoesNotExist
	}

	return val, nil
}

func (r *RepositoryMemory) GetAll(_ context.Context) ([]*models.Expression, error) {
	expressions := make([]*models.Expression, 0)

	r.mu.Lock()
	for _, val := range r.m {
		expressions = append(expressions, val)
	}
	r.mu.Unlock()

	if len(expressions) < 1 {
		return nil, models.ErrNoExpressions
	}

	return expressions, nil
}
