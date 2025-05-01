package mongo

import (
	"context"
	"github.com/distributed-calc/v1/internal/orchestrator/models"
	"github.com/distributed-calc/v1/pkg/mongo"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestRepository_Add(t *testing.T) {
	cases := []struct {
		name    string
		exp     *models.Expression
		wantErr bool
	}{
		{
			name: "success",
			exp: &models.Expression{
				Id:     uuid.NewString(),
				UserID: uuid.NewString(),
				Result: 0,
				Status: "pending",
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := mongo.NewMongoConfig()
	if err != nil {
		t.Fatal(err)
	}

	client, err := mongo.NewMongoClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		client.Database(cfg.DBName).Collection(collExp).Drop(ctx)
	})

	repo := NewMongoRepository(cfg, client)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			err := repo.Add(ctx, tc.exp)
			if tc.wantErr == false && err != nil {
				t.Errorf("expected no error got %v", err)
			}

			if tc.wantErr == true && err == nil {
				t.Errorf("expected error got none")
			}
		})
	}
}

func TestRepository_AddUser(t *testing.T) {
	cases := []struct {
		name    string
		user    *models.User
		wantErr bool
	}{
		{
			name: "success",
			user: &models.User{
				Id:             uuid.NewString(),
				Username:       "test:user:1",
				HashedPassword: []byte("test"),
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := mongo.NewMongoConfig()
	if err != nil {
		t.Fatal(err)
	}

	client, err := mongo.NewMongoClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		client.Database(cfg.DBName).Collection(collUsers).Drop(ctx)
	})

	repo := NewMongoRepository(cfg, client)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			err := repo.AddUser(ctx, tc.user)
			if tc.wantErr == false && err != nil {
				t.Errorf("expected no error got %v", err)
			}

			if tc.wantErr == true && err == nil {
				t.Errorf("expected error got none")
			}
		})
	}
}

func TestRepository_AddTasks(t *testing.T) {
	id := uuid.NewString()

	cases := []struct {
		name    string
		tasks   []*models.Task
		wantErr bool
	}{
		{
			name: "success",
			tasks: []*models.Task{
				{
					ID:       uuid.NewString(),
					ExpID:    uuid.NewString(),
					Op:       "+",
					LeftArg:  10,
					RightArg: 20,
					LeftID:   nil,
					RightID:  &id,
				},
				{
					ID:       uuid.NewString(),
					ExpID:    uuid.NewString(),
					Op:       "-",
					LeftArg:  10,
					RightArg: 20,
					LeftID:   nil,
					RightID:  nil,
				},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := mongo.NewMongoConfig()
	if err != nil {
		t.Fatal(err)
	}

	client, err := mongo.NewMongoClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		client.Database(cfg.DBName).Collection(collTasks).Drop(ctx)
	})

	repo := NewMongoRepository(cfg, client)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			err := repo.AddTasks(ctx, tc.tasks)
			if tc.wantErr == false && err != nil {
				t.Errorf("expected no error got %v", err)
			}

			if tc.wantErr == true && err == nil {
				t.Errorf("expected error got none")
			}
		})
	}
}

func TestRepository_Get(t *testing.T) {
	cases := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "success",
			id:      "test:get:1",
			wantErr: false,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := mongo.NewMongoConfig()
	if err != nil {
		t.Fatal(err)
	}

	client, err := mongo.NewMongoClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		client.Database(cfg.DBName).Collection(collExp).Drop(ctx)
	})

	repo := NewMongoRepository(cfg, client)

	err = repo.Add(ctx, &models.Expression{Id: "test:get:1", UserID: "test:get:1"})
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			_, err := repo.Get(ctx, tc.id)
			if tc.wantErr == false && err != nil {
				t.Errorf("expected no error got %v", err)
			}

			if tc.wantErr == true && err == nil {
				t.Errorf("expected error got none")
			}
		})
	}
}

func TestRepository_GetAll(t *testing.T) {
	cases := []struct {
		name    string
		userID  string
		cursor  string
		limit   int64
		wantErr bool
	}{
		{
			name:   "success",
			userID: "test:get:all:1",
			limit:  10,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := mongo.NewMongoConfig()
	if err != nil {
		t.Fatal(err)
	}

	client, err := mongo.NewMongoClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		client.Database(cfg.DBName).Collection(collExp).Drop(ctx)
	})

	repo := NewMongoRepository(cfg, client)

	err = repo.Add(ctx, &models.Expression{Id: "test:get:all:1", UserID: "test:get:all:1"})
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			_, err := repo.GetAll(ctx, tc.userID, tc.cursor, tc.limit)
			if tc.wantErr == false && err != nil {
				t.Errorf("expected no error got %v", err)
			}

			if tc.wantErr == true && err == nil {
				t.Errorf("expected error got none")
			}
		})
	}
}

func TestRepository_GetUser(t *testing.T) {
	cases := []struct {
		name     string
		userName string
		wantErr  bool
	}{
		{
			name:     "success",
			userName: "test:get:user:1",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := mongo.NewMongoConfig()
	if err != nil {
		t.Fatal(err)
	}

	client, err := mongo.NewMongoClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		client.Database(cfg.DBName).Collection(collUsers).Drop(ctx)
	})

	repo := NewMongoRepository(cfg, client)

	err = repo.AddUser(ctx, &models.User{Id: "test:get:user:1", Username: "test:get:user:1"})
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			_, err := repo.GetUser(ctx, tc.userName)
			if tc.wantErr == false && err != nil {
				t.Errorf("expected no error got %v", err)
			}

			if tc.wantErr == true && err == nil {
				t.Errorf("expected error got none")
			}
		})
	}
}

func TestRepository_GetTask(t *testing.T) {
	cases := []struct {
		name    string
		wantErr bool
	}{
		{
			name: "success",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := mongo.NewMongoConfig()
	if err != nil {
		t.Fatal(err)
	}

	client, err := mongo.NewMongoClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		client.Database(cfg.DBName).Collection(collTasks).Drop(ctx)
	})

	repo := NewMongoRepository(cfg, client)

	err = repo.AddTasks(ctx, []*models.Task{
		{
			ID:     uuid.NewString(),
			ExpID:  uuid.NewString(),
			Status: "ready",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			_, err := repo.GetTask(ctx)
			if tc.wantErr == false && err != nil {
				t.Errorf("expected no error got %v", err)
			}

			if tc.wantErr == true && err == nil {
				t.Errorf("expected error got none")
			}
		})
	}
}

func TestRepository_Update(t *testing.T) {
	cases := []struct {
		name    string
		exp     *models.Expression
		wantErr bool
	}{
		{
			name: "success",
			exp: &models.Expression{
				Id:     "test:update:1",
				Status: "pending",
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := mongo.NewMongoConfig()
	if err != nil {
		t.Fatal(err)
	}

	client, err := mongo.NewMongoClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		client.Database(cfg.DBName).Collection(collExp).Drop(ctx)
	})

	repo := NewMongoRepository(cfg, client)

	err = repo.Add(ctx, &models.Expression{Id: "test:update:1", UserID: "test:update:1"})
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			err := repo.Update(ctx, tc.exp)
			if tc.wantErr == false && err != nil {
				t.Errorf("expected no error got %v", err)
			}

			if tc.wantErr == true && err == nil {
				t.Errorf("expected error got none")
			}
		})
	}
}

func TestRepository_UpdateTask(t *testing.T) {
	cases := []struct {
		name    string
		task    *models.Task
		wantErr bool
	}{
		{
			name: "success",
			task: &models.Task{
				ID: "test:update:task:1",
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := mongo.NewMongoConfig()
	if err != nil {
		t.Fatal(err)
	}

	client, err := mongo.NewMongoClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		client.Database(cfg.DBName).Collection(collTasks).Drop(ctx)
	})

	repo := NewMongoRepository(cfg, client)

	leftID := "test:update:task:1"
	err = repo.AddTasks(ctx, []*models.Task{
		{
			ID:     uuid.NewString(),
			ExpID:  "test:2",
			Status: "pending",
			LeftID: &leftID,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			err := repo.UpdateTask(ctx, tc.task)
			if tc.wantErr == false && err != nil {
				t.Errorf("expected no error got %v", err)
			}

			if tc.wantErr == true && err == nil {
				t.Errorf("expected error got none")
			}
		})
	}
}

func TestRepository_DeleteTasks(t *testing.T) {
	cases := []struct {
		name    string
		expID   string
		wantErr bool
	}{
		{
			name:  "success",
			expID: "test:delete:tasks:1",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := mongo.NewMongoConfig()
	if err != nil {
		t.Fatal(err)
	}

	client, err := mongo.NewMongoClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		client.Database(cfg.DBName).Collection(collTasks).Drop(ctx)
	})

	repo := NewMongoRepository(cfg, client)

	err = repo.AddTasks(ctx, []*models.Task{
		{
			ID:     uuid.NewString(),
			ExpID:  "test:delete:tasks:1",
			Status: "ready",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			err := repo.DeleteTasks(ctx, tc.expID)
			if tc.wantErr == false && err != nil {
				t.Errorf("expected no error got %v", err)
			}

			if tc.wantErr == true && err == nil {
				t.Errorf("expected error got none")
			}
		})
	}
}
