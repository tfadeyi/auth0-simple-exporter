package logfetcher

import (
	"context"
	"github.com/auth0-simple-exporter/internal/auth0"
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"go.uber.org/multierr"
)

type (
	// LogFetcher Fetcher implementation for *management.Log
	LogFetcher struct {
		client *management.Management
	}
)

func (c *LogFetcher) FetchAll(ctx context.Context) (interface{}, error) {
	var allLogs []*management.Log
	i := 0
	for {
		logs, err := c.client.Log.List(
			management.Context(ctx),
			management.IncludeFields("type", "log_id"),
			management.Page(i),
			management.PerPage(100),
		)
		if err != nil {
			return nil, err
		}
		allLogs = append(allLogs, logs...)
		if len(logs) < 100 {
			return allLogs, nil
		}
		i++
	}

	return allLogs, nil
}

// NewFetcherWithOpts creates a new instance of the LogFetcher using the given options
func NewFetcherWithOpts(opts auth0.Options) (*LogFetcher, error) {
	var errs error
	var client *management.Management

	if opts.Domain == "" {
		errs = multierr.Append(errs, errors.New("missing auth0 domain"))
	}
	if opts.Token != "" {
		c, err := management.New(opts.Domain, management.WithStaticToken(opts.Token))
		if err != nil {
			errs = multierr.Append(errs, err)
		}
		client = c
	}
	if opts.ClientID != "" && opts.ClientSecret != "" {
		c, err := management.New(opts.Domain, management.WithClientCredentials(opts.ClientID, opts.ClientSecret))
		if err != nil {
			errs = multierr.Append(errs, err)
		}
		client = c
	}

	if errs != nil || client == nil {
		return nil, errs
	}

	return &LogFetcher{client: client}, nil
}
