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
	Options struct {
		ProfilingEnabled bool
		MetricsEndpoint  string
		HostPort         int

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

func New() *Options {
	return new(Options)
}

func (o *Options) Prepare(cmd *cobra.Command) *Options {
	o.addAppFlags(cmd.Flags())
	return o
}

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
		"profiling",
		false,
		"Enabled pprof profiling on the exporter on port :6060. (help: https://jvns.ca/blog/2017/09/24/profiling-go-with-pprof/)",
	)
	fs.StringVar(
		&o.cfg.Domain,
		"auth0.domain",
		"",
		"Auth0 tenant's domain.",
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
