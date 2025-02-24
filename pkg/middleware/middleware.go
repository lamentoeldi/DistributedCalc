package middleware

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

func MwLogger(log *zap.Logger, next http.Handler) http.Handler {
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

func MwRecover(log *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("panic recovered", zap.Error(err.(error)))
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
