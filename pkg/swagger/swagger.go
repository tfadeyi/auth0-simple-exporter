package swagger

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoswagger "github.com/swaggo/echo-swagger"
	_ "github.com/tfadeyi/auth0-simple-exporter/pkg/docs"
	"github.com/tfadeyi/auth0-simple-exporter/pkg/logging"
	"golang.org/x/sync/errgroup"
)

func Serve(ctx context.Context, port int) error {
	log := logging.LoggerFromContext(ctx)
	log.Info("starting swagger server")

	server := echo.New()
	server.HideBanner = true
	server.Use(middleware.Recover())

	server.GET("/swagger/*", echoswagger.WrapHandler)
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
