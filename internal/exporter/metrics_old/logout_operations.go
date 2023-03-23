package metrics_old

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"context"
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	failedLogoutCode     = "flo"
	successfulLogoutCode = "slo"

	TenantSuccessLogoutOperations = ctxKey("tenant_successful_logout_operations_total")
	TenantFailedLogoutOperations  = ctxKey("tenant_failed_logout_operations_total")
)

func NewSuccessLogoutOperationsMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(TenantSuccessLogoutOperations)),
			Help: "The number of successful logout operations. (codes: slo)",
		}, []string{})
}

func NewFailedLogoutOperationsMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(TenantFailedLogoutOperations)),
			Help: "The number of failed logout operations. (codes: flo)",
		}, []string{})
}

func LogoutOperationsEventHandler(ctx context.Context, log *management.Log) error {
	if ctx == nil {
		return errInvalidContext
	}
	if log == nil {
		return errInvalidLogEvent
	}
	slo, ok := ctx.Value(TenantSuccessLogoutOperations).(*prometheus.CounterVec)
	if !ok {
		return errors.Annotate(errMissingLogEventMetric, "logout metric is not in the context")
	}
	flo, ok := ctx.Value(TenantFailedLogoutOperations).(*prometheus.CounterVec)
	if !ok {
		return errors.Annotate(errMissingLogEventMetric, "logout metric is not in the context")
	}
	switch log.GetType() {
	case failedLogoutCode:
		flo.WithLabelValues().Inc()
	case successfulLogoutCode:
		slo.WithLabelValues().Inc()
	default:
		return errors.Annotate(errInvalidLogEvent, "logout event handler can't handle event")
	}

	return nil
}
