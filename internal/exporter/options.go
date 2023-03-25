package exporter

import (
	"context"

	"github.com/auth0-simple-exporter/internal/auth0"
)

func TLSHost(host string) Option {
	return func(e *exporter) {
		e.tlsHost = host
	}
}
func Context(ctx context.Context) Option {
	return func(e *exporter) {
		e.ctx = ctx
	}
}

func Client(client auth0.Fetcher) Option {
	return func(e *exporter) {
		e.client = client
	}
}

func Profiling(p bool) Option {
	return func(e *exporter) {
		e.profilingEnabled = p
	}
}

func MetricsAddr(addr string) Option {
	return func(e *exporter) {
		e.metricsAddr = addr
	}
}

func Port(port int) Option {
	return func(e *exporter) {
		e.hostPort = port
	}
}

func Namespace(namespace string) Option {
	return func(e *exporter) {
		e.namespace = namespace
	}
}

func ManagedTLS(t bool) Option {
	return func(e *exporter) {
		e.managedTLS = t
	}
}

func CertFile(filename string) Option {
	return func(e *exporter) {
		e.certFile = filename
	}
}

func KeyFile(filename string) Option {
	return func(e *exporter) {
		e.keyFile = filename
	}
}

func DisableTLS(t bool) Option {
	return func(e *exporter) {
		e.tlsDisabled = t
	}
}
