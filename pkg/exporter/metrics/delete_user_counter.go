package metrics

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	failedDeleteUserOperation     = "fdu"
	successfulDeleteUserOperation = "du"

	tenantSuccessDeleteUser = "tenant_success_delete_user_total"
	tenantFailedDeleteUser  = "tenant_failed_delete_user_total"
)

func successDeleteUserCounterMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(tenantSuccessDeleteUser)),
			Help: "The number of successful delete user operations on the tenant. (codes: du)",
		}, []string{})
}

func failDeleteUserCounterMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(tenantFailedDeleteUser)),
			Help: "The number of failed delete user operations on the tenant. (codes: fdu)",
		}, []string{})
}

func deleteUser(m *Metrics, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}
	switch log.GetType() {
	case failedDeleteUserOperation:
		increaseCounter(m.failedDeleteUserCounter)
	case successfulDeleteUserOperation:
		increaseCounter(m.successfulDeleteUserCounter)
	default:
		return errors.Annotate(errInvalidLogEvent, "delete user operations event handler can't handle event")
	}

	return nil
}
