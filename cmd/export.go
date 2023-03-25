package cmd

import (
	"github.com/auth0-simple-exporter/cmd/options"
	"github.com/auth0-simple-exporter/internal/exporter"
	"github.com/auth0-simple-exporter/internal/logging"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

func serveExporterCmd() *cobra.Command {
	opts := options.New()
	cmd := &cobra.Command{
		Use:           "export",
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

			e, err := exporter.New(
				exporter.Context(ctx),
				exporter.Port(opts.HostPort),
				exporter.MetricsAddr(opts.MetricsEndpoint),
				exporter.Profiling(opts.ProfilingEnabled),
				exporter.Client(opts.Client),
				exporter.Namespace(opts.Namespace),
				exporter.CertFile(opts.CertFile),
				exporter.DisableTLS(opts.TLSDisabled),
				exporter.KeyFile(opts.KeyFile),
				exporter.ManagedTLS(opts.ManagedTLS),
				exporter.TLSHost(opts.TLSHost))
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
