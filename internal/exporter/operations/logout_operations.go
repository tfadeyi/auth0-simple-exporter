package operations

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

	TenantLogoutOperations = CtxKey("tenant_logout_operations_total")
)

func NewLogoutOperationsMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(TenantLogoutOperations)),
			Help: "The number of logout operations. (codes: slo, flo)",
		}, []string{"status"})
}
func LogoutOperationsEventHandler(ctx context.Context, log *management.Log) error {
	if ctx == nil {
		return errInvalidContext
	}
	if log == nil {
		return errInvalidLogEvent
	}
	op, ok := ctx.Value(TenantLogoutOperations).(*prometheus.CounterVec)
	if !ok {
		return errors.Annotate(errMissingLogEventMetric, "login metric is not in the context")
	}
	switch log.GetType() {
	case failedLogoutCode:
		op.WithLabelValues(failed).Inc()
	case successfulLogoutCode:
		op.WithLabelValues(success).Inc()
	default:
		return errors.Annotate(errInvalidLogEvent, "logout event handler can't handle event")
	}

	return nil
}
