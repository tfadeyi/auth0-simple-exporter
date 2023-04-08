package metrics

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	failedChangeEmail     = "fce"
	successfulChangeEmail = "sce"

	tenantTotalChangeEmail  = "tenant_change_email_total"
	tenantFailedChangeEmail = "tenant_failed_change_email_total"
)

func changeEmailTotalCounterMetric(namespace, subsystem string, applications []*management.Client) *prometheus.CounterVec {
	m := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, tenantTotalChangeEmail),
			Help: "The total number of change_user_email operations on the tenant. (codes: sce,fce)",
		}, []string{"client"})
	for _, client := range applications {
		m.WithLabelValues(client.GetName())
	}
	return m
}

func changeEmailFailCounterMetric(namespace, subsystem string, applications []*management.Client) *prometheus.CounterVec {
	m := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, tenantFailedChangeEmail),
			Help: "The number of failed change user email operations on the tenant. (codes: fce)",
		}, []string{"client"})
	for _, client := range applications {
		m.WithLabelValues(client.GetName())
	}
	return m
}

func changeEmail(m *Metrics, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}

	switch log.GetType() {
	case successfulChangeEmail:
		increaseCounter(m.changeEmailTotalCounter, log.GetClientName())
	case failedChangeEmail:
		increaseCounter(m.changeEmailFailCounter, log.GetClientName())
		increaseCounter(m.changeEmailTotalCounter, log.GetClientName())
	default:
		return errors.Annotate(errInvalidLogEvent, "change user email operations event handler can't handle event")
	}

	return nil
}
