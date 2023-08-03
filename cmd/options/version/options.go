package version

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type (
	// Options is the list of options/flag available to the application,
	// plus the clients needed by the application to function.
	Options struct {
		Verbose bool
	}
)

// New creates a new instance of the application's options
func New() *Options {
	return new(Options)
}

// Prepare assigns the applications flag/options to the cobra cli
func (o *Options) Prepare(cmd *cobra.Command) *Options {
	o.addAppFlags(cmd.Flags())
	return o
}

// Validate validates the flag values given to the application
func (o *Options) Validate() error {
	return nil
}

// Complete initialises the components needed for the application to function given the options,
func (o *Options) Complete() error {
	return nil
}

func (o *Options) addAppFlags(fs *pflag.FlagSet) {
	fs.BoolVarP(
		&o.Verbose,
		"version",
		"v",
		false,
		"Verbose build information",
	)
}
