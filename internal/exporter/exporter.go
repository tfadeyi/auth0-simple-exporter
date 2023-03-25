package exporter

import (
	"context"
	"net/http"

	"github.com/auth0-simple-exporter/internal/auth0"
	"github.com/auth0-simple-exporter/internal/exporter/metrics"
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type (
	exporter struct {
		metricsAddr      string
		hostPort         int
		profilingEnabled bool

		// TLS
		tlsDisabled bool
		managedTLS  bool
		certFile    string
		keyFile     string
		tlsHost     string

		ctx context.Context

		client auth0.Fetcher

		totalScrapes              prometheus.Counter
		targetScrapeRequestErrors prometheus.Counter

		namespace string
		subsystem string
	}
	Option func(e *exporter)
)

func New(opts ...Option) (*exporter, error) {
	e := &exporter{
		namespace: "",
		subsystem: "auth0",
		//targetScrapeRequestErrors: prometheus.NewCounter(
		//	prometheus.CounterOpts{
		//		Namespace: namespace,
		//		Name:      "target_scrape_request_errors_total",
		//		Help:      "Errors in requests to the exporter",
		//	}),
		//totalScrapes: prometheus.NewCounter(
		//	prometheus.CounterOpts{
		//		Namespace: namespace,
		//		Name:      "target_scrape_request_total",
		//		Help:      "Errors in requests to the exporter",
		//	}),
	}
	for _, opt := range opts {
		// apply options
		opt(e)
	}

	//prometheus.MustRegister(e.targetScrapeRequestErrors)
	//prometheus.MustRegister(e.totalScrapes)

	return e, nil
}

// metrics godoc
// @Summary     exporter metrics.
// @Description Exposes Auth0 metrics.
// @Produce     json
// @Produce     text/plain; charset=utf-8
// @Router      /metrics [get]
// @tags        metrics, prometheus
//func (e *exporter) metrics() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		setHeaders(w)
//		ctx := r.Context()
//		log := logging.LoggerFromContext(ctx)
//
//		e.totalScrapes.Inc()
//		if _, ok := ctx.Value(metrics.RegistryKey).(*prometheus.Registry); !ok {
//			log.Error(errors.New("prometheus registry is missing from the request context"), "Error in the exporter metricsMiddleware")
//			http.Error(w, "error in the exporter metricsMiddleware, missing prometheus registry", http.StatusInternalServerError)
//			return
//		}
//
//		err := collect(ctx, e.client, metrics.Handlers...)
//		if err != nil {
//			log.Error(err, "Error collecting event logs metrics from the selected Auth0 tenant")
//			e.targetScrapeRequestErrors.Inc()
//		}
//
//		promhttp.HandlerFor(metrics.RegistryFromContext(ctx), promhttp.HandlerOpts{}).ServeHTTP(w, r)
//	}
//}

func (e *exporter) metrics() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		log := ctx.Logger()
		metrics := ctx.Get(metrics.ListCtxKey).(*metrics.Metrics)
		registry := prometheus.NewRegistry()
		registry.MustRegister(metrics.List()...)

		err := collect(metrics, ctx.Request().Context(), e.client)
		if err != nil {
			log.Error(err, "Error collecting event logs metrics from the selected Auth0 tenant")
		}

		promhttp.HandlerFor(registry, promhttp.HandlerOpts{}).ServeHTTP(ctx.Response(), ctx.Request())
		return nil
	}
}

func (e *exporter) probe() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		promhttp.Handler().ServeHTTP(writer, request)
	}
}

//func collect(ctx context.Context, client auth0.Fetcher, metricFss ...metrics.MetricFunc) error {
//	list, err := client.FetchAll(ctx)
//	if err != nil {
//		return errors.Annotate(err, "error fetching the log events from Auth0")
//	}
//
//	tenantLogEvents, ok := list.([]*management.Log)
//	if !ok {
//		return errors.New("log management mismatch")
//	}
//	for _, event := range tenantLogEvents {
//		for _, fs := range metricFss {
//			if err := fs(ctx, event); err != nil {
//				continue
//			}
//			// success in updating the metrics, goto next event/log
//			break
//		}
//	}
//	return nil
//}

func collect(m *metrics.Metrics, ctx context.Context, client auth0.Fetcher) error {
	list, err := client.FetchAll(ctx)
	if err != nil {
		return errors.Annotate(err, "error fetching the log events from Auth0")
	}

	tenantLogEvents, ok := list.([]*management.Log)
	if !ok {
		return errors.New("log management mismatch")
	}

	for _, event := range tenantLogEvents {
		m.Update(event)
	}
	return nil
}
