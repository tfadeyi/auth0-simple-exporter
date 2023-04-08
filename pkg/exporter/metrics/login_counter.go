package metrics

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	successfulLogin                  = "s"
	failedLogin                      = "f"
	failedLoginWithIncorrectPassword = "fp"
	failedLoginWithIncorrectUsername = "fu"

	tenantFailedLogin = "tenant_failed_login_operations_total"
	tenantTotalLogin  = "tenant_login_operations_total"
)

func loginTotalCounterMetric(namespace, subsystem string, applicationClients []*management.Client) *prometheus.CounterVec {
	m := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, tenantTotalLogin),
			Help: "The total number of login operations.",
		}, []string{"client"})
	for _, client := range applicationClients {
		m.WithLabelValues(client.GetName())
	}
	return m
}
func loginFailCounterMetric(namespace, subsystem string, applicationClients []*management.Client) *prometheus.CounterVec {
	m := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, tenantFailedLogin),
			Help: "The number of failed login operations. (codes: f,fp,fu)",
		}, []string{"client"})
	for _, client := range applicationClients {
		m.WithLabelValues(client.GetName())
	}
	return m
}

func login(m *Metrics, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}

	switch log.GetType() {
	case successfulLogin:
		increaseCounter(m.loginTotalCounter, log.GetClientName())
	case failedLogin:
		increaseCounter(m.loginFailCounter, log.GetClientName())
		increaseCounter(m.loginTotalCounter, log.GetClientName())
	case failedLoginWithIncorrectPassword:
		increaseCounter(m.loginFailCounter, log.GetClientName())
		increaseCounter(m.loginTotalCounter, log.GetClientName())
	case failedLoginWithIncorrectUsername:
		increaseCounter(m.loginFailCounter, log.GetClientName())
		increaseCounter(m.loginTotalCounter, log.GetClientName())
	default:
		return errors.Annotate(errInvalidLogEvent, "login event handler can't handle event")
	}
	return nil
}
