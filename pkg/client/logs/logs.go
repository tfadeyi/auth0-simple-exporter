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

// max number of items returned by Auth0 for each API call
const ItemCountPerPage = 50

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

	for {
		query := fmt.Sprintf("date:{%s TO *]", from.UTC().Format(time.RFC3339))

		logs, err := l.mgmt.List(
			ctx,
			management.IncludeFields("type", "log_id", "date", "client_name"),
			management.Query(query),
			management.Sort("date:1"),
			management.Take(ItemCountPerPage),
		)

		switch {
		case errors.Is(err, errors.QuotaLimitExceeded):
			return allLogs, ErrAPIRateLimitReached
		case err != nil:
			return allLogs, err
		}

		if len(logs) == 0 {
			break
		} else if len(logs) == ItemCountPerPage {
			// the last item is used as checkpoint (it will be the first
			// of the next response)
			from = *logs[len(logs)-1].Date
			logs = logs[:len(logs)-1]

			allLogs = append(allLogs, logs...)
		} else {
			allLogs = append(allLogs, logs...)
			break
		}
	}

	return allLogs, nil
}
