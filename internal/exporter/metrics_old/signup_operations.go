package metrics_old

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"context"

	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	failedSignupCode     = "fs"
	successfulSignupCode = "ss"

	TenantSignupOperations = ctxKey("tenant_sign_up_operations_total")
)

func NewSignupOperationsMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(TenantSignupOperations)),
			Help: "The number of signup operations. (codes: fs,ss)",
		}, []string{"status"})
}

func SignupOperationsEventHandler(ctx context.Context, log *management.Log) error {
	if ctx == nil {
		return errInvalidContext
	}
	if log == nil {
		return errInvalidLogEvent
	}
	op, ok := ctx.Value(TenantSignupOperations).(*prometheus.CounterVec)
	if !ok {
		return errors.Annotate(errMissingLogEventMetric, "signup metric is not in the context")
	}
	switch log.GetType() {
	case failedSignupCode:
		op.WithLabelValues(failed).Inc()
	case successfulSignupCode:
		op.WithLabelValues(success).Inc()
	default:
		return errors.Annotate(errInvalidLogEvent, "signup event handler can't handle event")
	}

	return nil
}
