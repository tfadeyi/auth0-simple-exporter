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

	tenantSuccessSendSMS = "tenant_success_send_sms_operations_total"
	tenantFailedSendSMS  = "tenant_failed_send_sms_operations_total"
)

func successSendSMSOperationsMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(tenantSuccessSendSMS)),
			Help: "The number of successful send_sms operations. (codes: gd_send_sms)",
		}, []string{})
}

func failSendSMSOperationsMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(tenantFailedSendSMS)),
			Help: "The number of failed send_sms operations. (codes: gd_send_sms_failure)",
		}, []string{})
}

func sendSMS(m *Metrics, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}
	switch log.GetType() {
	case successMfaSendSms:
		increaseCounter(m.successfulSendSMSCounter)
	case failureMfaSendFailure:
		increaseCounter(m.failedSendSMSCounter)
	default:
		return errors.Annotate(errInvalidLogEvent, "send_sms event handler can't handle event")
	}

	return nil
}
