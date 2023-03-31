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

	tenantSuccessfulLogoutOperations = "tenant_successful_logout_operations_total"
	tenantFailedLogoutOperations     = "tenant_failed_logout_operations_total"
)

func successLogoutCounterMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(tenantSuccessfulLogoutOperations)),
			Help: "The number of successful logout operations. (codes: slo)",
		}, []string{})
}

func failLogoutCounterMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(tenantFailedLogoutOperations)),
			Help: "The number of failed logout operations. (codes: flo)",
		}, []string{})
}

func logout(m *Metrics, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}

	switch log.GetType() {
	case failedLogoutCode:
		increaseCounter(m.failedLogoutCounter)
	case successfulLogoutCode:
		increaseCounter(m.successfulLogoutCounter)
	default:
		return errors.Annotate(errInvalidLogEvent, "logout event handler can't handle event")
	}

	return nil
}
