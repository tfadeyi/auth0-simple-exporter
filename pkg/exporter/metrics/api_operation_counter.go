package metrics

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	failAPIOperation    = "fapi"
	successAPIOperation = "sapi"

	tenantSuccessAPIOperations = "tenant_success_api_operations_total"
	tenantFailedAPIOperations  = "tenant_failed_api_operations_total"
)

func successAPIOperationCounterMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(tenantSuccessAPIOperations)),
			Help: "The number of successful API operations on the tenant. (codes: sapi)",
		}, []string{})
}

func failAPIOperationCounterMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(tenantFailedAPIOperations)),
			Help: "The number of failed API operations on the tenant. (codes: fapi)",
		}, []string{})
}

func apiOperations(m *Metrics, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}

	switch log.GetType() {
	case failAPIOperation:
		increaseCounter(m.failedAPIOperationCounter)
	case successAPIOperation:
		increaseCounter(m.successfulAPIOperationCounter)
	default:
		return errors.Annotate(errInvalidLogEvent, "api operations event handler can't handle event")
	}

	return nil
}
