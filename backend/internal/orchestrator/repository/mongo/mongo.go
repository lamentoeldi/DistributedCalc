package mongo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	errors2 "github.com/distributed-calc/v1/internal/orchestrator/errors"
	"github.com/distributed-calc/v1/internal/orchestrator/models"
	mongo2 "github.com/distributed-calc/v1/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collUsers = "users"
	collExp   = "expressions"
	collTasks = "tasks"
)

type Repository struct {
	client *mongo.Client
	cfg    *mongo2.Config
}

func NewMongoRepository(cfg *mongo2.Config, client *mongo.Client) *Repository {
	return &Repository{
		cfg:    cfg,
		client: client,
	}
}

func (r *Repository) AddUser(ctx context.Context, user *models.User) error {
	_, err := r.client.
		Database(r.cfg.DBName).
		Collection(collUsers).
		InsertOne(ctx, user)
	if err != nil {
		var e mongo.WriteException
		if errors.As(err, &e) {
			if e.HasErrorCode(11000) {
				return fmt.Errorf("failed to add user: %w", errors2.ErrUserAlreadyExists)
			}
		}
		return fmt.Errorf("failed to add user: %w", err)
	}

	return nil
}

func (r *Repository) GetUser(ctx context.Context, login string) (*models.User, error) {
	res := r.client.
		Database(r.cfg.DBName).
		Collection(collUsers).
		FindOne(ctx, bson.M{"username": login})
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

func (r *Repository) Add(ctx context.Context, exp *models.Expression) error {
	_, err := r.client.
		Database(r.cfg.DBName).
		Collection(collExp).
		InsertOne(ctx, &exp)
	if err != nil {
		return fmt.Errorf("failed to add exp: %w", err)
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, id string) (*models.Expression, error) {
	res := r.client.
		Database(r.cfg.DBName).
		Collection(collExp).
		FindOne(ctx, bson.M{"_id": id})
	err := res.Err()
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return nil, fmt.Errorf("failed to get exp: %w: expression not found", sql.ErrNoRows)
		default:
			return nil, fmt.Errorf("failed to get exp: %w", res.Err())
		}
	}

	var exp models.Expression
	err = res.Decode(&exp)
	if err != nil {
		return nil, fmt.Errorf("failed to get exp: %w", err)
	}

	return &exp, nil
}

func (r *Repository) GetAll(ctx context.Context, userID, cursor string, limit int64) ([]*models.Expression, error) {
	res, err := r.client.
		Database(r.cfg.DBName).
		Collection(collExp).
		Find(ctx, bson.M{
			"_id": bson.M{
				"$gt": cursor,
			},
			"user_id": userID,
		},
			options.Find().SetLimit(limit),
		)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("expressions not found: %w", sql.ErrNoRows)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get expressions for user %s: %w", userID, err)
	}

	expressions := make([]*models.Expression, 0, res.RemainingBatchLength())
	for res.Next(ctx) {
		var exp models.Expression
		err := res.Decode(&exp)
		if err != nil {
			return nil, fmt.Errorf("failed to get expressions for user %s: %w", userID, err)
		}

		expressions = append(expressions, &exp)
	}

	if len(expressions) < 1 {
		return nil, fmt.Errorf("expressions not found: %w", sql.ErrNoRows)
	}

	return expressions, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	res := r.client.
		Database(r.cfg.DBName).
		Collection(collUsers).
		FindOne(ctx, bson.M{"_id": id})
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

func (r *Repository) Update(ctx context.Context, exp *models.Expression) error {
	_, err := r.client.
		Database(r.cfg.DBName).
		Collection(collExp).
		UpdateByID(ctx, exp.Id, bson.M{
			"$set": bson.M{
				"result": exp.Result,
				"status": exp.Status},
		})
	if err != nil {
		return fmt.Errorf("failed to update exp: %w", err)
	}

	return nil
}

func (r *Repository) AddTasks(ctx context.Context, tasks []*models.Task) error {
	docs := make([]interface{}, 0, len(tasks))
	for _, task := range tasks {
		docs = append(docs, task)
	}

	_, err := r.client.
		Database(r.cfg.DBName).
		Collection(collTasks).
		InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("failed to add tasks: %w", err)
	}

	return nil
}

func (r *Repository) GetTask(ctx context.Context) (*models.Task, error) {
	res := r.client.
		Database(r.cfg.DBName).
		Collection(collTasks).
		FindOne(ctx, bson.M{"status": "ready"})
	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("task not found: %w", sql.ErrNoRows)
	}

	if res.Err() != nil {
		return nil, fmt.Errorf("failed to get task: %w", res.Err())
	}

	var task models.Task
	err := res.Decode(&task)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &task, nil
}

func (r *Repository) UpdateTask(ctx context.Context, task *models.Task) error {
	session, err := r.client.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	err = session.StartTransaction()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	client := session.Client()

	_, err = client.
		Database(r.cfg.DBName).
		Collection(collTasks).DeleteOne(ctx, bson.M{"_id": task.ID})
	//UpdateByID(ctx, task.ID, bson.M{
	//	"$set": bson.M{
	//		"result": task.Result,
	//		"status": task.Status,
	//	},
	//})
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	_, err = client.
		Database(r.cfg.DBName).
		Collection(collTasks).
		UpdateMany(ctx, bson.M{"left_id": task.ID}, bson.M{
			"$unset": bson.M{
				"left_id": "",
			},
			"$set": bson.M{
				"left_arg": task.Result,
			},
		})
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	_, err = client.
		Database(r.cfg.DBName).
		Collection(collTasks).
		UpdateMany(ctx, bson.M{"right_id": task.ID}, bson.M{
			"$unset": bson.M{
				"right_id": "",
			},
			"$set": bson.M{
				"right_arg": task.Result,
			},
		})
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	_, err = client.
		Database(r.cfg.DBName).
		Collection(collTasks).
		UpdateMany(ctx,
			bson.M{
				"left_id":  bson.M{"$exists": false},
				"right_id": bson.M{"$exists": false},
			},
			bson.M{
				"$set": bson.M{
					"status": "ready",
				},
			})
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	err = session.CommitTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to update task: failed to commmit transaction %w", err)
	}

	return nil
}

func (r *Repository) DeleteTasks(ctx context.Context, expID string) error {
	_, err := r.client.
		Database(r.cfg.DBName).
		Collection(collTasks).
		DeleteMany(ctx, bson.M{
			"exp_id": expID,
		})
	if err != nil {
		return fmt.Errorf("failed to delete tasks: %w", err)
	}

	return nil
}
