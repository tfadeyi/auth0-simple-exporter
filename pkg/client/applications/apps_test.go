package applications

import (
	"context"
	"testing"
	"time"

	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	t.Parallel()

	t.Run("failed to initialise the applications client, if auth0.Domain is missing", func(t *testing.T) {
		_, err := New("", "some id", "some secret", "random token")
		require.Error(t, err)
	})
	t.Run("failed to initialise the applications client, with invalid client credentials and token secrets", func(t *testing.T) {
		_, err := New("domain", "some id", "", "")
		require.Error(t, err)
	})
	t.Run("successfully to initialise the applications client, if auth0.Token is present and client credentials are missing", func(t *testing.T) {
		_, err := New("domain", "", "", "random token")
		require.NoError(t, err)
	})
	t.Run("successfully to initialise the applications client, if auth0 client creds are present and token is missing", func(t *testing.T) {
		_, err := New("domain", "id", "secret", "")
		require.NoError(t, err)
	})

	t.Run("fail if client returns rate limit errors", func(t *testing.T) {
		c := applicationClient{mgmt: &applicationManagementMock{
			ListFunc: func(ctx context.Context, opts ...management.RequestOption) (*management.ClientList, error) {
				return nil, errors.QuotaLimitExceededf("api request limit was reached")
			},
		}}
		_, err := c.List(context.Background(), time.Now())
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrAPIRateLimitReached)
	})
}
