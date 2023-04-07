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

	tenantTotalSignup  = "tenant_sign_up_operations_total"
	tenantFailedSignup = "tenant_failed_sign_up_operations_total"
)

func signupTotalCounterMetric(namespace, subsystem string, applications []*management.Client) *prometheus.CounterVec {
	m := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, tenantTotalSignup),
			Help: "The total number of signup operations.",
		}, []string{"client"})
	for _, client := range applications {
		initCounter(m, client.GetName())
	}
	return m
}

func signupFailCounterMetric(namespace, subsystem string, applications []*management.Client) *prometheus.CounterVec {
	m := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, tenantFailedSignup),
			Help: "The number of failed signup operations. (codes: fs)",
		}, []string{"client"})
	for _, client := range applications {
		initCounter(m, client.GetName())
	}
	return m
}

func signup(m *Metrics, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}

	switch log.GetType() {
	case failedSignupCode:
		increaseCounter(m.signupFailCounter, log.GetClientName())
		increaseCounter(m.signupTotalCounter, log.GetClientName())
	case successfulSignupCode:
		increaseCounter(m.signupTotalCounter, log.GetClientName())
	default:
		return errors.Annotate(errInvalidLogEvent, "signup event handler can't handle event")
	}

	return nil
}
