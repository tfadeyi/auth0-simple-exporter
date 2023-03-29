package auth0

import (
	"context"
	"errors"
	"time"
)

//go:generate rm -f ./mock.go
//go:generate moq -out mock.go . Fetcher

type (
	Fetcher interface {
		// FetchAll returns all resources of a specif type
		FetchAll(ctx context.Context, startTime time.Time) (interface{}, error)
	}

	// Options to configure the Fetcher
	Options struct {
		Domain       string
		Token        string
		ClientSecret string
		ClientID     string
	}
)

var (
	ErrAPIRateLimitReached = errors.New("auth0 api rate limit was reached")
)