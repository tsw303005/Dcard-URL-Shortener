package logkit

import (
	"context"
	"log"
)

type loggerContextKey int8

const contextKeyLogger loggerContextKey = iota

func WithContext(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, contextKeyLogger, logger)
}

func FromContext(ctx context.Context) *Logger {
	logger, ok := ctx.Value(contextKeyLogger).(*Logger)

	if !ok {
		log.Fatal("logger is not found in this context")
	}

	return logger
}
