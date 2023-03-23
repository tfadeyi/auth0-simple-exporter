package metrics_old

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"context"
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	failedChangeEmailOperation     = "fce"
	successfulChangeEmailOperation = "sce"

	TenantSuccessChangeEmailOperations = ctxKey("tenant_success_change_email_operations_total")
	TenantFailedChangeEmailOperations  = ctxKey("tenant_failed_change_email_operations_total")
)

func successChangeEmailOperationsMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(TenantSuccessChangeEmailOperations)),
			Help: "The number of successful change user email operations on the tenant. (codes: sce)",
		}, []string{})
}

func failedChangeEmailOperationsMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(TenantFailedChangeEmailOperations)),
			Help: "The number of failed change user email operations on the tenant. (codes: fce)",
		}, []string{})
}

func changeEmailOperationsEventHandler(ctx context.Context, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}
	if ctx == nil {
		return errInvalidContext
	}
	sce, ok := ctx.Value(TenantSuccessChangeEmailOperations).(*prometheus.CounterVec)
	if !ok {
		return errors.Annotate(errMissingLogEventMetric, "change user email operations metric is not in the context")
	}
	fce, ok := ctx.Value(TenantFailedChangeEmailOperations).(*prometheus.CounterVec)
	if !ok {
		return errors.Annotate(errMissingLogEventMetric, "change user email operations metric is not in the context")
	}
	switch log.GetType() {
	case failedChangeEmailOperation:
		fce.WithLabelValues().Inc()
	case successfulChangeEmailOperation:
		sce.WithLabelValues().Inc()
	default:
		return errors.Annotate(errInvalidLogEvent, "change user email operations event handler can't handle event")
	}

	return nil
}
