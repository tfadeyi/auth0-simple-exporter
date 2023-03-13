package exporter

import (
	"context"
	"github.com/auth0-simple-exporter/internal/exporter/operations"
	"github.com/prometheus/client_golang/prometheus"
)

func registryFromContext(ctx context.Context) *prometheus.Registry {
	registry, ok := ctx.Value(registryKey).(*prometheus.Registry)
	if !ok {
		return prometheus.NewRegistry()
	}
	return registry
}

func contextWithRegistry(ctx context.Context, registry *prometheus.Registry) context.Context {
	return context.WithValue(ctx, registryKey, registry)
}

func contextWithMetrics(ctx context.Context, mss map[operations.CtxKey]prometheus.Collector) context.Context {
	registry := registryFromContext(ctx)

	for name, metric := range mss {
		ctx = context.WithValue(ctx, name, metric)
		registry.MustRegister(metric)
	}
	return ctx
}
