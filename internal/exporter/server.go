package exporter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/auth0-simple-exporter/internal/exporter/metrics"
	"github.com/auth0-simple-exporter/internal/logging"
	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"
	_ "net/http/pprof"
)

// Export configures the exporter Router and starts the server with the given configuration
func (e *exporter) Export() error {
	log := logging.LoggerFromContext(e.ctx)
	server := echo.New()

	server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return metrics.NamespaceMiddleware(next, e.namespace)
	})
	server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return metrics.SubsystemMiddleware(next, e.subsystem)
	})
	server.Use(metrics.MetricsMiddleware)
	//server.Use(middleware.Recover())

	server.GET("/", func(ctx echo.Context) error {
		return ctx.HTML(http.StatusOK, `<html>
			<head><title>Auth0 Exporter</title></head>
			<body>
			<h1>Auth0 Exporter</h1>
			<p><a href="`+e.metricsAddr+`">Metrics</a></p>
			</body>
			</html>`)
	})

	// the exporter is using the log mgmt api
	// Metrics path route
	server.GET(fmt.Sprintf("/%s", e.metricsAddr), e.metrics())

	log.Info("starting exporter",
		"port", e.hostPort,
		"metrics-address", "/"+e.metricsAddr,
	)

	grp, ctx := errgroup.WithContext(e.ctx)
	grp.Go(func() error {
		return server.Start(fmt.Sprintf(":%d", e.hostPort))
	})
	grp.Go(func() error {
		<-ctx.Done()
		return server.Shutdown(context.Background())
	})

	// if profiling is enabled start the pprof server
	if e.profilingEnabled {
		log.Info(
			"pprof profiling has been activated on port :6060",
		)
		profilingServer := &http.Server{
			Addr: fmt.Sprintf(":%d", 6060),
		}
		grp.Go(func() error {
			return profilingServer.ListenAndServe()
		})
		grp.Go(func() error {
			<-ctx.Done()
			return profilingServer.Shutdown(context.Background())
		})
	}

	return grp.Wait()
}

//func (e *exporter) Export() error {
//	log := logging.LoggerFromContext(e.ctx)
//	// configure and start sever
//	router := chi.NewRouter()
//	router.Use(logging.Middleware)
//	router.Use(chiMiddleware.Heartbeat("/ping"))
//
//	// home path route. It will always return a static html with guide on how to use the exporter
//	router.With(chiMiddleware.AllowContentType("text/html")).Route("/", func(r chi.Router) {
//		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
//			log := logging.LoggerFromContext(r.Context())
//			if _, err := w.Write([]byte(`<html>
//			<head><title>Auth0 Exporter</title></head>
//			<body>
//			<h1>Auth0 Exporter</h1>
//			<p><a href="` + e.metricsAddr + `">Metrics</a></p>
//			</body>
//			</html>`)); err != nil {
//				log.Error(err, "Error generating the index page")
//				http.Error(w, "Missing index page", http.StatusInternalServerError)
//				return
//			}
//		})
//	})
//
//	// the exporter is using the log mgmt api
//	// Metrics path route
//	router.With(metrics.Middleware).Route(fmt.Sprintf("/%s", e.metricsAddr), func(r chi.Router) {
//		r.Get("/", e.metrics())
//	})
//	router.Get("/probe", e.probe())
//
//	log.Info("starting exporter",
//		"port", e.hostPort,
//		"metrics-address", e.metricsAddr,
//		"probe-address", "probe")
//
//	// start server with the router with setup before
//	server := &http.Server{
//		Addr:    fmt.Sprintf(":%d", e.hostPort),
//		Handler: router,
//	}
//
//	grp, ctx := errgroup.WithContext(e.ctx)
//	grp.Go(func() error {
//		return server.ListenAndServe()
//	})
//	grp.Go(func() error {
//		<-ctx.Done()
//		return server.Shutdown(context.Background())
//	})
//
//	// if profiling is enabled start the pprof server
//	if e.profilingEnabled {
//		log.Info(
//			"pprof profiling has been activated on port :6060",
//		)
//		profilingServer := &http.Server{
//			Addr: fmt.Sprintf(":%d", 6060),
//		}
//		grp.Go(func() error {
//			return profilingServer.ListenAndServe()
//		})
//		grp.Go(func() error {
//			<-ctx.Done()
//			return profilingServer.Shutdown(context.Background())
//		})
//	}
//
//	return grp.Wait()
//}
