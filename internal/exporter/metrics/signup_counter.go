package metrics

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	failedSignupCode     = "fs"
	successfulSignupCode = "ss"

	tenantSuccessSignup = "tenant_successful_sign_up_total"
	tenantFailedSignup = "tenant_failed_sign_up_total"
)

func successSignupCounterMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(tenantSuccessSignup)),
			Help: "The number of successful signup operations. (codes: ss)",
		}, []string{})
}

func failSignupCounterMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(tenantFailedSignup)),
			Help: "The number of failed signup operations. (codes: fs)",
		}, []string{})
}

func signup(m *Metrics, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}

	switch log.GetType() {
	case failedSignupCode:
		increaseCounter(m.failedSignupCounter)
	case successfulSignupCode:
		increaseCounter(m.successfulSignupCounter)
	default:
		return errors.Annotate(errInvalidLogEvent, "signup event handler can't handle event")
	}

	return nil
}
