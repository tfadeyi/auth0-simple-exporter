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

	tenantTotalAPIOperations  = "tenant_api_operations_total"
	tenantFailedAPIOperations = "tenant_failed_api_operations_total"
)

func APIOperationTotalCounterMetric(namespace, subsystem string, applications []*management.Client) *prometheus.CounterVec {
	m := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, tenantTotalAPIOperations),
			Help: "The total number of API operations on the tenant. (codes: sapi,fapi)",
		}, []string{"client"})
	for _, client := range applications {
		m.WithLabelValues(client.GetName())
	}
	return m
}

func APIOperationFailCounterMetric(namespace, subsystem string, applications []*management.Client) *prometheus.CounterVec {
	m := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, tenantFailedAPIOperations),
			Help: "The number of failed API operations on the tenant. (codes: fapi)",
		}, []string{"client"})
	for _, client := range applications {
		m.WithLabelValues(client.GetName())
	}
	return m
}

func apiOperations(m *Metrics, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}

	switch log.GetType() {
	case successAPIOperation:
		increaseCounter(m.apiOperationTotalCounter, log.GetClientName())
	case failAPIOperation:
		increaseCounter(m.apiOperationFailCounter, log.GetClientName())
		increaseCounter(m.apiOperationTotalCounter, log.GetClientName())
	default:
		return errors.Annotate(errInvalidLogEvent, "api operations event handler can't handle event")
	}

	return nil
}
