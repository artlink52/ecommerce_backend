package logger

import (
	"context"
	"log/slog"
	"os"
)

type LoggerContextKey struct{}

var (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"

	key = LoggerContextKey{}
)

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func ToContext(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, key, log)
}

func FromContext(ctx context.Context) *slog.Logger {
	log, ok := ctx.Value(key).(*slog.Logger)
	if !ok {
		panic("logger not found in context")
	}
	return log
}
