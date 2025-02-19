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

type Service interface {
	StartEvaluation(ctx context.Context, expression string) error
	Get(ctx context.Context, id int) (*models.Expression, error)
	GetAll(ctx context.Context) ([]*models.Expression, error)
	FinishTask(ctx context.Context, result *models.TaskResult) error
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
	s      Service
	q      Queue[models.Task]
	mux    *http.ServeMux
	server *http.Server
	log    *zap.Logger
	Host   string
	Port   int
}

func NewTransportHttp(s Service, log *zap.Logger, cfg *TransportHttpConfig, queue Queue[models.Task]) *TransportHttp {
	t := &TransportHttp{
		s:    s,
		q:    queue,
		mux:  http.NewServeMux(),
		log:  log,
		Host: cfg.Host,
		Port: cfg.Port,
	}

	t.mux.HandleFunc("/internal/ping", t.handlePing)
	t.mux.Handle("/internal/task", mwLogger(log, http.HandlerFunc(t.handleTask)))

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
		t.handlePostResult(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (t *TransportHttp) handleGetTask(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value("request_id")
	log := t.log.With(zap.Any("request_id", reqID))

	task, err := t.q.Dequeue()
	switch {
	case errors.Is(err, models.ErrNoTasks):
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	case err != nil:
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(struct {
		Task *models.Task `json:"task"`
	}{
		Task: task,
	})
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (t *TransportHttp) handlePostResult(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value("request_id")
	log := t.log.With(zap.Any("request_id", reqID))

	defer r.Body.Close()

	var result *models.TaskResult

	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = t.s.FinishTask(r.Context(), result)
	switch {
	case errors.Is(err, models.ErrTaskDoesNotExist):
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	case err != nil:
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
