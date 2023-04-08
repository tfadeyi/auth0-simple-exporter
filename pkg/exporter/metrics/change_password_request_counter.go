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

	tenantTotalChangePasswordRequest        = "tenant_change_password_request_total"
	tenantFailedLogoutChangePasswordRequest = "tenant_failed_change_password_request_total"
)

func changePasswordRequestTotalCounterMetric(namespace, subsystem string, applications []*management.Client) *prometheus.CounterVec {
	m := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, tenantTotalChangePasswordRequest),
			Help: "The total number of change_password_request operations. (codes: scpr,fcpr)",
		}, []string{"client"})
	for _, client := range applications {
		m.WithLabelValues(client.GetName())
	}
	return m
}

func changePasswordRequestFailCounterMetric(namespace, subsystem string, applications []*management.Client) *prometheus.CounterVec {
	m := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, tenantFailedLogoutChangePasswordRequest),
			Help: "The number of failed change password request operations. (codes: fcpr)",
		}, []string{"client"})
	for _, client := range applications {
		m.WithLabelValues(client.GetName())
	}
	return m
}

func changePasswordRequest(m *Metrics, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}

	switch log.GetType() {
	case successfulChangePasswordRequest:
		increaseCounter(m.changePasswordRequestTotalCounter, log.GetClientName())
	case failedChangePasswordRequest:
		increaseCounter(m.changePasswordRequestFailCounter, log.GetClientName())
		increaseCounter(m.changePasswordRequestTotalCounter, log.GetClientName())
	default:
		return errors.Annotate(errInvalidLogEvent, "change password request event handler can't handle event")
	}

	return nil
}
