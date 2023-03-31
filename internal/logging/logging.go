package logging

import (
	"context"
	"net/http"

	kitlog "github.com/go-kit/log"
	"github.com/go-logr/logr"
	echolog "github.com/labstack/gommon/log"
	"github.com/prometheus/common/promlog"
	"github.com/tonglil/gokitlogr"
)

type (
	Logger struct {
		logr.Logger
	}
)

// ContextWithLogger wraps the logr NewContext function
func ContextWithLogger(ctx context.Context, l Logger) context.Context {
	return logr.NewContext(ctx, l.Logger)
}

// LoggerFromContext wraps the LoggerFromContext or creates a Zap production logger
func LoggerFromContext(ctx context.Context) Logger {
	l, err := logr.FromContext(ctx)
	if err != nil {
		return NewPromLogger()
	}
	return Logger{l}
}

// ContextWithEchoLogger wraps the logr NewContext function
func ContextWithEchoLogger(ctx context.Context, l Logger) context.Context {
	return logr.NewContext(ctx, l.Logger)
}

func EchoLoggerFromContext(ctx context.Context) Logger {
	l, err := logr.FromContext(ctx)
	if err != nil {
		return NewPromLogger()
	}
	return Logger{l}
}

// NewEchoLogger wraps the creation of a new gommon production logger
func NewEchoLogger() *echolog.Logger {
	return echolog.New("auth0-simple-exporter")
}

// NewPromLogger wraps the creation of a new prometheus production logger
func NewPromLogger() Logger {
	l := kitlog.NewSyncLogger(promlog.New(&promlog.Config{}))
	return Logger{gokitlogr.New(&l)}
}

// Middleware creates a production logger and sets it in the context for the incoming requests.
// incoming requests will have a logger that can be use in the handler.
// example: log := logging.LoggerFromContext(r.Context())
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := NewPromLogger()
		ctx := ContextWithLogger(r.Context(), log)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
