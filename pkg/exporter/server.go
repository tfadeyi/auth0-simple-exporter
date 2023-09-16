package exporter

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tfadeyi/auth0-simple-exporter/pkg/exporter/metrics"
	"github.com/tfadeyi/auth0-simple-exporter/pkg/logging"
	"github.com/tfadeyi/auth0-simple-exporter/pkg/version"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/sync/errgroup"
)

//	@title			Auth0 simple exporter
//	@version		0.1.1
//	@description	A simple Prometheus exporter for Auth0 log [events](https://auth0.com/docs/api/management/v2#!/Logs/get_logs),
//	@description	which allows you to collect metrics from Auth0 and expose them in a format that can be consumed by Prometheus.

//	@contact.name	Oluwole Fadeyi (@tfadeyi)

//	@license.name	Apache 2.0
//	@license.url	https://github.com/tfadeyi/auth0-simple-exporter/blob/main/LICENSE

//	@host		localhost:8080
//	@BasePath	/

// Export configures the exporter Router and starts the server with the given configuration
func (e *exporter) Export() error {
	log := e.logger
	log.Info("initialising exporter", "build", version.Info())

	// collect all the auth0 apps/clients, this is required for initialising the prometheus counter metrics to zero
	// all labels must be known before.
	list, err := e.client.App.List(e.ctx)
	if err != nil {
		return errors.Annotate(err, "error fetching the auth0 tenant client applications.")
	}
	applications, ok := list.([]*management.Client)
	if !ok {
		return errors.New("auth0 client applications fetch didn't return the expected list of applications client type")
	}

	server := echo.New()
	server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return metrics.NamespaceMiddleware(next, e.namespace)
	})
	server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return metrics.SubsystemMiddleware(next, e.subsystem)
	})
	server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return metrics.Middleware(next, applications)
	})
	server.Use(middleware.Recover())
	server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return logging.Middleware(next, log)
	})
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

	log.Info("starting exporter",
		"port", e.hostPort,
		"metrics-address", "/"+e.metricsAddr,
	)

	// N.B: we want to fail hard if users don't specify the autoTLS and cert/key pair isn't there.
	// and warn users to use the tls.auto flag if they don't want to manage their certs.
	grp, ctx := errgroup.WithContext(e.ctx)
	switch {
	case !e.tlsDisabled && !e.autoTLS:
		e.logger.Info("running exporter with TLS connection")
		// start server using the given certs
		grp.Go(func() error {
			return server.StartTLS(fmt.Sprintf(":%d", e.hostPort), e.certFile, e.keyFile)
		})
	case !e.tlsDisabled && e.autoTLS:
		e.logger.Info("running exporter with auto TLS connection")
		// start server using managed certs
		// Cache certificates to avoid issues with rate limits (https://letsencrypt.org/docs/rate-limits)
		server.AutoTLSManager.HostPolicy = autocert.HostWhitelist(e.tlsHosts...)
		server.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
		grp.Go(func() error {
			return server.StartAutoTLS(fmt.Sprintf(":%d", e.hostPort))
		})
	default:
		e.logger.Info("running exporter with insecure connection!")
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
		e.logger.Info("pprof profiling is activate", "port", e.profilingPort)
		profilingServer := echo.New()
		grp.Go(func() error {
			return profilingServer.Start(fmt.Sprintf(":%d", e.profilingPort))
		})
		grp.Go(func() error {
			<-ctx.Done()
			return profilingServer.Shutdown(context.Background())
		})
	}

	// exporter's own metrics
	e.logger.Info("starting metrics server", "port", e.probePort, "endpoint", e.probeAddr)
	probeServer := echo.New()
	probeServer.HideBanner = true
	probeServer.HidePort = true
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
