package metrics

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// successMfaSendSms SMS for MFA successfully sent.
	successMfaSendSms = "gd_send_sms"
	// failureMfaSendFailure Attempt to send SMS for MFA failed.
	failureMfaSendFailure = "gd_send_sms_failure"

	tenantTotalSendSMS  = "tenant_send_sms_operations_total"
	tenantFailedSendSMS = "tenant_failed_send_sms_operations_total"
)

func sendSMSTotalCounterMetric(namespace, subsystem string, applications []*management.Client) *prometheus.CounterVec {
	m := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, tenantTotalSendSMS),
			Help: "The total number of successful send_sms operations.",
		}, []string{"client"})
	for _, client := range applications {
		initCounter(m, client.GetName())
	}
	return m
}

func sendSMSFailCounterMetric(namespace, subsystem string, applications []*management.Client) *prometheus.CounterVec {
	m := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, tenantFailedSendSMS),
			Help: "The number of failed send_sms operations. (codes: gd_send_sms_failure)",
		}, []string{"client"})
	for _, client := range applications {
		initCounter(m, client.GetName())
	}
	return m
}

func sendSMS(m *Metrics, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}
	switch log.GetType() {
	case successMfaSendSms:
		increaseCounter(m.sendSMSTotalCounter, log.GetClientName())
	case failureMfaSendFailure:
		increaseCounter(m.sendSMSFailCounter, log.GetClientName())
		increaseCounter(m.sendSMSTotalCounter, log.GetClientName())
	default:
		return errors.Annotate(errInvalidLogEvent, "send_sms event handler can't handle event")
	}

	return nil
}
