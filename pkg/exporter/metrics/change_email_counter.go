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

	tenantSuccessChangeEmail = "tenant_success_change_email_total"
	tenantFailedChangeEmail  = "tenant_failed_change_email_total"
)

func successChangeEmailCounterMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(tenantSuccessChangeEmail)),
			Help: "The number of successful change user email operations on the tenant. (codes: sce)",
		}, []string{})
}

func failedChangeEmailCounterMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(tenantFailedChangeEmail)),
			Help: "The number of failed change user email operations on the tenant. (codes: fce)",
		}, []string{})
}

func changeEmail(m *Metrics, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}

	switch log.GetType() {
	case failedChangeEmail:
		increaseCounter(m.failedChangeEmailCounter)
	case successfulChangeEmail:
		increaseCounter(m.successfulChangeEmailCounter)
	default:
		return errors.Annotate(errInvalidLogEvent, "change user email operations event handler can't handle event")
	}

	return nil
}
