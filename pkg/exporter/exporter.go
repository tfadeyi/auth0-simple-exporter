package exporter

import (
	"context"
	"net/http"
	"time"

	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tfadeyi/auth0-simple-exporter/pkg/client"
	"github.com/tfadeyi/auth0-simple-exporter/pkg/client/logs"
	"github.com/tfadeyi/auth0-simple-exporter/pkg/exporter/metrics"
	"github.com/tfadeyi/auth0-simple-exporter/pkg/logging"
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
		metricsObject             *metrics.Metrics
		metricsRegistry           *prometheus.Registry
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
		probeRegistry:   prometheus.NewRegistry(),
		metricsRegistry: prometheus.NewRegistry(),
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

		e.totalScrapes.Inc()
		log.Info("handling request for the auth0 tenant metrics")
		err := e.collect(ctx.Request().Context(), e.metricsObject)
		switch {
		case errors.Is(err, logs.ErrAPIRateLimitReached):
			log.Error(err, "reached the Auth0 rate limit, fetching should resume shortly")
			e.targetScrapeRequestErrors.Inc()
		case err != nil:
			log.Error(err, "error collecting event metrics from the selected Auth0 tenant")
			e.targetScrapeRequestErrors.Inc()
			return echo.NewHTTPError(http.StatusInternalServerError, "Exporter encountered some issues when collecting events from Auth0, please check the exporter logs with --log.level debug")
		}

		log.Info("successfully collected metrics from the Auth0 tenant")
		promhttp.HandlerFor(e.metricsRegistry, promhttp.HandlerOpts{}).ServeHTTP(ctx.Response(), ctx.Request())
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
//
//nolint:goconst
func (e *exporter) collect(ctx context.Context, m *metrics.Metrics) error {
	// Process logs
	list, err := e.client.Log.List(ctx, e.startTime)
	switch {
	case errors.Is(err, context.Canceled):
		eventLogs := list.([]*management.Log)
		e.logger.V(0).Error(err, "Request was terminated by the client,"+
			"the exporter could not finish polling the Auth0 log client to fetch the tenant logs."+
			"Please increase the client timeout or try adding the --auth0.from flag", "logs_events_found", len(eventLogs), "from", e.startTime)
	case err != nil:
		return errors.Annotate(err, "error fetching the log events from Auth0")
	}

	tenantLogEvents, ok := list.([]*management.Log)
	if !ok {
		return errors.New("Auth0 log client did not return the expected list of Log type")
	}

	for _, event := range tenantLogEvents {
		if err := m.Update(event); err != nil {
			e.logger.V(0).Error(err, err.Error())
			continue
		}
	}

	// Process users
	list, err = e.client.User.List(ctx)
	switch {
	case errors.Is(err, context.Canceled):
		eventUsers := list.([]*management.User)
		e.logger.V(0).Error(err, "Request was terminated by the client,"+
			"the exporter could not finish polling the Auth0 user client to fetch the tenant users."+
			"Please increase the client timeout", "users_found", len(eventUsers))
	case err != nil:
		return errors.Annotate(err, "error fetching the users from Auth0")
	}
	tenantUsers, ok := list.([]*management.User)
	if !ok {
		return errors.New("auth0 client users fetch didn't return the expected list of User type")
	}
	if err := m.ProcessUsers(tenantUsers); err != nil {
		e.logger.V(0).Error(err, err.Error())
	}

	return nil
}
