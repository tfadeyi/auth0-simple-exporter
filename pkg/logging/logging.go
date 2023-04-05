package logging

import (
	"context"

	kitlog "github.com/go-kit/log"
	"github.com/go-logr/logr"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/common/promlog"
	"github.com/tonglil/gokitlogr"
)

type (
	Logger struct {
		logr.Logger
	}
)

var (
	// contextKey is how we find Loggers in a context.Context.
	contextKey string = "prom-logger"
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

// NewPromLogger wraps the creation of a new prometheus production logger
func NewPromLogger() Logger {
	l := kitlog.NewSyncLogger(promlog.New(&promlog.Config{}))
	return Logger{gokitlogr.New(&l)}
}

func NewPromLoggerWithOpts(lvl string) Logger {
	level := &promlog.AllowedLevel{}
	if level.Set(lvl) != nil {
		return NewPromLogger()
	}
	l := kitlog.NewSyncLogger(promlog.New(&promlog.Config{Level: level}))
	return Logger{gokitlogr.New(&l)}
}

// Middleware creates a production logger and sets it in the context for the incoming requests.
// incoming requests will have a logger that can be use in the handler.
func Middleware(next echo.HandlerFunc, loggers ...Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := NewPromLogger()
		if len(loggers) > 0 {
			log = loggers[0]
		}
		ctx := EchoContextWithLogger(c, log)

		req := ctx.Request()
		res := ctx.Response()
		if err := next(ctx); err != nil {
			ctx.Error(err)
		}

		id := req.Header.Get(echo.HeaderXRequestID)
		if id == "" {
			id = res.Header().Get(echo.HeaderXRequestID)
		}
		p := req.URL.Path
		if p == "" {
			p = "/"
		}

		log.WithValues(
			"requestID", id,
			"remote_ip", ctx.RealIP(),
			"host", req.Host,
			"uri", req.RequestURI,
			"method", req.Method,
			"path", p,
			"route", c.Path(),
			"user_agent", req.UserAgent(),
			"status", res.Status,
		).Info("incoming request")
		return nil
	}
}

// EchoContextWithLogger the logger in the given echo context
func EchoContextWithLogger(ctx echo.Context, l Logger) echo.Context {
	ctx.Set(contextKey, l)
	return ctx
}

// LoggerFromEchoContext find the logger in the echo context or returns a new prom logger
func LoggerFromEchoContext(ctx echo.Context) Logger {
	l, ok := ctx.Get(contextKey).(Logger)
	if !ok {
		return NewPromLogger()
	}
	return l
}
