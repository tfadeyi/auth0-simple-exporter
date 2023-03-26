package exporter

import (
	"context"
	"time"

	"github.com/auth0-simple-exporter/internal/auth0"
	"github.com/auth0-simple-exporter/internal/exporter/metrics"
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type (
	exporter struct {
		// webserver
		metricsAddr      string
		hostPort         int
		profilingEnabled bool
		profilingPort    int

		// exporter
		namespace string
		subsystem string
		// checkpoint from where to start fetching logs
		startTime time.Time

		// webserver TLS
		tlsDisabled bool
		autoTLS     bool
		certFile    string
		keyFile     string
		tlsHosts    []string

		ctx context.Context

		// auth0
		client auth0.Fetcher

		// probe
		probeAddr                 string
		probePort                 int
		totalScrapes              prometheus.Counter
		targetScrapeRequestErrors prometheus.Counter
	}
	Option func(e *exporter)
)

// New returns an instance of the exporter
func New(opts ...Option) (*exporter, error) {
	e := &exporter{
		namespace: "",
		subsystem: "auth0",
		startTime: time.Now(),
		targetScrapeRequestErrors: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: "",
				Name:      "target_scrape_request_errors_total",
				Help:      "Errors in requests to the exporter",
			}),
		totalScrapes: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: "",
				Name:      "target_scrape_request_total",
				Help:      "Total requests to the exporter",
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

// metrics godoc
// @Summary     exporter's collected auth0 metrics.
// @Description Exposes the Auth0 metrics collected by the exporter.
// @Produce     json
// @Produce     text/plain; charset=utf-8
// @Router      /metrics [get]
// @tags        metrics, auth0
func (e *exporter) metrics() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		log.Debug("handling request for the auth0 tenant metrics")
		metrics := ctx.Get(metrics.ListCtxKey).(*metrics.Metrics)
		registry := prometheus.NewRegistry()
		registry.MustRegister(metrics.List()...)

		e.totalScrapes.Inc()
		log.Debug("handling request for the auth0 tenant metrics")
		err := collect(ctx.Request().Context(), metrics, e.client, e.startTime)
		if err != nil {
			// TODO check if the error is a rate limit error
			log.Errorf("error collecting event logs metrics from the selected Auth0 tenant: %s", err)
			e.targetScrapeRequestErrors.Inc()
		}

		log.Debug("successfully collected metrics from the auth0 tenant")
		promhttp.HandlerFor(registry, promhttp.HandlerOpts{}).ServeHTTP(ctx.Response(), ctx.Request())
		return nil
	}
}

// probe godoc
// @Summary     exporter's own metrics.
// @Description Exposes the exporter's own metrics, i.e: target_scrape_request_total.
// @Produce     json
// @Produce     text/plain; charset=utf-8
// @Router      /probe [get]
// @tags        metrics, prometheus
func (e *exporter) probe() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		promhttp.Handler().ServeHTTP(ctx.Response(), ctx.Request())
		return nil
	}
}

// collect collects all logs from Auth0 using startTime as the initial checkpoint
func collect(ctx context.Context, m *metrics.Metrics, client auth0.Fetcher, from time.Time) error {
	list, err := client.FetchAll(ctx, from)
	if err != nil {
		return errors.Annotate(err, "error fetching the log events from Auth0")
	}

	tenantLogEvents, ok := list.([]*management.Log)
	if !ok {
		return errors.New("log management mismatch")
	}

	for _, event := range tenantLogEvents {
		if err := m.Update(event); err != nil {
			log.Error(err)
			continue
		}
	}
	return nil
}
