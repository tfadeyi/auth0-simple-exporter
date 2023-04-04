package exporter

import (
	"context"
	"net/http"
	"time"

	"github.com/auth0-simple-exporter/pkg/client"
	"github.com/auth0-simple-exporter/pkg/client/logs"
	"github.com/auth0-simple-exporter/pkg/exporter/metrics"
	"github.com/auth0-simple-exporter/pkg/logging"
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

		ctx    context.Context
		logger logging.Logger

		// auth0
		client client.Client

		// probe
		probeAddr                 string
		probePort                 int
		totalScrapes              prometheus.Counter
		targetScrapeRequestErrors prometheus.Counter
		probeRegistry             *prometheus.Registry
	}
	Option func(e *exporter)
)

// New returns an instance of the exporter
func New(ctx context.Context, opts ...Option) *exporter {
	e := &exporter{
		namespace: "auth0",
		subsystem: "",
		ctx:       ctx,
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
		probeRegistry: prometheus.NewRegistry(),
	}
	for _, opt := range opts {
		// apply options
		opt(e)
	}

	if ctx == nil {
		e.ctx = context.Background()
	}

	e.probeRegistry.MustRegister(e.targetScrapeRequestErrors, e.totalScrapes)

	return e
}

// metrics godoc
//
//	@Summary				Auth0 metrics in Prometheus format.
//	@Description.markdown	metrics.md
//	@Produce				json
//	@Produce				text/plain; charset=utf-8
//	@Router					/metrics [get]
func (e *exporter) metrics() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		log := logging.LoggerFromEchoContext(ctx)
		log.Info("handling request for the auth0 tenant metrics")
		metrics := ctx.Get(metrics.ListCtxKey).(*metrics.Metrics)
		registry := prometheus.NewRegistry()
		registry.MustRegister(metrics.List()...)

		e.totalScrapes.Inc()
		log.Info("handling request for the auth0 tenant metrics")
		err := e.collect(ctx.Request().Context(), metrics)
		switch {
		case errors.Is(err, logs.ErrAPIRateLimitReached):
			log.Error(err, "reached the Auth0 rate limit, fetching should resume shortly")
			e.targetScrapeRequestErrors.Inc()
		case err != nil:
			log.Error(err, "error collecting event logs metrics from the selected Auth0 tenant")
			e.targetScrapeRequestErrors.Inc()
			return echo.NewHTTPError(http.StatusInternalServerError, "Exporter encountered some issues when collecting logs from Auth0")
		}

		log.Info("successfully collected metrics from the auth0 tenant")
		promhttp.HandlerFor(registry, promhttp.HandlerOpts{}).ServeHTTP(ctx.Response(), ctx.Request())
		return nil
	}
}

// probe godoc
//
//	@Summary		Exporter's own metrics for its operations.
//	@Description	Exposes the exporter's own metrics, i.e: target_scrape_request_total.
//	@Produce		json
//	@Produce		text/plain; charset=utf-8
//	@Router			/probe [get]
func (e *exporter) probe() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		promhttp.HandlerFor(e.probeRegistry, promhttp.HandlerOpts{}).ServeHTTP(ctx.Response(), ctx.Request())
		return nil
	}
}

// collect collects all logs from Auth0 using startTime as the initial checkpoint
func (e *exporter) collect(ctx context.Context, m *metrics.Metrics) error {
	list, err := e.client.List(ctx, e.startTime)
	if err != nil {
		return errors.Annotate(err, "error fetching the log events from Auth0")
	}

	tenantLogEvents, ok := list.([]*management.Log)
	if !ok {
		return errors.New("client FetchAll didn't return the expect list types, expected Log type")
	}

	for _, event := range tenantLogEvents {
		if err := m.Update(event); err != nil {
			log.Error(err)
			continue
		}
	}
	return nil
}
