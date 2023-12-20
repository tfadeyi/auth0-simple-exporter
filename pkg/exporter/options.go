package exporter

import (
	"time"

	"github.com/tfadeyi/auth0-simple-exporter/pkg/client"
	"github.com/tfadeyi/auth0-simple-exporter/pkg/logging"
)

func Client(client client.Client) Option {
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

func RequestTimeout(timeout time.Duration) Option {
	return func(e *exporter) {
		e.requestTimeout = timeout
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

func AutoTLS(t bool) Option {
	return func(e *exporter) {
		e.autoTLS = t
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

func TLSHosts(hosts []string) Option {
	return func(e *exporter) {
		e.tlsHosts = hosts
	}
}

func ProbePort(port int) Option {
	return func(e *exporter) {
		e.probePort = port
	}
}

func ProbeAddr(addr string) Option {
	return func(e *exporter) {
		e.probeAddr = addr
	}
}

func ProfilingPort(port int) Option {
	return func(e *exporter) {
		e.profilingPort = port
	}
}

func From(time time.Time) Option {
	return func(e *exporter) {
		e.startTime = time
	}
}

func Subsystem(sub string) Option {
	return func(e *exporter) {
		e.subsystem = sub
	}
}

// Logger for the exporter
func Logger(l logging.Logger) Option {
	return func(e *exporter) {
		e.logger = l
	}
}

func DisableUserMetrics(flag bool) Option {
	return func(e *exporter) {
		e.userMetricDisabled = flag
	}
}
