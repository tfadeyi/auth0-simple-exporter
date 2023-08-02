package applications

import (
	"context"
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"go.uber.org/multierr"
)

//go:generate moq -out apps_mock.go . Client

var ErrAPIRateLimitReached = errors.New("client reached api rate limit")

type (
	Client interface {
		List(ctx context.Context, args ...interface{}) (interface{}, error)
	}

	applicationClient struct {
		mgmt applicationManagement
	}
)

func (l *applicationClient) List(ctx context.Context, args ...interface{}) (interface{}, error) {
	var allApplications []*management.Client
	page := 0
	hasNext := true
	for hasNext {
		apps, err := l.mgmt.List(
			ctx,
			management.IncludeFields("name"),
			management.PerPage(100),
			management.Page(page),
		)
		switch {
		case errors.Is(err, errors.QuotaLimitExceeded):
			return nil, ErrAPIRateLimitReached
		case err != nil:
			return nil, err
		}
		allApplications = append(allApplications, apps.Clients...)
		hasNext = apps.HasNext()
		page++
	}
	return allApplications, nil
}

// New returns a new instance of the application fetching client, plus possible errors
func New(domain, clientID, clientSecret, token string) (*applicationClient, error) {
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
	return &applicationClient{client.Client}, nil
}
