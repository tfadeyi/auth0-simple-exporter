package metrics_old

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"context"
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	failedChangePasswordOperation     = "fcp"
	successfulChangePasswordOperation = "sce"

	TenantSuccessChangePasswordOperations = ctxKey("tenant_success_change_password_operations_total")
	TenantFailedChangePasswordOperations  = ctxKey("tenant_failed_change_password_operations_total")
)

func successChangePasswordOperationsMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(TenantSuccessChangeEmailOperations)),
			Help: "The number of successful change user password operations on the tenant. (codes: scp)",
		}, []string{})
}

func failedChangePasswordOperationsMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(TenantFailedChangeEmailOperations)),
			Help: "The number of failed change user password operations on the tenant. (codes: fcp)",
		}, []string{})
}

func changePasswordOperationsEventHandler(ctx context.Context, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}
	if ctx == nil {
		return errInvalidContext
	}
	scp, ok := ctx.Value(TenantSuccessChangePasswordOperations).(*prometheus.CounterVec)
	if !ok {
		return errors.Annotate(errMissingLogEventMetric, "change user password operations metric is not in the context")
	}
	fcp, ok := ctx.Value(TenantFailedChangePasswordOperations).(*prometheus.CounterVec)
	if !ok {
		return errors.Annotate(errMissingLogEventMetric, "change user password operations metric is not in the context")
	}
	switch log.GetType() {
	case failedChangePasswordOperation:
		fcp.WithLabelValues().Inc()
	case successfulChangePasswordOperation:
		scp.WithLabelValues().Inc()
	default:
		return errors.Annotate(errInvalidLogEvent, "change user password operations event handler can't handle event")
	}

	return nil
}
