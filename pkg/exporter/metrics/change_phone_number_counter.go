package metrics

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	failedChangePhoneNumber     = "fcpn"
	successfulChangePhoneNumber = "scpn"

	tenantSuccessChangePhoneNumber = "tenant_success_change_phone_number_total"
	tenantFailedChangePhoneNumber  = "tenant_failed_change_phone_number_total"
)

func successChangePhoneNumberMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(tenantSuccessChangePhoneNumber)),
			Help: "The number of successful change phone number operations on the tenant. (codes: scpn)",
		}, []string{})
}

func failChangePhoneNumberMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(tenantFailedChangePhoneNumber)),
			Help: "The number of failed change phone number operations on the tenant. (codes: fcpn)",
		}, []string{})
}

func changePhoneNumber(m *Metrics, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}
	switch log.GetType() {
	case failedChangePhoneNumber:
		increaseCounter(m.failedChangePhoneNumberCounter)
	case successfulChangePhoneNumber:
		increaseCounter(m.successfulChangePhoneNumberCounter)
	default:
		return errors.Annotate(errInvalidLogEvent, "change phone number operations event handler can't handle event")
	}

	return nil
}
