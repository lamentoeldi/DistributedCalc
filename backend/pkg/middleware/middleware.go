package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

var (
	ErrTokenWasRevoked = fmt.Errorf("token was revoked")
)

type Auth interface {
	VerifyJWT(ctx context.Context, jwt string) error
	RefreshTokens(ctx context.Context, jwt string) (string, string, error)
}

func MwLogger(log *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, _ := uuid.NewV7()

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

func MwAuth(log *zap.Logger, auth Auth, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		authorization := r.Header.Get("Authorization")

		access, ok := ejectJWT(authorization)
		if err := auth.VerifyJWT(ctx, access); ok && err == nil {
			next.ServeHTTP(w, r)
			return
		}

		refresh := r.Header.Get("Refresh-Token")
		if err := auth.VerifyJWT(ctx, refresh); err != nil {
			http.Error(w, "Unauthorized: No valid JWT was provided", http.StatusUnauthorized)
			return
		}

		newAccess, newRefresh, err := auth.RefreshTokens(ctx, refresh)
		if err != nil {
			log.Error("failed to refresh access token", zap.Error(err))

			switch {
			case errors.Is(err, ErrTokenWasRevoked):
				http.Error(w, "token was revoked", http.StatusUnauthorized)
			default:
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", newAccess))

		w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", newAccess))
		w.Header().Set("Refresh-Token", newRefresh)

		next.ServeHTTP(w, r)
	})
}

func ejectJWT(authHeader string) (string, bool) {
	const prefix = "Bearer "

	if !strings.HasPrefix(authHeader, prefix) {
		return "", false
	}

	return strings.TrimSpace(authHeader[len(prefix):]), true
}
