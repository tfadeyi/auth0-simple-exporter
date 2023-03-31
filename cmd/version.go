package cmd

import (
	"github.com/auth0-simple-exporter/pkg/logging"
	"github.com/auth0-simple-exporter/pkg/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Returns the binary build information.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		log := logging.LoggerFromContext(ctx)
		log.Info(version.BuildInfo())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
