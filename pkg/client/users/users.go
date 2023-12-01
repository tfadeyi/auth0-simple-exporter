package users

import (
	"context"

	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"go.uber.org/multierr"
)

//go:generate moq -out users_mock.go . Client

var ErrAPIRateLimitReached = errors.New("client reached api rate limit")

const ItemCountPerPage = 100

type (
	Client interface {
		List(ctx context.Context, args ...interface{}) (interface{}, error)
	}

	usersClient struct {
		mgmt userManagement
	}
)

func (l *usersClient) List(ctx context.Context, args ...interface{}) (interface{}, error) {
	var allUsers []*management.User
	page := 0
	hasNext := true
	for hasNext {
		users, err := l.mgmt.List(
			ctx,
			management.IncludeFields("user_id", "blocked", "last_login"),
			management.PerPage(ItemCountPerPage),
			management.Page(page),
		)
		switch {
		case errors.Is(err, errors.QuotaLimitExceeded):
			return nil, ErrAPIRateLimitReached
		case err != nil:
			return nil, err
		}
		allUsers = append(allUsers, users.Users...)
		hasNext = users.HasNext()
		page++
	}
	return allUsers, nil
}

// New returns a new instance of the clients fetching client, plus possible errors
func New(domain, clientID, clientSecret, token string) (*usersClient, error) {
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
	return &usersClient{client.User}, nil
}
