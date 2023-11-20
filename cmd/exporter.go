package cmd

import (
	"fmt"
	"time"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
	exporteroptions "github.com/tfadeyi/auth0-simple-exporter/cmd/options/exporter"
	"github.com/tfadeyi/auth0-simple-exporter/pkg/exporter"
	"github.com/tfadeyi/auth0-simple-exporter/pkg/exporter/format"
	"github.com/tfadeyi/auth0-simple-exporter/pkg/logging"
)

// serveExporterCmd is the entrypoint to the exporter
func serveExporterCmd() *cobra.Command {
	opts := exporteroptions.New()
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
			log := logging.NewPromLoggerWithOpts(opts.LogLevel)
			from, err := time.Parse(format.TimeFormat, opts.FromFetchTime)
			if err != nil {
				opts.FromFetchTime = fmt.Sprintf("%sT%d:%d", opts.FromFetchTime, 23, 59)
				from, err = time.Parse(format.TimeFormat, opts.FromFetchTime)
				if err != nil {
					return errors.Annotate(err, "could not parse time in --auth0.from, make sure the format is correct")
				}
			}

			e := exporter.New(
				ctx,
				exporter.Port(opts.HostPort),
				exporter.MetricsAddr(opts.MetricsEndpoint),
				exporter.Profiling(opts.ProfilingEnabled),
				exporter.ProfilingPort(opts.ProfilingPort),
				exporter.Client(opts.Client),
				exporter.Namespace(opts.Namespace),
				exporter.CertFile(opts.CertFile),
				exporter.DisableTLS(opts.TLSDisabled),
				exporter.KeyFile(opts.KeyFile),
				exporter.AutoTLS(opts.AutoTLS),
				exporter.TLSHosts(opts.TLSHosts),
				exporter.ProbeAddr(opts.ProbeAddr),
				exporter.ProbePort(opts.ProbePort),
				exporter.From(from),
				exporter.Logger(log),
			)
			return e.Export()
		},
	}
	opts = opts.Prepare(cmd)
	return cmd
}

func init() {
	rootCmd.AddCommand(serveExporterCmd())
}
