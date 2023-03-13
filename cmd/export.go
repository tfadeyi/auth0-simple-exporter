package cmd

import (
	"github.com/auth0-simple-exporter/cmd/options"
	"github.com/auth0-simple-exporter/internal/exporter"
	"github.com/auth0-simple-exporter/internal/logging"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

// serveAppCmd represents the serve command for the application server
func serveExporterCmd() *cobra.Command {
	opts := options.New()
	cmd := &cobra.Command{
		Use:           "run",
		Short:         "Start serving the auth0 metrics",
		Long:          `This starts the exporter HTTP server on the given port.`,
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.Complete()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			log := logging.LoggerFromContext(ctx)
			log.Info("initialising exporter")

			//go run main.go run --auth0.domain jetstack-stage.eu.auth0.com --profiling
			e, err := exporter.New(
				exporter.Context(ctx),
				exporter.Port(opts.HostPort),
				exporter.MetricsAddr(opts.MetricsEndpoint),
				exporter.Profiling(opts.ProfilingEnabled),
				exporter.Client(opts.Client))
			if err != nil {
				return errors.Annotate(err, "failed to initialise the exporter")
			}
			return e.Export()
		},
	}
	opts = opts.Prepare(cmd)
	return cmd
}

func init() {
	rootCmd.AddCommand(serveExporterCmd())
}
