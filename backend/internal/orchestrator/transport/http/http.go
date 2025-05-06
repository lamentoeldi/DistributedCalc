package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	e "github.com/distributed-calc/v1/internal/orchestrator/errors"
	"github.com/distributed-calc/v1/internal/orchestrator/models"
	"github.com/distributed-calc/v1/pkg/middleware"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	methodNotAllowed       = "method not allowed"
	requestTimeout         = 5 * time.Second
	defaultLimit     int64 = 10
)

type Service interface {
	Evaluate(ctx context.Context, expression, userID string) (string, error)
	Get(ctx context.Context, id, userID string) (*models.Expression, error)
	GetAll(ctx context.Context, userID, cursor string, limit int64) ([]*models.Expression, error)

	GetTask(ctx context.Context) (*models.AgentTask, error)
	FinishTask(ctx context.Context, result *models.TaskResult) error

	Register(ctx context.Context, creds *models.UserCredentials) error
	Login(ctx context.Context, creds *models.UserCredentials) (*models.JWTTokens, error)
	GetUserID(_ context.Context, token string) (string, error)

	GetUser(ctx context.Context, id string) (*models.UserView, error)

	middleware.Auth
}

type Config struct {
	Host string
	Port int
}

type Server struct {
	s      Service
	mux    *http.ServeMux
	server *http.Server
	log    *zap.Logger
	Host   string
	Port   int
}

func NewServer(cfg *Config, s Service, log *zap.Logger) *Server {
	t := &Server{
		s:    s,
		mux:  http.NewServeMux(),
		log:  log,
		Host: cfg.Host,
		Port: cfg.Port,
	}

	t.mux.
		Handle(
			"/api/v1/calculate",
			middleware.MwLogger(log,
				middleware.MwRecover(log,
					middleware.MwAuth(log, s, http.HandlerFunc(t.handleCalculate)))))
	t.mux.
		Handle(
			"/api/v1/expressions",
			middleware.MwLogger(log,
				middleware.MwRecover(log,
					middleware.MwAuth(log, s, http.HandlerFunc(t.handleExpressions)))))
	t.mux.
		Handle(
			"/api/v1/expressions/",
			middleware.MwLogger(log,
				middleware.MwRecover(log,
					middleware.MwAuth(log, s, http.HandlerFunc(t.handleExpression)))))

	t.mux.
		Handle(
			"/api/v1/register",
			middleware.MwLogger(log, middleware.MwRecover(log, http.HandlerFunc(t.handleRegister))))
	t.mux.
		Handle(
			"/api/v1/login",
			middleware.MwLogger(log, middleware.MwRecover(log, http.HandlerFunc(t.handleLogin))))
	t.mux.
		Handle(
			"/api/v1/authorize",
			middleware.MwLogger(log, middleware.MwRecover(log, http.HandlerFunc(t.handleAuthorize))))

	return t
}

func (t *Server) handlePing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, methodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (t *Server) handleCalculate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	if r.Method != http.MethodPost {
		http.Error(w, methodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	var exp *models.CalculateRequest
	err := json.NewDecoder(r.Body).Decode(&exp)
	if err != nil {
		t.log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authorization := r.Header.Get("Authorization")
	if len(authorization) < len("Bearer ") {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	accessToken := strings.TrimPrefix(authorization, "Bearer ")

	userID, err := t.s.GetUserID(ctx, accessToken)
	if err != nil {
		t.log.Error("failed to get user id", zap.Error(err))
		return
	}

	expID, err := t.s.Evaluate(ctx, exp.Expression, userID)
	if err != nil {
		t.log.Error(err.Error())

		switch {
		case errors.Is(err, e.ErrInvalidExpression):
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	data, err := json.Marshal(map[string]any{
		"id": expID,
	})
	if err != nil {
		t.log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(data)
}

func (t *Server) handleExpressions(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	if r.Method != http.MethodGet {
		http.Error(w, methodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	authorization := r.Header.Get("Authorization")
	if len(authorization) < len("Bearer ") {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	accessToken := strings.TrimPrefix(authorization, "Bearer ")

	userID, err := t.s.GetUserID(ctx, accessToken)
	if err != nil {
		t.log.Error("failed to get user id", zap.Error(err))
		return
	}

	cursor := r.URL.Query().Get("cursor")
	limitQuery := r.URL.Query().Get("limit")

	var limit int64
	limit, err = strconv.ParseInt(limitQuery, 10, 64)
	if err != nil {
		limit = defaultLimit
	}

	exp, err := t.s.GetAll(ctx, userID, cursor, limit)
	if err != nil {
		t.log.Error(err.Error())

		switch {
		case errors.Is(err, sql.ErrNoRows):
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
		t.log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(data)
}

func (t *Server) handleExpression(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	if r.Method != http.MethodGet {
		http.Error(w, methodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	routes := strings.Split(r.URL.Path, "/")

	// To ensure uuid is valid
	id, err := uuid.Parse(routes[len(routes)-1])
	if err != nil {
		t.log.Error(err.Error())
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	authorization := r.Header.Get("Authorization")
	if len(authorization) < len("Bearer ") {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	accessToken := strings.TrimPrefix(authorization, "Bearer ")

	userID, err := t.s.GetUserID(ctx, accessToken)
	if err != nil {
		t.log.Error("failed to get user id", zap.Error(err))
		return
	}

	exp, err := t.s.Get(ctx, id.String(), userID)
	if err != nil {
		t.log.Error(err.Error(), zap.String("exp_id", id.String()))

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

	_, _ = w.Write(data)
}

func (t *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	if r.Method != http.MethodPost {
		http.Error(w, methodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var creds models.UserCredentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		t.log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = t.s.Register(ctx, &creds)
	if err != nil {
		t.log.Error(err.Error())

		switch {
		case errors.Is(err, e.ErrBadRequest):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, e.ErrConflict):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (t *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	if r.Method != http.MethodPost {
		http.Error(w, methodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var creds models.UserCredentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		t.log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokens, err := t.s.Login(ctx, &creds)
	if err != nil {
		t.log.Error(err.Error())

		switch {
		case errors.Is(err, e.ErrUnauthorized):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	err = json.NewEncoder(w).Encode(tokens)
	if err != nil {
		t.log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (t *Server) handleAuthorize(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	if r.Method != http.MethodPost {
		http.Error(w, methodNotAllowed, http.StatusMethodNotAllowed)
	}

	defer r.Body.Close()

	authorization := r.Header.Get("Authorization")
	if len(authorization) < len("Bearer ") {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	accessToken := strings.TrimPrefix(authorization, "Bearer ")

	userID, err := t.s.GetUserID(ctx, accessToken)
	if err != nil {
		t.log.Error("failed to get user id", zap.Error(err))
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := t.s.GetUser(ctx, userID)
	if err != nil {
		t.log.Error("failed to get user", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		t.log.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (t *Server) Run() {
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

func (t *Server) Shutdown(ctx context.Context) {
	t.log.Info("Shutting down...")
	err := t.server.Shutdown(ctx)
	if err != nil {
		t.log.Error(err.Error())
	}
}
