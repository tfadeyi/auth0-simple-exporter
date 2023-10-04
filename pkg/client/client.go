package client

import (
	"github.com/tfadeyi/auth0-simple-exporter/pkg/client/applications"
	"github.com/tfadeyi/auth0-simple-exporter/pkg/client/logs"
	"github.com/tfadeyi/auth0-simple-exporter/pkg/client/users"
)

type (
	// Options to configure the Client
	Options struct {
		Domain       string
		Token        string
		ClientSecret string
		ClientID     string
	}

	Client struct {
		Log  logs.Client
		App  applications.Client
		User users.Client
	}
)

// NewWithOpts creates a new instance of the client using the given options
func NewWithOpts(opts Options) (Client, error) {
	// create client to retrieve logs
	c, err := logs.New(opts.Domain, opts.ClientID, opts.ClientSecret, opts.Token)
	if err != nil {
		return Client{}, err
	}
	// create client to retrieve apps
	appClient, err := applications.New(opts.Domain, opts.ClientID, opts.ClientSecret, opts.Token)
	if err != nil {
		return Client{}, err
	}
	// create client to retrieve users
	userClient, err := users.New(opts.Domain, opts.ClientID, opts.ClientSecret, opts.Token)
	if err != nil {
		return Client{}, err
	}
	return Client{Log: c, App: appClient, User: userClient}, nil
}
