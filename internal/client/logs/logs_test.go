package logs

import (
	"context"
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	t.Parallel()

	t.Run("failed to initialize the log client, if auth0.Domain is missing", func(t *testing.T) {
		_, err := New("", "some id", "some secret", "random token")
		require.Error(t, err)
	})
	t.Run("failed to initialize the log client, with invalid client credentials and token secrets", func(t *testing.T) {
		_, err := New("domain", "some id", "", "")
		require.Error(t, err)
	})
	t.Run("successfully to initialize the log client, if auth0.Token is present and client credentials are missing", func(t *testing.T) {
		_, err := New("domain", "", "", "random token")
		require.NoError(t, err)
	})
	t.Run("successfully to initialize the log client, if auth0 client creds are present and token is missing", func(t *testing.T) {
		_, err := New("domain", "id", "secret", "")
		require.NoError(t, err)
	})

	t.Run("fail if List \"from\" argument is not a time.Time type", func(t *testing.T) {
		c := logClient{mgmt: &logManagementMock{
			ListFunc: func(opts ...management.RequestOption) ([]*management.Log, error) {
				return nil, nil
			},
		}}
		_, err := c.List(context.Background(), "not a valid time type")
		require.Error(t, err)
	})

	t.Run("fail if client returns rate limit errors", func(t *testing.T) {
		c := logClient{mgmt: &logManagementMock{
			ListFunc: func(opts ...management.RequestOption) ([]*management.Log, error) {
				return nil, errors.QuotaLimitExceededf("api request limit was reached")
			},
		}}
		_, err := c.List(context.Background(), time.Now())
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrAPIRateLimitReached)
	})

	t.Run("successfully fetch all logs across multiple pages with client.List", func(t *testing.T) {
		totalLogNumber := 4299
		pageSize := 100
		currentPage := 0
		storedLogs := make([]*management.Log, totalLogNumber)
		for i := 0; i < totalLogNumber; i++ {
			var code string = "f"
			storedLogs[i] = &management.Log{Type: &code}
		}

		c := logClient{mgmt: &logManagementMock{
			ListFunc: func(opts ...management.RequestOption) ([]*management.Log, error) {
				var result []*management.Log
				for i := range storedLogs {
					index := i + (currentPage * pageSize)
					if i == pageSize || (index >= len(storedLogs)) {
						break
					}
					result = append(result, storedLogs[index])
				}

				currentPage++
				return result, nil
			},
		}}
		actualLogs, err := c.List(context.Background(), time.Now())
		require.NoError(t, err)
		assert.Len(t, actualLogs, totalLogNumber)
	})
}
