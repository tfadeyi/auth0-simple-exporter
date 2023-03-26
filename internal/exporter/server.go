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
	server := echo.New()

	server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return metrics.NamespaceMiddleware(next, e.namespace)
	})
	server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return metrics.SubsystemMiddleware(next, e.subsystem)
	})
	server.Use(metrics.Middleware)
	server.Use(middleware.Recover())
	server.HideBanner = true
	server.HidePort = true

	server.GET("/", func(ctx echo.Context) error {
		return ctx.HTML(http.StatusOK, fmt.Sprintf(`<html>
			<head><title>Auth0 Exporter</title></head>
			<body>
			<h1>Auth0 Exporter</h1>
			<p><a href="%s">Metrics</a></p>
			</body>
			</html>`, e.metricsAddr))
	})
	server.GET("/healthz", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, "ok")
	})

	// the exporter is using the log mgmt api
	// Metrics path route
	server.GET(fmt.Sprintf("/%s", e.metricsAddr), e.metrics())

	log.Infof("starting exporter",
		"port", e.hostPort,
		"metrics-address", "/"+e.metricsAddr,
	)

	// N.B: we want to fail hard if users don't specify the autoTLS and cert/key pair isn't there.
	// and warn users to use the tls.auto flag if they don't want to manage their certs.
	grp, ctx := errgroup.WithContext(e.ctx)
	switch {
	case !e.tlsDisabled && !e.autoTLS:
		log.Info("running exporter with TLS connection")
		// start server using the given certs
		grp.Go(func() error {
			return server.StartTLS(fmt.Sprintf(":%d", e.hostPort), e.certFile, e.keyFile)
		})
	case !e.tlsDisabled && e.autoTLS:
		log.Info("running exporter with auto TLS connection")
		// start server using managed certs
		// Cache certificates to avoid issues with rate limits (https://letsencrypt.org/docs/rate-limits)
		server.AutoTLSManager.HostPolicy = autocert.HostWhitelist(e.tlsHosts...)
		server.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
		grp.Go(func() error {
			return server.StartAutoTLS(fmt.Sprintf(":%d", e.hostPort))
		})
	default:
		log.Info("running exporter with insecure connection!")
		grp.Go(func() error {
			return server.Start(fmt.Sprintf(":%d", e.hostPort))
		})
	}
	grp.Go(func() error {
		<-ctx.Done()
		return server.Shutdown(context.Background())
	})

	// if profiling is enabled start the pprof server
	if e.profilingEnabled {
		log.Infof(
			"pprof profiling has been activated on port :%d", e.profilingPort,
		)
		profilingServer := echo.New()
		grp.Go(func() error {
			return profilingServer.Start(fmt.Sprintf(":%d",  e.profilingPort))
		})
		grp.Go(func() error {
			<-ctx.Done()
			return profilingServer.Shutdown(context.Background())
		})
	}

	// exporter's own metrics
	log.Infof("starting metrics server on port :%d at /%s", e.probePort, e.probeAddr)
	probeServer := echo.New()
	probeServer.GET(fmt.Sprintf("/%s", e.probeAddr), e.probe())
	grp.Go(func() error {
		return probeServer.Start(fmt.Sprintf(":%d", e.probePort))
	})
	grp.Go(func() error {
		<-ctx.Done()
		return probeServer.Shutdown(context.Background())
	})

	return grp.Wait()
}
