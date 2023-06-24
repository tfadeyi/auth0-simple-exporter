package metrics

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	failedChangePassword     = "fcp"
	successfulChangePassword = "sce"

	tenantTotalChangePassword  = "tenant_change_password_total"
	tenantFailedChangePassword = "tenant_failed_change_password_total"
)

// @sloth.slo      name change_password_service_availability
// @sloth.slo      objective 99.0
// @sloth.sli      error_query sum(rate(tenant_failed_change_password_total[{{.window}}])) OR on() vector(0)
// @sloth.sli      total_query sum(rate(tenant_change_password_total[{{.window}}]))
// @sloth.slo      description SLO describing the availability of the Auth0 tenant change password service, setting the objective to 99%.
// @sloth.alerting name Auth0ChangePasswordAvailability
func changePasswordTotalCounterMetric(namespace, subsystem string, applications []*management.Client) *prometheus.CounterVec {
	m := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, tenantTotalChangePassword),
			Help: "The total number of change_user_password operations on the tenant. (codes: scp,fcp)",
		}, []string{"client"})
	for _, client := range applications {
		m.WithLabelValues(client.GetName())
	}
	return m
}

func changePasswordFailCounterMetric(namespace, subsystem string, applications []*management.Client) *prometheus.CounterVec {
	m := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, tenantFailedChangePassword),
			Help: "The number of failed change_user_password operations on the tenant. (codes: fcp)",
		}, []string{"client"})
	for _, client := range applications {
		m.WithLabelValues(client.GetName())
	}
	return m
}

func changePassword(m *Metrics, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}

	switch log.GetType() {
	case successfulChangePassword:
		increaseCounter(m.changePasswordTotalCounter, log.GetClientName())
	case failedChangePassword:
		increaseCounter(m.changePasswordFailCounter, log.GetClientName())
		increaseCounter(m.changePasswordTotalCounter, log.GetClientName())
	default:
		return errors.Annotate(errInvalidLogEvent, "change user password operations event handler can't handle event")
	}

	return nil
}
