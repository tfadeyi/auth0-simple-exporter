package operations

import (
	"context"
	"github.com/auth0/go-auth0/management"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLoginOperations(t *testing.T) {
	t.Run("fail if context is not set", func(t *testing.T) {
		require.ErrorIs(t, LoginOperationsEventHandler(nil, &management.Log{}), errInvalidContext)
	})
	t.Run("fail if tenantLoginsOperations metric is missing from the context", func(t *testing.T) {
		require.ErrorIs(t, LoginOperationsEventHandler(context.TODO(), &management.Log{}), errMissingLogEventMetric)
	})
	t.Run("fails on nil log event", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, TenantLoginsOperations, NewApiOperationsMetric())
		require.ErrorIs(t, LoginOperationsEventHandler(ctx, nil), errInvalidLogEvent)
	})
	t.Run("fails on invalid log type", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, TenantLoginsOperations, NewApiOperationsMetric())
		var ty string
		ty = "random"
		require.ErrorIs(t, LoginOperationsEventHandler(ctx, &management.Log{Type: &ty}), errInvalidLogEvent)
	})
	t.Run("success on failedLoginCode log type", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, TenantLoginsOperations, NewApiOperationsMetric())
		var ty string
		ty = "f"
		require.NoError(t, LoginOperationsEventHandler(ctx, &management.Log{Type: &ty}))
	})
	t.Run("success on successfulLoginCode log type", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, TenantLoginsOperations, NewApiOperationsMetric())
		var ty string
		ty = "s"
		require.NoError(t, LoginOperationsEventHandler(ctx, &management.Log{Type: &ty}))
	})
	t.Run("success on failedLoginWithIncorrectPassword log type", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, TenantLoginsOperations, NewApiOperationsMetric())
		var ty string
		ty = "fp"
		require.NoError(t, LoginOperationsEventHandler(ctx, &management.Log{Type: &ty}))
	})
}
