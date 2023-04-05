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
		c, err := management.New(domain, management.WithClientCredentials(clientID, clientSecret))
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

	i := 0
	for {
		var logs []*management.Log
		var err error

		if from.IsZero() {
			logs, err = l.mgmt.List(
				management.Context(ctx),
				management.IncludeFields("type", "log_id", "date"),
				management.Page(i),
				management.PerPage(100),
			)
		} else {
			logs, err = l.mgmt.List(
				management.Context(ctx),
				management.IncludeFields("type", "log_id", "date"),
				management.Page(i),
				management.PerPage(100),
				management.Query(fmt.Sprintf("date:[%s TO *]", from.Format("2006-01-02"))),
			)
		}
		switch {
		case errors.Is(err, errors.QuotaLimitExceeded):
			return nil, ErrAPIRateLimitReached
		case err != nil:
			return nil, err
		}
		allLogs = append(allLogs, logs...)
		if len(logs) < 100 {
			return allLogs, nil
		}
		i++
	}
}
