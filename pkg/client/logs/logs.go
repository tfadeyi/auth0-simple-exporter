package logs

import (
	"context"
	"fmt"
	"time"

	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"go.uber.org/multierr"
)

//go:generate moq -out logs_mock.go . Client

type (
	Client interface {
		List(ctx context.Context, args ...interface{}) (interface{}, error)
	}

	logClient struct {
		mgmt logManagement
	}
)

var ErrAPIRateLimitReached = errors.New("client reached api rate limit")

// New returns a new instance of the log fetching client, plus possible errors
func New(domain, clientID, clientSecret, token string) (*logClient, error) {
	var errs error
	var client *management.Management
	if domain == "" {
		errs = multierr.Append(errs, errors.New("missing auth0 domain"))
	}
	if token != "" {
		c, err := management.New(domain, management.WithStaticToken(token))
		if err != nil {
			errs = multierr.Append(errs, err)
		}
		client = c
	}
	if clientID != "" && clientSecret != "" {
		c, err := management.New(domain, management.WithClientCredentials(context.Background(), clientID, clientSecret))
		if err != nil {
			errs = multierr.Append(errs, err)
		}
		client = c
	}
	if client == nil {
		errs = multierr.Append(errs, errors.New("unable to initialise the auth0 client, check the credentials are correct."))
	}

	if errs != nil {
		return nil, errs
	}
	return &logClient{client.Log}, nil
}

func (l *logClient) List(ctx context.Context, args ...interface{}) (interface{}, error) {
	var allLogs []*management.Log
	var from time.Time
	if len(args) > 0 {
		var ok bool
		from, ok = args[0].(time.Time)
		if !ok {
			return nil, errors.New("invalid \"from\" argument passed to the client")
		}
	}

	// Get the last log from the list of logs for the previous day.
	// This is used as the starting point for the fetching of the logs.
	// This allows us to use the checkpoint pagination style
	var checkpoint *management.Log
	if !from.IsZero() {
		previousDay := from.Add(-24 * time.Hour).Format("2006-01-02")
		logs, err := l.mgmt.List(
			ctx,
			management.IncludeFields("type", "log_id", "date", "client_name"),
			management.PerPage(1),
			management.Query(fmt.Sprintf("date:[* TO %s]", previousDay)),
		)
		switch {
		case errors.Is(err, errors.QuotaLimitExceeded):
			return allLogs, ErrAPIRateLimitReached
		case err != nil:
			return allLogs, err
		case len(logs) > 0:
			checkpoint = logs[0]
		}
	}

	for {
		logs, err := l.fetchLogs(ctx, checkpoint)
		switch {
		case errors.Is(err, errors.QuotaLimitExceeded):
			return allLogs, ErrAPIRateLimitReached
		case err != nil:
			return allLogs, err
		}
		allLogs = append(allLogs, logs...)

		if len(logs) == 0 {
			return allLogs, nil
		}
		checkpoint = logs[len(logs)-1]
	}
}

// fetchLogs returns the list of logs given a starting checkpoint. If no checkpoint is passed it returns the list of the
// latest logs. (Default: 100 items)
func (l *logClient) fetchLogs(ctx context.Context, checkpoint *management.Log) ([]*management.Log, error) {
	if checkpoint != nil {
		return l.mgmt.List(
			ctx,
			management.IncludeFields("type", "log_id", "date", "client_name"),
			management.From(checkpoint.GetLogID()),
			management.Take(100),
		)
	}
	return l.mgmt.List(
		ctx,
		management.IncludeFields("type", "log_id", "date", "client_name"),
		management.Take(100),
	)
}
