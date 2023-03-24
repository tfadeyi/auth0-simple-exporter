package options

import (
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"

	"github.com/auth0-simple-exporter/internal/auth0"
	fetch "github.com/auth0-simple-exporter/internal/auth0/logfetcher"
)

type (
	// Options is the list of options/flag available to the application,
	// plus the clients needed by the application to function.
	Options struct {
		ProfilingEnabled bool
		MetricsEndpoint  string
		HostPort    int
		TLSDisabled bool
		ManagedTLS  bool
		CertFile    string
		KeyFile          string

		// exporter
		Namespace string

		// Auth0 setup
		cfg    auth0.Options
		Client auth0.Fetcher
	}
)

const (
	envClientSecret = "CLIENT_SECRET"
	envClientID     = "CLIENT_ID"
	envMgmtToken    = "TOKEN"
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

// Complete initialises the components needed for the application to function given the options,
// like the auth0 client.
func (o *Options) Complete() error {
	var err error

	o.Client, err = fetch.NewFetcherWithOpts(o.cfg)
	if err != nil {
		return errors.Annotate(err, "failed to initialise exporter's connection to auth0")
	}

	return nil
}

func (o *Options) addAppFlags(fs *pflag.FlagSet) {
	fs.BoolVar(
		&o.ProfilingEnabled,
		"pprof",
		false,
		"Enabled pprof profiling on the exporter on port :6060. (help: https://jvns.ca/blog/2017/09/24/profiling-go-with-pprof/)",
	)
	fs.BoolVar(
		&o.TLSDisabled,
		"tls.disabled",
		false,
		"",
	)
	fs.BoolVar(
		&o.ManagedTLS,
		"tls.managed",
		false,
		"Allow the exporter manage its own certificates.",
	)
	fs.StringVar(
		&o.CertFile,
		"tls.cert-file",
		"",
		"The certificate file for the exporter.",
	)
	fs.StringVar(
		&o.KeyFile,
		"tls.key-file",
		"",
		"The key file for the exporter.",
	)
	fs.StringVar(
		&o.cfg.Domain,
		"auth0.domain",
		"",
		"Auth0 tenant's domain. (i.e: <tenant_name>.eu.auth0.com)",
	)
	fs.StringVar(
		&o.cfg.Token,
		"auth0.token",
		os.Getenv(envMgmtToken),
		"Auth0 management api static token",
	)
	fs.StringVar(
		&o.cfg.ClientID,
		"auth0.client-id",
		os.Getenv(envClientID),
		"Auth0 management api client-id",
	)
	fs.StringVar(
		&o.cfg.ClientSecret,
		"auth0.client-secret",
		os.Getenv(envClientSecret),
		"Auth0 management api static token.",
	)
	fs.StringVar(
		&o.Namespace,
		"namespace",
		"",
		"Exporter's namespace",
	)
	fs.IntVar(
		&o.HostPort,
		"web.listen-address",
		8081,
		"Port where the server will listen.",
	)
	fs.StringVar(
		&o.MetricsEndpoint,
		"web.metrics-path",
		"metrics",
		"URL Path under which to expose metrics.",
	)
}
