package memory

import (
	"DistributedCalc/pkg/models"
	"context"
	"sync"
)

type RepositoryMemory struct {
	m sync.Map
}

func NewRepositoryMemory() *RepositoryMemory {
	return &RepositoryMemory{}
}

func (r *RepositoryMemory) Add(_ context.Context, exp *models.Expression) error {
	r.m.Store(exp.Id, exp)
	return nil
}

func (r *RepositoryMemory) Get(_ context.Context, id int) (*models.Expression, error) {
	val, ok := r.m.Load(id)
	if !ok {
		return nil, models.ErrTaskDoesNotExist
	}

	if exp, ok := val.(*models.Expression); ok {
		return exp, nil
	}

	return nil, models.ErrTaskDoesNotExist
}

func (r *RepositoryMemory) GetAll(_ context.Context) ([]*models.Expression, error) {
	expressions := make([]*models.Expression, 0)
	var err error

	r.m.Range(func(_, value any) bool {
		if exp, ok := value.(*models.Expression); ok {
			expressions = append(expressions, exp)
			return true
		}

		err = models.ErrTaskDoesNotExist
		return false
	})

	return expressions, err
}
