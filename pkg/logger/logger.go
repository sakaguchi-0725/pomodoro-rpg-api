package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"pomodoro-rpg-api/pkg/contextkey"
)

var Logger *slog.Logger

type LogLevel string

const (
	INFO  LogLevel = "INFO"
	ERROR LogLevel = "ERROR"
	WARN  LogLevel = "WARN"
	DEBUG LogLevel = "DEBUG"
)

func Init() {
	if Logger == nil {
		Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	}
}

func Event(ctx context.Context, level LogLevel, msg string, err error) {
	userID, ok := ctx.Value(contextkey.UserID).(string)
	if !ok {
		userID = ""
	}

	requestID, ok := ctx.Value(contextkey.RequestID).(string)
	if !ok {
		requestID = ""
	}

	switch level {
	case INFO:
		Logger.Info(msg,
			slog.String("requestID", requestID),
			slog.String("userID", userID),
			slog.String("Stack Trace", fmt.Sprintf("%+v", err)),
		)
	case ERROR:
		Logger.Error(msg,
			slog.String("requestID", requestID),
			slog.String("userID", userID),
			slog.String("Stack Trace", fmt.Sprintf("%+v", err)),
		)
	case WARN:
		Logger.Warn(msg,
			slog.String("requestID", requestID),
			slog.String("userID", userID),
			slog.String("Stack Trace", fmt.Sprintf("%+v", err)),
		)
	case DEBUG:
		Logger.Debug(msg,
			slog.String("requestID", requestID),
			slog.String("userID", userID),
			slog.String("Stack Trace", fmt.Sprintf("%+v", err)),
		)
	}
}
