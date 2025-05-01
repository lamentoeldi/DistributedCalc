package mongo

import (
	"context"
	"fmt"
	"github.com/distributed-calc/v1/internal/orchestrator/models"
	mongo2 "github.com/distributed-calc/v1/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	collUsers = "users"
	collExp   = "expressions"
)

type Repository struct {
	client *mongo.Client
	cfg    mongo2.Config
}

func NewMongoRepository(cfg mongo2.Config) *Repository {
	return &Repository{
		cfg: cfg,
	}
}

func (r Repository) AddUser(ctx context.Context, user *models.User) error {
	_, err := r.client.
		Database(r.cfg.DBName).
		Collection(collUsers).
		InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to add user: %w", err)
	}

	return nil
}

func (r Repository) GetUser(ctx context.Context, login string) (*models.User, error) {
	res := r.client.
		Database(r.cfg.DBName).
		Collection(collUsers).
		FindOne(ctx, bson.M{"login": login})
	if res.Err() != nil {
		return nil, fmt.Errorf("failed to get user: %w", res.Err())
	}

	var user models.User
	err := res.Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r Repository) Add(ctx context.Context, exp *models.Expression) error {
	_, err := r.client.
		Database(r.cfg.DBName).
		Collection(collExp).
		InsertOne(ctx, &exp)
	if err != nil {
		return fmt.Errorf("failed to add exp: %w", err)
	}

	return nil
}

func (r Repository) Get(ctx context.Context, id string) (*models.Expression, error) {
	res := r.client.
		Database(r.cfg.DBName).
		Collection(collExp).
		FindOne(ctx, bson.M{"_id": id})
	if res.Err() != nil {
		return nil, fmt.Errorf("failed to get exp: %w", res.Err())
	}

	var exp models.Expression
	err := res.Decode(&exp)
	if err != nil {
		return nil, fmt.Errorf("failed to get exp: %w", err)
	}

	return &exp, nil
}

func (r Repository) GetAll(ctx context.Context) ([]*models.Expression, error) {
	//TODO implement me
	panic("implement me")
}

func (r Repository) Update(ctx context.Context, exp *models.Expression) error {
	//TODO implement me
	panic("implement me")
}
