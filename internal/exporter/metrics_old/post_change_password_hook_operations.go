package metrics_old

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"context"
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	failedPostChangePasswordHookOperation     = "fcph"
	successfulPostChangePasswordHookOperation = "scph"

	TenantSuccessPostChangePasswordHookOperations = ctxKey("tenant_success_post_change_password_hook_operations_total")
	TenantFailedPostChangePasswordHookOperations  = ctxKey("tenant_failed_post_change_password_hook_operations_total")
)

func successPostChangePasswordHookOperationsMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(TenantSuccessChangeEmailOperations)),
			Help: "The number of successful post change user password hook operations on the tenant. (codes: scp)",
		}, []string{})
}

func failedPostChangePasswordHookOperationsMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(TenantFailedChangeEmailOperations)),
			Help: "The number of failed post change user password hook operations on the tenant. (codes: fcp)",
		}, []string{})
}

func postChangePasswordHookOperationsEventHandler(ctx context.Context, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}
	if ctx == nil {
		return errInvalidContext
	}
	scph, ok := ctx.Value(TenantSuccessPostChangePasswordHookOperations).(*prometheus.CounterVec)
	if !ok {
		return errors.Annotate(errMissingLogEventMetric, "post change user password hook operations metric is not in the context")
	}
	fcph, ok := ctx.Value(TenantFailedPostChangePasswordHookOperations).(*prometheus.CounterVec)
	if !ok {
		return errors.Annotate(errMissingLogEventMetric, "post change user password hook operations metric is not in the context")
	}
	switch log.GetType() {
	case failedPostChangePasswordHookOperation:
		fcph.WithLabelValues().Inc()
	case successfulPostChangePasswordHookOperation:
		scph.WithLabelValues().Inc()
	default:
		return errors.Annotate(errInvalidLogEvent, "post change user password hook operations event handler can't handle event")
	}

	return nil
}
