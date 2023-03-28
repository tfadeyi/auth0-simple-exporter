package options

import (
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
	"time"

	"github.com/auth0-simple-exporter/internal/auth0"
	fetch "github.com/auth0-simple-exporter/internal/auth0/logfetcher"
)

type (
	// Options is the list of options/flag available to the application,
	// plus the clients needed by the application to function.
	Options struct {
		ProfilingEnabled    bool
		ProfilingPort       int
		MetricsEndpoint     string
		HostPort            int
		LogLevel            int8
		StartTimeCheckpoint string

		// probe
		ProbePort int
		ProbeAddr string

		// TLS
		TLSDisabled bool
		AutoTLS     bool
		CertFile    string
		KeyFile     string
		TLSHosts    []string

		// exporter
		Namespace string
		Subsystem string

		// Auth0 setup
		cfg    auth0.Options
		Client auth0.Fetcher
	}
)

const (
	envClientSecret = "CLIENT_SECRET"
	envClientID     = "CLIENT_ID"
	envMgmtToken    = "TOKEN"
	envDomainToken  = "DOMAIN"
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

	if !o.TLSDisabled {
		// only check cert/key if they don't exist if TLS is not disabled
		if _, err := os.Stat(o.KeyFile); errors.Is(err, os.ErrNotExist) {
			return errors.New("failed to find the exporter's private key file. TLS can be disabled with --tls.disabled")
		}

		if _, err := os.Stat(o.CertFile); errors.Is(err, os.ErrNotExist) {
			return errors.New("failed to find the exporter's certificate file. TLS can be disabled with --tls.disabled")
		}
	}

	return nil
}

func (o *Options) addAppFlags(fs *pflag.FlagSet) {
	fs.Int8Var(
		&o.LogLevel,
		"log.level",
		1,
		"Exporter log level.",
	)

	fs.IntVar(
		&o.ProfilingPort,
		"pprof.listen-address",
		6060,
		"Port where the pprof webserver will listen on.",
	)
	fs.BoolVar(
		&o.ProfilingEnabled,
		"pprof.enabled",
		false,
		"Enabled pprof profiling on the exporter on port :6060. (help: https://jvns.ca/blog/2017/09/24/profiling-go-with-pprof/).",
	)

	fs.BoolVar(
		&o.TLSDisabled,
		"tls.disabled",
		false,
		"Run exporter without TLS.",
	)
	fs.BoolVar(
		&o.AutoTLS,
		"tls.auto",
		false,
		`Allow the exporter to use autocert to renew its certificates with letsencrypt.
(Can only be used if the exporter is publicly accessible by the internet)`,
	)
	fs.StringVar(
		&o.CertFile,
		"tls.cert-file",
		"",
		"The certificate file for the exporter TLS connection.",
	)
	fs.StringVar(
		&o.KeyFile,
		"tls.key-file",
		"",
		"The key file for the exporter TLS connection.",
	)
	fs.StringSliceVar(
		&o.TLSHosts,
		"tls.hosts",
		[]string{},
		"The different allowed hosts for the exporter. Only works when --tls.auto has been enabled.",
	)

	fs.StringVar(
		&o.StartTimeCheckpoint,
		"auth0.checkpoint",
		time.Now().Format("2006-01-02"),
		"Point in time from were to start fetching auth0 logs. (format: YYYY-MM-DD)",
	)
	fs.StringVar(
		&o.cfg.Domain,
		"auth0.domain",
		os.Getenv(envDomainToken),
		"Auth0 tenant's domain. (i.e: <tenant_name>.eu.auth0.com).",
	)
	fs.StringVar(
		&o.cfg.Token,
		"auth0.token",
		os.Getenv(envMgmtToken),
		"Auth0 management api static token. (the token can be used instead of client credentials).",
	)
	fs.StringVar(
		&o.cfg.ClientID,
		"auth0.client-id",
		os.Getenv(envClientID),
		"Auth0 management api client-id.",
	)
	fs.StringVar(
		&o.cfg.ClientSecret,
		"auth0.client-secret",
		os.Getenv(envClientSecret),
		"Auth0 management api client-secret.",
	)

	fs.StringVar(
		&o.Namespace,
		"namespace",
		"",
		"Exporter's namespace.",
	)
	fs.StringVar(
		&o.Subsystem,
		"subsystem",
		"",
		"Exporter's subsystem.",
	)

	fs.IntVar(
		&o.HostPort,
		"web.listen-address",
		8080,
		"Port where the exporter webserver will listen on.",
	)
	fs.StringVar(
		&o.MetricsEndpoint,
		"web.path",
		"metrics",
		"URL Path under which to expose the collected auth0 metrics.",
	)

	fs.IntVar(
		&o.ProbePort,
		"probe.listen-address",
		8081,
		"Port where the probe webserver will listen on.",
	)
	fs.StringVar(
		&o.ProbeAddr,
		"probe.path",
		"probe",
		"URL Path under which to expose the probe metrics.",
	)
}
