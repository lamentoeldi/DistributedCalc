package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	e "github.com/distributed-calc/v1/internal/orchestrator/errors"
	"github.com/distributed-calc/v1/internal/orchestrator/models"
	"github.com/distributed-calc/v1/pkg/middleware"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

const (
	methodNotAllowed = "method not allowed"
)

type Service interface {
	Evaluate(ctx context.Context, expression string) (string, error)
	Get(ctx context.Context, id string) (*models.Expression, error)
	GetAll(ctx context.Context) ([]*models.Expression, error)

	GetTask(ctx context.Context) (*models.AgentTask, error)
	FinishTask(ctx context.Context, task *models.TaskResult) error
	Finalize(ctx context.Context, expID string, result float64) error
}

type TransportHttpConfig struct {
	Host string
	Port int
}

type TransportHttp struct {
	s      Service
	mux    *http.ServeMux
	server *http.Server
	log    *zap.Logger
	Host   string
	Port   int
}

func NewTransportHttp(s Service, log *zap.Logger, cfg *TransportHttpConfig) *TransportHttp {
	t := &TransportHttp{
		s:    s,
		mux:  http.NewServeMux(),
		log:  log,
		Host: cfg.Host,
		Port: cfg.Port,
	}

	t.mux.Handle("/api/v1/calculate", middleware.MwLogger(log, middleware.MwRecover(log, http.HandlerFunc(t.handleCalculate))))
	t.mux.Handle("/api/v1/expressions", middleware.MwLogger(log, middleware.MwRecover(log, http.HandlerFunc(t.handleExpressions))))
	t.mux.Handle("/api/v1/expressions/", middleware.MwLogger(log, middleware.MwRecover(log, http.HandlerFunc(t.handleExpression))))

	t.mux.Handle("/internal/ping", middleware.MwRecover(log, http.HandlerFunc(t.handlePing)))
	t.mux.Handle("/internal/task", middleware.MwRecover(log, http.HandlerFunc(t.handleTask)))

	return t
}

func (t *TransportHttp) handlePing(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, methodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (t *TransportHttp) handleCalculate(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value("request_id")
	log := t.log.With(zap.Any("request_id", reqID))

	if r.Method != "POST" {
		http.Error(w, methodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	var exp *models.CalculateRequest
	err := json.NewDecoder(r.Body).Decode(&exp)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := t.s.Evaluate(r.Context(), exp.Expression)
	if err != nil {
		log.Error(err.Error())

		switch {
		case errors.Is(err, e.ErrInvalidExpression):
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	data, err := json.Marshal(map[string]any{
		"id": id,
	})
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(data)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (t *TransportHttp) handleExpressions(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value("request_id")
	log := t.log.With(zap.Any("request_id", reqID))

	if r.Method != "GET" {
		http.Error(w, methodNotAllowed, http.StatusMethodNotAllowed)
	}

	exp, err := t.s.GetAll(r.Context())
	if err != nil {
		log.Error(err.Error())

		switch {
		case errors.Is(err, e.ErrNoExpressions):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	if len(exp) < 1 {
		http.Error(w, e.ErrNoExpressions.Error(), http.StatusNotFound)
		return
	}

	data, err := json.Marshal(map[string]any{
		"expressions": exp,
	})
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (t *TransportHttp) handleExpression(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value("request_id")
	log := t.log.With(zap.Any("request_id", reqID))

	if r.Method != "GET" {
		http.Error(w, methodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	routes := strings.Split(r.URL.Path, "/")

	// To ensure uuid is valid
	id, err := uuid.Parse(routes[len(routes)-1])
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	exp, err := t.s.Get(r.Context(), id.String())
	if err != nil {
		log.Error(err.Error(), zap.String("exp_id", id.String()))

		switch {
		case errors.Is(err, e.ErrExpressionDoesNotExist):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	data, err := json.Marshal(map[string]any{
		"expression": exp,
	})

	_, err = w.Write(data)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
	task, err := t.s.GetTask(r.Context())
	switch {
	case errors.Is(err, e.ErrNoTasks):
		t.log.Debug("no tasks", zap.Error(err))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	case err != nil:
		t.log.Error("error getting task", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(struct {
		Task *models.AgentTask `json:"task"`
	}{
		Task: task,
	})
	if err != nil {
		t.log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		t.log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (t *TransportHttp) handlePostResult(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var result *models.TaskResult

	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		t.log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = t.s.FinishTask(r.Context(), result)
	switch {
	case errors.Is(err, e.ErrExpressionDoesNotExist):
		t.log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	case err != nil:
		t.log.Error(err.Error())
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
