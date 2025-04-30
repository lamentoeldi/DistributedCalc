package memory

import (
	"context"
	models2 "github.com/distributed-calc/v1/internal/orchestrator/models"
	"github.com/distributed-calc/v1/pkg/models"
	"sync"
)

type RepositoryMemory struct {
	exp     map[string]*models.Expression
	expMu   sync.RWMutex
	users   map[string]*models2.User
	usersMu sync.RWMutex
}

func NewRepositoryMemory() *RepositoryMemory {
	return &RepositoryMemory{
		exp:   make(map[string]*models.Expression),
		users: make(map[string]*models2.User),
	}
}

func (r *RepositoryMemory) Add(_ context.Context, exp *models.Expression) error {
	r.expMu.Lock()
	defer r.expMu.Unlock()

	r.exp[exp.Id] = exp
	return nil
}

func (r *RepositoryMemory) Get(_ context.Context, id string) (*models.Expression, error) {
	r.expMu.RLock()
	val, ok := r.exp[id]
	r.expMu.RUnlock()
	if !ok {
		return nil, models.ErrExpressionDoesNotExist
	}

	return val, nil
}

func (r *RepositoryMemory) GetAll(_ context.Context) ([]*models.Expression, error) {
	expressions := make([]*models.Expression, 0)

	r.expMu.RLock()
	for _, val := range r.exp {
		expressions = append(expressions, val)
	}
	r.expMu.RUnlock()

	if len(expressions) < 1 {
		return nil, models.ErrNoExpressions
	}

	return expressions, nil
}

func (r *RepositoryMemory) AddUser(_ context.Context, user *models2.User) error {
	r.usersMu.Lock()
	defer r.usersMu.Unlock()

	r.users[user.Username] = user
	return nil
}

func (r *RepositoryMemory) GetUser(_ context.Context, login string) (*models2.User, error) {
	r.usersMu.RLock()
	user, ok := r.users[login]
	r.usersMu.RUnlock()
	if !ok {
		return nil, models.ErrExpressionDoesNotExist
	}

	return user, nil
}
