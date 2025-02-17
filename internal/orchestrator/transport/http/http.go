package http

import (
	"DistributedCalc/pkg/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

type Calculator interface {
	Evaluate(ctx context.Context, expression string) (float64, error)
}

type Queue[T any] interface {
	Enqueue(obj *T)
	Dequeue() (*T, error)
}

type TransportHttpConfig struct {
	Host string
	Port int
}

type TransportHttp struct {
	Calculator
	Queue[models.Task]
	mux    *http.ServeMux
	server *http.Server
	log    *zap.Logger
	Host   string
	Port   int
}

func NewTransportHttp(calculator Calculator, log *zap.Logger, cfg *TransportHttpConfig, queue Queue[models.Task]) *TransportHttp {
	t := &TransportHttp{
		Calculator: calculator,
		Queue:      queue,
		mux:        http.NewServeMux(),
		log:        log,
		Host:       cfg.Host,
		Port:       cfg.Port,
	}

	t.mux.HandleFunc("/internal/ping", t.handlePing)
	t.mux.HandleFunc("/internal/task", t.handleTask)

	return t
}

func (t *TransportHttp) handlePing(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (t *TransportHttp) handleTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		t.handleGetTask(w, r)
	case "POST":
		http.Error(w, "not implemented", http.StatusNotImplemented)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (t *TransportHttp) handleGetTask(w http.ResponseWriter, r *http.Request) {
	task, err := t.Queue.Dequeue()
	switch {
	case errors.Is(err, models.ErrNoTasks):
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(struct {
		Task *models.Task `json:"task"`
	}{
		Task: task,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (t *TransportHttp) Run() {
	addr := fmt.Sprintf("%s:%d", t.Host, t.Port)
	t.server = &http.Server{
		Addr:    addr,
		Handler: t.mux,
	}

	go func() {
		t.log.Info("Starting server on " + addr)
		err := t.server.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				t.log.Info(err.Error())
				return
			}
			t.log.Fatal(err.Error())
		}
	}()
}

func (t *TransportHttp) Shutdown(ctx context.Context) {
	t.log.Info("Shutting down...")
	err := t.server.Shutdown(ctx)
	if err != nil {
		t.log.Error(err.Error())
	}
}
