package client

import (
	"github.com/tfadeyi/auth0-simple-exporter/pkg/client/logs"
)

type (
	// Options to configure the Client
	Options struct {
		Domain       string
		Token        string
		ClientSecret string
		ClientID     string
		// Client, if set it overrides the auth0 client used in the wrapper
		Client logs.Client
	}

	Client struct {
		logs.Client
	}
)

// NewWithOpts creates a new instance of the client using the given options
func NewWithOpts(opts Options) (Client, error) {
	if opts.Client != nil {
		return Client{Client: opts.Client}, nil
	}
	// create client to retrieve logs
	c, err := logs.New(opts.Domain, opts.ClientID, opts.ClientSecret, opts.Token)
	if err != nil {
		return Client{}, err
	}
	return Client{Client: c}, nil
}
