package metrics_old

import (
	"context"
	"github.com/auth0/go-auth0/management"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLoginOperations(t *testing.T) {
	t.Run("fail if tenantLoginsOperations metric is missing from the context", func(t *testing.T) {
		require.ErrorIs(t, FailedLoginOperationsEventHandler(context.TODO(), &management.Log{}), errMissingLogEventMetric)
	})
	t.Run("fails on nil log event", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, TenantSendSMSOperations, SapiOperationsMetric("", ""))
		require.ErrorIs(t, FailedLoginOperationsEventHandler(ctx, nil), errInvalidLogEvent)
	})
	t.Run("fails on invalid log type", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, TenantSendSMSOperations, SapiOperationsMetric("", ""))
		var ty = "random"
		require.ErrorIs(t, FailedLoginOperationsEventHandler(ctx, &management.Log{Type: &ty}), errInvalidLogEvent)
	})
	t.Run("success on failedLogin log type", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, TenantSendSMSOperations, SapiOperationsMetric("", ""))
		var ty = "f"
		require.NoError(t, FailedLoginOperationsEventHandler(ctx, &management.Log{Type: &ty}))
	})
	t.Run("success on successfulLogin log type", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, TenantSendSMSOperations, SapiOperationsMetric("", ""))
		var ty = "s"
		require.NoError(t, FailedLoginOperationsEventHandler(ctx, &management.Log{Type: &ty}))
	})
	t.Run("success on failedLoginWithIncorrectPassword log type", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, TenantSendSMSOperations, SapiOperationsMetric("", ""))
		var ty = "fp"
		require.NoError(t, FailedLoginOperationsEventHandler(ctx, &management.Log{Type: &ty}))
	})
}
