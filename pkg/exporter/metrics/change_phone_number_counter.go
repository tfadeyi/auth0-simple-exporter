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

	tenantTotalChangePhoneNumber  = "tenant_change_phone_number_total"
	tenantFailedChangePhoneNumber = "tenant_failed_change_phone_number_total"
)

func changePhoneNumberTotalCounterMetric(namespace, subsystem string, applications []*management.Client) *prometheus.CounterVec {
	m := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, tenantTotalChangePhoneNumber),
			Help: "The total number of change_phone_number operations on the tenant. (codes: scpn,fcpn)",
		}, []string{"client"})
	for _, client := range applications {
		m.WithLabelValues(client.GetName())
	}
	return m
}

func changePhoneNumberFailCounterMetric(namespace, subsystem string, applications []*management.Client) *prometheus.CounterVec {
	m := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, tenantFailedChangePhoneNumber),
			Help: "The number of failed change phone number operations on the tenant. (codes: fcpn)",
		}, []string{"client"})
	for _, client := range applications {
		m.WithLabelValues(client.GetName())
	}
	return m
}

func changePhoneNumber(m *Metrics, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}

	switch log.GetType() {
	case successfulChangePhoneNumber:
		increaseCounter(m.changePhoneNumberTotalCounter, log.GetClientName())
	case failedChangePhoneNumber:
		increaseCounter(m.changePhoneNumberFailCounter, log.GetClientName())
		increaseCounter(m.changePhoneNumberTotalCounter, log.GetClientName())
	default:
		return errors.Annotate(errInvalidLogEvent, "change phone number operations event handler can't handle event")
	}

	return nil
}
