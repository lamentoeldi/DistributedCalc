package orchestrator

import (
	"DistributedCalc/pkg/models"
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
)

func handleGetTask(w http.ResponseWriter, _ *http.Request) {
	task := &models.Task{
		Id:            rand.Int(),
		Arg1:          3.14,
		Arg2:          2,
		Operation:     "+",
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

	o := NewOrchestrator(&http.Client{}, ts.URL, 3)
	result := &models.TaskResult{
		Id:     90131329,
		Result: 5,
	}

	err := o.PostResult(context.Background(), result)
	if err != nil {
		t.Error("failed to post result")
	}
}
