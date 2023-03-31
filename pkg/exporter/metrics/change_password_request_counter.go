package metrics

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	failedChangePasswordRequest     = "fcpr"
	successfulChangePasswordRequest = "scpr"

	tenantSuccessfulChangePasswordRequest   = "tenant_successful_change_password_request_total"
	tenantFailedLogoutChangePasswordRequest = "tenant_failed_change_password_request_total"
)

func successChangePasswordRequestCounterMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(tenantSuccessfulChangePasswordRequest)),
			Help: "The number of successful change password request operations. (codes: scpr)",
		}, []string{})
}

func failChangePasswordRequestCounterMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(tenantFailedLogoutChangePasswordRequest)),
			Help: "The number of failed change password request operations. (codes: fcpr)",
		}, []string{})
}

func changePasswordRequest(m *Metrics, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}

	switch log.GetType() {
	case failedChangePasswordRequest:
		increaseCounter(m.failedChangePasswordRequestCounter)
	case successfulChangePasswordRequest:
		increaseCounter(m.successfulChangePasswordRequestCounter)
	default:
		return errors.Annotate(errInvalidLogEvent, "change password request event handler can't handle event")
	}

	return nil
}
