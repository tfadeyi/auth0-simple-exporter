package operations

import (
	"context"
	"testing"

	"github.com/auth0/go-auth0/management"
	"github.com/stretchr/testify/require"
)

func TestApiOperations(t *testing.T) {
	t.Run("fail if TenantApiOperations metric is missing from the context", func(t *testing.T) {
		require.ErrorIs(t, ApiOperationsEventHandler(context.TODO(), &management.Log{}), errMissingLogEventMetric)
	})
	t.Run("fails on nil log event", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, TenantApiOperations, NewApiOperationsMetric("", ""))
		require.ErrorIs(t, ApiOperationsEventHandler(ctx, nil), errInvalidLogEvent)
	})
	t.Run("fails on invalid log type", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, TenantApiOperations, NewApiOperationsMetric("", ""))
		var ty = "random"
		require.ErrorIs(t, ApiOperationsEventHandler(ctx, &management.Log{Type: &ty}), errInvalidLogEvent)
	})
	t.Run("success on failedAPIOperation log type", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, TenantApiOperations, NewApiOperationsMetric("", ""))
		var ty = "fapi"
		require.NoError(t, ApiOperationsEventHandler(ctx, &management.Log{Type: &ty}))
	})
	t.Run("success on successfulAPIOperation log type", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, TenantApiOperations, NewApiOperationsMetric("", ""))
		var ty = "sapi"
		require.NoError(t, ApiOperationsEventHandler(ctx, &management.Log{Type: &ty}))
	})
}
