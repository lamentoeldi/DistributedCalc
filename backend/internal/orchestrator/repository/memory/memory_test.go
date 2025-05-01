package memory

import (
	"context"
	"github.com/distributed-calc/v1/internal/orchestrator/models"
	"github.com/google/uuid"
	"log"
	"testing"
)

func TestRepositoryMemory_Add(t *testing.T) {
	rep := NewRepositoryMemory()

	exp := &models.Expression{
		Id:     uuid.NewString(),
		Status: "testing",
		Result: 64,
	}

	err := rep.Add(context.Background(), exp)
	if err != nil {
		t.Error("Failed to add expression to memory")
	}
}

func TestRepositoryMemory_Get(t *testing.T) {
	rep := NewRepositoryMemory()

	exp := &models.Expression{
		Id:     uuid.NewString(),
		Status: "testing",
		Result: 64,
	}

	err := rep.Add(context.Background(), exp)
	if err != nil {
		t.Fatal("Failed to add expression to memory")
	}

	got, err := rep.Get(context.Background(), exp.Id)
	if err != nil {
		t.Error("Failed to get expression from memory")
	}

	log.Println(got)
}

func TestRepositoryMemory_GetAll(t *testing.T) {
	rep := NewRepositoryMemory()

	exp1 := &models.Expression{
		Id:     uuid.NewString(),
		Status: "testing",
		Result: 64,
	}

	exp2 := &models.Expression{
		Id:     uuid.NewString(),
		Status: "testing",
		Result: 64,
	}

	err := rep.Add(context.Background(), exp1)
	err = rep.Add(context.Background(), exp2)
	if err != nil {
		t.Fatal("Failed to add expression to memory")
	}

	got, err := rep.GetAll(context.Background())
	if err != nil {
		t.Error("Failed to get all expressions from memory")
	}

	for _, val := range got {
		log.Println(val)
	}
}
