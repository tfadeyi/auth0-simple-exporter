package metrics_old

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"context"
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	failedDeleteUserOperation     = "fdu"
	successfulDeleteUserOperation = "du"

	TenantSuccessDeleteUserOperations = ctxKey("tenant_success_delete_user_operations_total")
	TenantFailedDeleteUserOperations  = ctxKey("tenant_failed_delete_user_operations_total")
)

func successDeleteUserOperationsMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(TenantSuccessDeleteUserOperations)),
			Help: "The number of successful delete user operations on the tenant. (codes: du)",
		}, []string{})
}

func failedDeleteUserOperationsMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(TenantFailedDeleteUserOperations)),
			Help: "The number of failed delete user operations on the tenant. (codes: fdu)",
		}, []string{})
}

func deleteUserOperationsEventHandler(ctx context.Context, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}
	if ctx == nil {
		return errInvalidContext
	}
	du, ok := ctx.Value(TenantSuccessDeleteUserOperations).(*prometheus.CounterVec)
	if !ok {
		return errors.Annotate(errMissingLogEventMetric, "delete user operations metric is not in the context")
	}
	fdu, ok := ctx.Value(TenantFailedDeleteUserOperations).(*prometheus.CounterVec)
	if !ok {
		return errors.Annotate(errMissingLogEventMetric, "delete user operations metric is not in the context")
	}
	switch log.GetType() {
	case failedDeleteUserOperation:
		fdu.WithLabelValues().Inc()
	case successfulDeleteUserOperation:
		du.WithLabelValues().Inc()
	default:
		return errors.Annotate(errInvalidLogEvent, "delete user operations event handler can't handle event")
	}

	return nil
}
