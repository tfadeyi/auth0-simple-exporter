package operations

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"context"
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

type CtxKey string

const (
	failedAPIOperation     = "fapi"
	successfulAPIOperation = "sapi"

	TenantApiOperations = CtxKey("tenant_api_operations_total")
	success             = "success"
	failed              = "failed"
)

var (
	errInvalidContext        = errors.New("invalid or nil context object")
	errInvalidLogEvent       = errors.New("event handler doesn't accept the event log type")
	errMissingLogEventMetric = errors.New("couldn't find the prometheus metric for the required log event")
)

func NewApiOperationsMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(TenantApiOperations)),
			Help: "The number of API operations on the tenant.",
		}, []string{"status"})
}

func ApiOperationsEventHandler(ctx context.Context, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}
	if ctx == nil {
		return errInvalidContext
	}
	apiOperations, ok := ctx.Value(TenantApiOperations).(*prometheus.CounterVec)
	if !ok {
		return errors.Annotate(errMissingLogEventMetric, "api operations metric is not in the context")
	}

	switch log.GetType() {
	case failedAPIOperation:
		apiOperations.WithLabelValues(failed).Inc()
	case successfulAPIOperation:
		apiOperations.WithLabelValues(success).Inc()
	default:
		return errors.Annotate(errInvalidLogEvent, "api operations event handler can't handle event")
	}

	return nil
}
