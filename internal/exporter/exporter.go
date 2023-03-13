package exporter

import (
	"context"
	"github.com/auth0-simple-exporter/internal/auth0"
	"github.com/auth0-simple-exporter/internal/exporter/operations"
	"github.com/auth0-simple-exporter/internal/logging"
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type (
	exporter struct {
		metricsAddr      string
		hostPort         int
		ctx              context.Context
		profilingEnabled bool
		client           auth0.Fetcher

		totalScrapes              prometheus.Counter
		targetScrapeRequestErrors prometheus.Counter
	}
	metricFunc func(ctx context.Context, log *management.Log) error
	Option     func(e *exporter)
)

const (
	registryKey = operations.CtxKey("registry")
	namespace   = "auth0"
	subsystem   = ""
)

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

func New(opts ...Option) (*exporter, error) {
	e := &exporter{
		targetScrapeRequestErrors: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "target_scrape_request_errors_total",
				Help:      "Errors in requests to the exporter",
			}),
		totalScrapes: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "target_scrape_request_total",
				Help:      "Errors in requests to the exporter",
			}),
	}
	for _, opt := range opts {
		// apply options
		opt(e)
	}

	prometheus.MustRegister(e.targetScrapeRequestErrors)
	prometheus.MustRegister(e.totalScrapes)

	return e, nil
}

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Bearer realm="Echo"`)
	w.Header().Set("Content-Type", "application/json")
}

// metrics godoc
// @Summary     exporter metrics metrics.
// @Description Exposes Auth0 metrics.
// @Produce     json
// @Produce     text/plain; charset=utf-8
// @Router      /metrics [get]
// @tags        metrics, prometheus
func (e *exporter) metrics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setHeaders(w)
		ctx := r.Context()
		log := logging.LoggerFromContext(ctx)

		e.totalScrapes.Inc()
		if _, ok := ctx.Value(registryKey).(*prometheus.Registry); !ok {
			log.Error(errors.New("prometheus registry is missing from the request context"), "Error in the exporter metricsMiddleware")
			http.Error(w, "error in the exporter metricsMiddleware, missing prometheus registry", http.StatusInternalServerError)
			return
		}

		err := collect(ctx, e.client,
			operations.LoginOperationsEventHandler,
			operations.ApiOperationsEventHandler,
			operations.LogoutOperationsEventHandler)
		if err != nil {
			log.Error(err, "Error collecting event logs metrics from the selected Auth0 tenant")
			e.targetScrapeRequestErrors.Inc()
		}

		promhttp.HandlerFor(registryFromContext(ctx), promhttp.HandlerOpts{}).ServeHTTP(w, r)
	}
}

func (e *exporter) probe() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		promhttp.Handler().ServeHTTP(writer, request)
	}
}

func collect(ctx context.Context, client auth0.Fetcher, metricFss ...metricFunc) error {
	list, err := client.FetchAll(ctx)
	if err != nil {
		return errors.Annotate(err, "error fetching the log events from Auth0")
	}

	tenantLogEvents, ok := list.([]*management.Log)
	if !ok {
		return errors.New("log management mismatch")
	}
	for _, event := range tenantLogEvents {
		for _, fs := range metricFss {
			if err := fs(ctx, event); err != nil {
				continue
			}
			// success in updating the metrics, goto next event/log
			break
		}
	}
	return nil
}
