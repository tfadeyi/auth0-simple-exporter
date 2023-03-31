package metrics

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	failChangePassword    = "fcp"
	successChangePassword = "sce"

	tenantSuccessChangePassword = "tenant_success_change_password_total"
	tenantFailedChangePassword  = "tenant_failed_change_password_total"
)

func successChangePasswordMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(tenantSuccessChangePassword)),
			Help: "The number of successful change user password operations on the tenant. (codes: scp)",
		}, []string{})
}

func failChangePasswordCounterMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(tenantFailedChangePassword)),
			Help: "The number of failed change user password operations on the tenant. (codes: fcp)",
		}, []string{})
}

func changePassword(m *Metrics, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}
	switch log.GetType() {
	case failChangePassword:
		increaseCounter(m.successfulLoginCnt)
	case successChangePassword:
		increaseCounter(m.successfulLoginCnt)
	default:
		return errors.Annotate(errInvalidLogEvent, "change user password operations event handler can't handle event")
	}

	return nil
}
