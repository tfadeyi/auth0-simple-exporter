package logs

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	t.Parallel()

	t.Run("failed to initialise the log client, if auth0.Domain is missing", func(t *testing.T) {
		_, err := New("", "some id", "some secret", "random token")
		require.Error(t, err)
	})
	t.Run("failed to initialise the log client, with invalid client credentials and token secrets", func(t *testing.T) {
		_, err := New("domain", "some id", "", "")
		require.Error(t, err)
	})
	t.Run("successfully to initialise the log client, if auth0.Token is present and client credentials are missing", func(t *testing.T) {
		_, err := New("domain", "", "", "random token")
		require.NoError(t, err)
	})
	t.Run("successfully to initialise the log client, if auth0 client creds are present and token is missing", func(t *testing.T) {
		_, err := New("domain", "id", "secret", "")
		require.NoError(t, err)
	})

	t.Run("fail if List \"from\" argument is not a time.Time type", func(t *testing.T) {
		c := logClient{mgmt: &logManagementMock{
			ListFunc: func(ctx context.Context, opts ...management.RequestOption) ([]*management.Log, error) {
				return nil, nil
			},
		}}
		_, err := c.List(context.Background(), "not a valid time type")
		require.Error(t, err)
	})

	t.Run("fail if client returns rate limit errors", func(t *testing.T) {
		c := logClient{mgmt: &logManagementMock{
			ListFunc: func(ctx context.Context, opts ...management.RequestOption) ([]*management.Log, error) {
				return nil, errors.QuotaLimitExceededf("api request limit was reached")
			},
		}}
		_, err := c.List(context.Background(), time.Now())
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrAPIRateLimitReached)
	})

	t.Run("successfully fetch all logs across multiple pages with client.List", func(t *testing.T) {
		totalLogNumber := 220
		take := 1
		checkpoint := 0
		firstCall := true

		storedLogs := make([]*management.Log, totalLogNumber)
		for i := 0; i < totalLogNumber; i++ {
			var code = "f"
			var logID = fmt.Sprintf("log-%d", i)
			storedLogs[i] = &management.Log{LogID: &logID, Type: &code}
		}

		c := logClient{mgmt: &logManagementMock{
			ListFunc: func(ctx context.Context, opts ...management.RequestOption) ([]*management.Log, error) {
				var result []*management.Log
				if !firstCall {
					take = 100
				}

				if checkpoint >= totalLogNumber {
					return result, nil
				}

				if (checkpoint + take) >= totalLogNumber {
					result = storedLogs[checkpoint:totalLogNumber]
				} else {
					result = storedLogs[checkpoint:(checkpoint + take)]
				}

				checkpoint += take
				firstCall = false

				return result, nil
			},
		}}

		totalActualLogs, err := c.List(context.Background(), time.Now())
		require.NoError(t, err)
		assert.Len(t, totalActualLogs, totalLogNumber-1)
	})
}

func TestFindLatestCheckpoint(t *testing.T) {
	var checkpointID = "foo"
	t.Run("successfully find latest checkpoint, 2 days before auth0.from", func(t *testing.T) {
		from := time.Now()
		maxAttempts := 30
		expected := &management.Log{LogID: &checkpointID}
		var globalCounter = 2

		client := logClient{mgmt: &logManagementMock{
			ListFunc: func(ctx context.Context, opts ...management.RequestOption) ([]*management.Log, error) {
				var result []*management.Log
				if globalCounter == 2 {
					return append(result, &management.Log{LogID: &checkpointID}), nil
				}
				return result, nil
			},
		}}

		checkpoint, err := client.findLatestCheckpoint(context.TODO(), from, globalCounter, maxAttempts)
		require.NoError(t, err)
		assert.EqualValues(t, expected.LogID, checkpoint.LogID)
	})

	t.Run("fails to find latest checkpoint, max attempt are reached", func(t *testing.T) {
		from := time.Now()
		maxAttempts := 10
		var globalCounter = 12

		client := logClient{mgmt: &logManagementMock{
			ListFunc: func(ctx context.Context, opts ...management.RequestOption) ([]*management.Log, error) {
				var result []*management.Log
				if globalCounter == 2 {
					return append(result, &management.Log{LogID: &checkpointID}), nil
				}
				return result, nil
			},
		}}

		_, err := client.findLatestCheckpoint(context.TODO(), from, globalCounter, maxAttempts)
		require.Error(t, err)
		assert.ErrorIs(t, err, errLastCheckpointMaxAttemptsReached)
	})
}
