package swagger

import (
	"context"
	"fmt"

	_ "github.com/auth0-simple-exporter/pkg/docs"
	"github.com/auth0-simple-exporter/pkg/logging"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"golang.org/x/sync/errgroup"
)

func Serve(ctx context.Context, port int) error {
	log := logging.LoggerFromContext(ctx)
	log.Info("starting swagger server")

	server := echo.New()
	server.HideBanner = true
	server.Use(middleware.Recover())

	server.GET("/swagger/*", echoSwagger.WrapHandler)
	grp, ctx := errgroup.WithContext(ctx)
	grp.Go(func() error {
		return server.Start(fmt.Sprintf(":%d", port))
	})
	grp.Go(func() error {
		<-ctx.Done()
		return server.Shutdown(context.Background())
	})

	return grp.Wait()
}
