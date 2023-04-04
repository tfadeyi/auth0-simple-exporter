package cmd

import (
	swaggeroptions "github.com/auth0-simple-exporter/cmd/options/swagger"
	"github.com/auth0-simple-exporter/pkg/swagger"
	"github.com/spf13/cobra"
)

func serveSwaggerCmd() *cobra.Command {
	opts := swaggeroptions.New()
	cmd := &cobra.Command{
		Use:           "swagger",
		Short:         "Starts a swagger docs local server.",
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
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return swagger.Serve(ctx, opts.HostPort)
		},
	}
	opts = opts.Prepare(cmd)
	return cmd
}

func init() {
	rootCmd.AddCommand(serveSwaggerCmd())
}
