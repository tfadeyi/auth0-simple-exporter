package metrics_old

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"context"
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	failedAPIOperation     = "fapi"
	successfulAPIOperation = "sapi"

	TenantSuccessApiOperations = ctxKey("tenant_success_api_operations_total")
	TenantFailedApiOperations  = ctxKey("tenant_failed_api_operations_total")
)

func SapiOperationsMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(TenantSuccessApiOperations)),
			Help: "The number of successful API operations on the tenant. (codes: sapi)",
		}, []string{})
}

func FapiOperationsMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(TenantFailedApiOperations)),
			Help: "The number of failed API operations on the tenant. (codes: fapi)",
		}, []string{})
}

func ApiOperationsEventHandler(ctx context.Context, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}
	if ctx == nil {
		return errInvalidContext
	}
	sapi, ok := ctx.Value(TenantSuccessApiOperations).(*prometheus.CounterVec)
	if !ok {
		return errors.Annotate(errMissingLogEventMetric, "api operations metric is not in the context")
	}
	fapi, ok := ctx.Value(TenantFailedApiOperations).(*prometheus.CounterVec)
	if !ok {
		return errors.Annotate(errMissingLogEventMetric, "api operations metric is not in the context")
	}
	switch log.GetType() {
	case failedAPIOperation:
		fapi.WithLabelValues().Inc()
	case successfulAPIOperation:
		sapi.WithLabelValues().Inc()
	default:
		return errors.Annotate(errInvalidLogEvent, "api operations event handler can't handle event")
	}

	return nil
}
