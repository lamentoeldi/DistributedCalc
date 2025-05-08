package mongo

import (
	"context"
	"testing"
	"time"
)

func TestNewMongoClient(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 1*time.Second)
	defer cancel()

	_, err := NewMongoClient(ctx)
	if err != nil {
		t.Fatal(err)
	}
}
