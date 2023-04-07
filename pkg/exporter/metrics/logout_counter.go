package metrics

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	failedLogoutCode     = "flo"
	successfulLogoutCode = "slo"

	tenantFailedLogoutOperations = "tenant_failed_logout_operations_total"
	tenantTotalLogoutOperations  = "tenant_logout_operations_total"
)

func logoutTotalCounterMetric(namespace, subsystem string, applicationClients []*management.Client) *prometheus.CounterVec {
	m := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, tenantTotalLogoutOperations),
			Help: "The total number of logout operations.",
		}, []string{"client"})
	for _, client := range applicationClients {
		m.WithLabelValues(client.GetName())
	}
	return m
}

func logoutFailCounterMetric(namespace, subsystem string, applicationClients []*management.Client) *prometheus.CounterVec {
	m := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, tenantFailedLogoutOperations),
			Help: "The number of failed logout operations. (codes: flo)",
		}, []string{"client"})
	for _, client := range applicationClients {
		m.WithLabelValues(client.GetName())
	}
	return m
}

func logout(m *Metrics, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}

	switch log.GetType() {
	case successfulLogoutCode:
		increaseCounter(m.logoutTotalCounter, log.GetClientName())
	case failedLogoutCode:
		increaseCounter(m.logoutFailCounter, log.GetClientName())
		increaseCounter(m.logoutTotalCounter, log.GetClientName())
	default:
		return errors.Annotate(errInvalidLogEvent, "logout event handler can't handle event")
	}

	return nil
}
