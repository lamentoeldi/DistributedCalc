package orchestrator

import (
	"context"
	"encoding/json"
	"github.com/distributed-calc/v1/internal/agent/models"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"
)

func handleGetTask(w http.ResponseWriter, _ *http.Request) {
	id, _ := uuid.NewV7()

	task := &models.AgentTask{
		Id:            id.String(),
		LeftArg:       3.14,
		RightArg:      2,
		Op:            "+",
		OperationTime: 10,
	}

	w.Header().Set("Content-Type", "application/json")
	data, _ := json.Marshal(task)
	_, _ = w.Write(data)
}

func handlePostResult(w http.ResponseWriter, r *http.Request) {
	var result models.TaskResult
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handlePing(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestOrchestrator_Ping(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handlePing))
	defer ts.Close()

	o := NewOrchestrator(&http.Client{}, ts.URL, 3)
	err := o.Ping()
	if err != nil {
		t.Error("failed to ping")
	}
}

func TestOrchestrator_GetTask(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handleGetTask))
	defer ts.Close()

	o := NewOrchestrator(&http.Client{}, ts.URL, 3)
	_, err := o.GetTask(context.Background())
	if err != nil {
		t.Error("failed to get task")
	}
}

func TestOrchestrator_PostResult(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handlePostResult))
	defer ts.Close()

	id, _ := uuid.NewV7()

	o := NewOrchestrator(&http.Client{}, ts.URL, 3)
	result := &models.TaskResult{
		Id:     id.String(),
		Result: 5,
	}

	err := o.PostResult(context.Background(), result)
	if err != nil {
		t.Error("failed to post result")
	}
}
