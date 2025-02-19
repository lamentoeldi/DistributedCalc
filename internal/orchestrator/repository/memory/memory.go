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

func (r *RepositoryMemory) Add(ctx context.Context, exp *models.Expression) error {
	//TODO implement me
	return nil
}

func (r *RepositoryMemory) Get(ctx context.Context, id int) (*models.Expression, error) {
	//TODO implement me
	return nil, nil
}

func (r *RepositoryMemory) GetAll(ctx context.Context) ([]*models.Expression, error) {
	//TODO implement me
	return nil, nil
}
