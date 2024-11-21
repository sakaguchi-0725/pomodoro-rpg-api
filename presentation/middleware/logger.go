package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"pomodoro-rpg-api/pkg/contextkey"
	"pomodoro-rpg-api/pkg/logger"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		requestID := getRequstID(r.Context())
		userID := getUserID(r.Context())

		logger.Logger.Info("request recived",
			slog.String("requestID", requestID),
			slog.String("userID", userID),
			slog.String("method", r.Method),
			slog.String("host", r.Host),
			slog.String("uri", r.RequestURI),
			slog.String("remote", r.RemoteAddr),
		)
	})
}

func getRequstID(ctx context.Context) string {
	if reqID, ok := ctx.Value(contextkey.RequestID).(string); ok {
		return reqID
	}
	return ""
}

func getUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(contextkey.UserID).(string); ok {
		return userID
	}
	return ""
}
