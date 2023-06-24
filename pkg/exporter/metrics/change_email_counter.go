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

// @sloth.slo      name change_email_service_availability
// @sloth.slo      objective 99.0
// @sloth.sli      error_query sum(rate(tenant_failed_change_email_total[{{.window}}])) OR on() vector(0)
// @sloth.sli      total_query sum(rate(tenant_change_email_total[{{.window}}]))
// @sloth.slo      description SLO describing the availability of the Auth0 tenant change email service, setting the objective to 99%.
// @sloth.alerting name Auth0ChangeEmailAvailability
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
