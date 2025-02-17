package http

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

type Calculator interface {
	Evaluate(ctx context.Context, expression string) (float64, error)
}

type TransportHttpConfig struct {
	Host string
	Port int
}

type TransportHttp struct {
	Calculator
	mux    *http.ServeMux
	server *http.Server
	log    *zap.Logger
	Host   string
	Port   int
}

func NewTransportHttp(calculator Calculator, log *zap.Logger, cfg *TransportHttpConfig) *TransportHttp {
	t := &TransportHttp{
		Calculator: calculator,
		mux:        http.NewServeMux(),
		log:        log,
		Host:       cfg.Host,
		Port:       cfg.Port,
	}

	t.mux.HandleFunc("/internal/ping", t.handlePing)
	t.mux.HandleFunc("internal/task", t.handleTask)

	return t
}

func (t *TransportHttp) handlePing(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (t *TransportHttp) handleTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		//get logic
		w.WriteHeader(http.StatusNotImplemented)
	case "POST":
		//post logic
		w.WriteHeader(http.StatusNotImplemented)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (t *TransportHttp) Run() {
	addr := fmt.Sprintf("%s:%d", t.Host, t.Port)
	t.server = &http.Server{
		Addr:    addr,
		Handler: t.mux,
	}

	go func() {
		err := t.server.ListenAndServe()
		if err != nil {
			t.log.Fatal(err.Error())
		}
	}()
}

func (t *TransportHttp) Shutdown(ctx context.Context) {
	err := t.server.Shutdown(ctx)
	if err != nil {
		t.log.Error(err.Error())
	}
}
