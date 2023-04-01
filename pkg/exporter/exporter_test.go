package exporter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/auth0-simple-exporter/pkg/client"
	"github.com/auth0-simple-exporter/pkg/client/logs"
	"github.com/auth0-simple-exporter/pkg/exporter/metrics"
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExporter(t *testing.T) {
	t.Parallel()

	t.Run("exporter options configuration", func(t *testing.T) {
		current := time.Now()

		exp := &exporter{
			metricsAddr:      "metrics",
			hostPort:         8080,
			profilingEnabled: true,
			profilingPort:    6060,
			namespace:        "auth0",
			subsystem:        "my_app",
			startTime:        current,
			tlsDisabled:      false,
			autoTLS:          false,
			certFile:         "cert.crt",
			keyFile:          "tls.key",
			probeAddr:        "probe",
			probePort:        8081,
		}

		act := New(
			context.Background(),
			Profiling(true),
			ProfilingPort(6060),

			ProbePort(8081),
			ProbeAddr("probe"),

			KeyFile("tls.key"),
			CertFile("cert.crt"),

			Namespace("auth0"),
			Subsystem("my_app"),

			MetricsAddr("metrics"),
			Port(8080),

			AutoTLS(false),
			DisableTLS(false),

			From(current),
		)

		assert.Equal(t, exp.profilingPort, act.profilingPort)
		assert.Equal(t, exp.profilingEnabled, act.profilingEnabled)

		assert.Equal(t, exp.subsystem, act.subsystem)
		assert.Equal(t, exp.namespace, act.namespace)

		assert.Equal(t, exp.probeAddr, act.probeAddr)
		assert.Equal(t, exp.probePort, act.probePort)

		assert.Equal(t, exp.metricsAddr, act.metricsAddr)
		assert.Equal(t, exp.hostPort, act.hostPort)

		assert.Equal(t, exp.startTime, act.startTime)

		assert.Equal(t, exp.autoTLS, act.autoTLS)
		assert.Equal(t, exp.tlsDisabled, act.tlsDisabled)

		assert.Equal(t, exp.keyFile, act.keyFile)
		assert.Equal(t, exp.certFile, act.certFile)
	})
	t.Run("fail exporter collect if error occurs in auth0 client", func(t *testing.T) {
		ctx := context.Background()
		client, err := client.NewWithOpts(client.Options{Client: &logs.ClientMock{ListFunc: func(ctx context.Context, args ...interface{}) (interface{}, error) {
			return nil, errors.New("some error")
		}}})
		require.NoError(t, err)

		current := time.Now()
		e := exporter{
			startTime: current,
			ctx:       ctx,
			client:    client,
		}
		require.Error(t, e.collect(ctx, nil))
	})
	t.Run("successful execute exporter collect if auth0 client returns a empty log list", func(t *testing.T) {
		ctx := context.Background()
		client, err := client.NewWithOpts(client.Options{Client: &logs.ClientMock{ListFunc: func(ctx context.Context, args ...interface{}) (interface{}, error) {
			return []*management.Log{}, nil
		}}})
		require.NoError(t, err)
		current := time.Now()
		e := exporter{
			startTime: current,
			ctx:       ctx,
			client:    client,
		}
		require.NoError(t, e.collect(ctx, nil))
	})
	t.Run("fail exporter collect if auth0 client didn't return a list of logs", func(t *testing.T) {
		ctx := context.Background()
		client, err := client.NewWithOpts(client.Options{Client: &logs.ClientMock{ListFunc: func(ctx context.Context, args ...interface{}) (interface{}, error) {
			return []string{}, nil
		}}})
		require.NoError(t, err)
		current := time.Now()
		e := exporter{
			startTime: current,
			ctx:       ctx,
			client:    client,
		}
		require.Error(t, e.collect(ctx, nil))
	})
}

func TestExporterHandler(t *testing.T) {
	t.Parallel()

	t.Run("don't fail if API rate limit is reached", func(t *testing.T) {
		ctx := context.Background()
		client, err := client.NewWithOpts(client.Options{Client: &logs.ClientMock{ListFunc: func(ctx context.Context, args ...interface{}) (interface{}, error) {
			return []string{}, logs.ErrAPIRateLimitReached
		}}})
		require.NoError(t, err)
		current := time.Now()
		exporter := New(ctx, From(current), Client(client))

		metricsServer := echo.New()
		metricsServer.Use(metrics.Middleware)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		echoCtx := metricsServer.NewContext(req, rec)
		echoCtx.Set(metrics.ListCtxKey, metrics.New(exporter.namespace, exporter.subsystem))

		require.NoError(t, exporter.metrics()(echoCtx))
		assert.Equal(t, http.StatusOK, rec.Code)
	})
	t.Run("successful request if the auth0 client returns 0 items", func(t *testing.T) {
		ctx := context.Background()
		client, err := client.NewWithOpts(client.Options{Client: &logs.ClientMock{ListFunc: func(ctx context.Context, args ...interface{}) (interface{}, error) {
			return []*management.Log{}, nil
		}}})
		require.NoError(t, err)
		current := time.Now()
		exporter := New(ctx, From(current), Client(client))

		metricsServer := echo.New()
		metricsServer.Use(metrics.Middleware)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		echoCtx := metricsServer.NewContext(req, rec)
		echoCtx.Set(metrics.ListCtxKey, metrics.New(exporter.namespace, exporter.subsystem))

		require.NoError(t, exporter.metrics()(echoCtx))
		assert.Equal(t, http.StatusOK, rec.Code)
	})
	t.Run("fail if Auth0 client errors with an unexpected error", func(t *testing.T) {
		ctx := context.Background()
		client, err := client.NewWithOpts(client.Options{Client: &logs.ClientMock{ListFunc: func(ctx context.Context, args ...interface{}) (interface{}, error) {
			return []string{}, errors.New("unexpected error")
		}}})
		require.NoError(t, err)
		current := time.Now()
		exporter := New(ctx, From(current), Client(client))

		metricsServer := echo.New()
		metricsServer.Use(metrics.Middleware)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		echoCtx := metricsServer.NewContext(req, rec)
		echoCtx.Set(metrics.ListCtxKey, metrics.New(exporter.namespace, exporter.subsystem))

		require.Error(t, exporter.metrics()(echoCtx))
	})
}
