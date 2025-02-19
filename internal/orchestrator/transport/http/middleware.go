package http

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

func mwLogger(log *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New()

		log.Info("request",
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.String("request_id", id.String()),
		)

		newCtx := context.WithValue(r.Context(), "request_id", id)

		newR := r.WithContext(newCtx)

		next.ServeHTTP(w, newR)
	})
}
