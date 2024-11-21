package middleware

import (
	"context"
	"net/http"
	"pomodoro-rpg-api/pkg/contextkey"

	"github.com/google/uuid"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()

		ctx := context.WithValue(r.Context(), contextkey.RequestID, requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
