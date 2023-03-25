package exporter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/auth0-simple-exporter/internal/exporter/metrics"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/sync/errgroup"
	_ "net/http/pprof"
)

// Export configures the exporter Router and starts the server with the given configuration
func (e *exporter) Export() error {
	//log := logging.LoggerFromContext(e.ctx)
	server := echo.New()

	server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return metrics.NamespaceMiddleware(next, e.namespace)
	})
	server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return metrics.SubsystemMiddleware(next, e.subsystem)
	})
	server.Use(metrics.Middleware)
	server.Use(middleware.Recover())
	server.Use(middleware.Logger())
	server.HideBanner = true
	server.HidePort = true

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

	log.Infof("starting exporter",
		"port", e.hostPort,
		"metrics-address", "/"+e.metricsAddr,
	)

	// N.B: we want to fail hard if users don't specify the autoTLS and cert/key pair isn't there.
	// and warn users to use the tls.managed flag if they don't want to manage their certs.
	grp, ctx := errgroup.WithContext(e.ctx)
	switch {
	case !e.tlsDisabled && !e.autoTLS:
		log.Info("exporter's TLS option was enabled, using certificates in the local machine")
		// start server using the given certs
		grp.Go(func() error {
			return server.StartTLS(fmt.Sprintf(":%d", e.hostPort), e.certFile, e.keyFile)
		})
	case !e.tlsDisabled && e.autoTLS:
		log.Info("exporter's managed TLS option was enabled")
		// start server using managed certs
		// Cache certificates to avoid issues with rate limits (https://letsencrypt.org/docs/rate-limits)
		server.AutoTLSManager.HostPolicy = autocert.HostWhitelist(e.tlsHosts...)
		server.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
		grp.Go(func() error {
			return server.StartAutoTLS(fmt.Sprintf(":%d", e.hostPort))
		})
	default:
		log.Info("Starting exporter with an insecure connection!")
		grp.Go(func() error {
			return server.Start(fmt.Sprintf(":%d", e.hostPort))
		})
	}

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

	grp.Go(func() error {
		<-ctx.Done()
		return server.Shutdown(context.Background())
	})

	return grp.Wait()
}
