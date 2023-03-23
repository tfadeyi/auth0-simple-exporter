package metrics

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	failedLogin                      = "f"
	successfulLogin                  = "s"
	failedLoginWithIncorrectPassword = "fp"
	failedLoginWithIncorrectUsername = "fu"

	tenantSuccessfulLogin = "tenant_successful_login_operations_total"
	tenantFailedLogin     = "tenant_failed_login_operations_total"
)

func successLoginCounterMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, tenantSuccessfulLogin),
			Help: "The number of successful login operations. (codes: s)",
		}, []string{})
}
func failLoginCounterMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, tenantFailedLogin),
			Help: "The number of failed login operations. (codes: f,fp,fu)",
		}, []string{"type"})
}

func login(m *Metrics, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}
	switch log.GetType() {
	case successfulLogin:
		increaseCounter(m.successfulLoginCnt)
	case failedLogin:
		increaseCounter(m.failedLoginCnt, log.GetType())
	case failedLoginWithIncorrectPassword:
		increaseCounter(m.failedLoginCnt, log.GetType())
	case failedLoginWithIncorrectUsername:
		increaseCounter(m.failedLoginCnt, log.GetType())
	default:
		return errors.Annotate(errInvalidLogEvent, "login event handler can't handle event")
	}
	return nil
}
