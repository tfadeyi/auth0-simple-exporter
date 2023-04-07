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

	tenantTotalDeleteUser  = "tenant_delete_user_total"
	tenantFailedDeleteUser = "tenant_failed_delete_user_total"
)

func deleteUserTotalCounterMetric(namespace, subsystem string, applications []*management.Client) *prometheus.CounterVec {
	m := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, tenantTotalDeleteUser),
			Help: "The total number of delete user operations on the tenant. (codes: du,fdu)",
		}, []string{"client"})
	for _, client := range applications {
		m.WithLabelValues(client.GetName())
	}
	return m
}

func deleteUserFailCounterMetric(namespace, subsystem string, applications []*management.Client) *prometheus.CounterVec {
	m := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, tenantFailedDeleteUser),
			Help: "The number of failed delete user operations on the tenant. (codes: fdu)",
		}, []string{"client"})
	for _, client := range applications {
		m.WithLabelValues(client.GetName())
	}
	return m
}

func deleteUser(m *Metrics, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}

	switch log.GetType() {
	case successfulDeleteUserOperation:
		increaseCounter(m.deleteUserTotalCounter, log.GetClientName())
	case failedDeleteUserOperation:
		increaseCounter(m.deleteUserFailCounter, log.GetClientName())
		increaseCounter(m.deleteUserTotalCounter, log.GetClientName())
	default:
		return errors.Annotate(errInvalidLogEvent, "delete user operations event handler can't handle event")
	}

	return nil
}
