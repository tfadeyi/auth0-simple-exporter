package cmd

import (
	"github.com/spf13/cobra"
	versionoptions "github.com/tfadeyi/auth0-simple-exporter/cmd/options/version"
	"github.com/tfadeyi/auth0-simple-exporter/pkg/logging"
	"github.com/tfadeyi/auth0-simple-exporter/pkg/version"
)

func versionCmd() *cobra.Command {
	opts := versionoptions.New()
	cmd := &cobra.Command{
		Use:           "version",
		Short:         "Returns the binary build information.",
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := opts.Validate()
			if err != nil {
				return err
			}
			err = opts.Complete()
			return err
		},
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			log := logging.LoggerFromContext(ctx)
			if opts.Verbose {
				log.Info(version.BuildInfo())
				return
			}
			log.Info(version.Info())
		},
	}
	opts = opts.Prepare(cmd)
	return cmd
}

func init() {
	rootCmd.AddCommand(versionCmd())
}
