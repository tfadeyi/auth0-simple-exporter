package metrics

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// successMfaSendVoiceCall Voice call for MFA successfully made.
	successMfaSendVoiceCall = "gd_send_voice"
	// failedMfaSendVoiceCall Attempt to make Voice call for MFA failed.
	failedMfaSendVoiceCall = "gd_send_voice_failure"

	tenantSuccessSendVoiceCall = "tenant_success_send_voice_call_operations_total"
	tenantFailedSendVoiceCall = "tenant_failed_send_voice_call_operations_total"
)

func successSendVoiceCallCounterMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(tenantSuccessSendVoiceCall)),
			Help: "The number of voice_call operations. (codes: gd_send_voice)",
		}, []string{})
}

func failSendVoiceCallCounterMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(tenantFailedSendVoiceCall)),
			Help: "The number of voice_call operations. (codes: gd_send_voice_failure)",
		}, []string{})
}

func sendVoiceCall(m *Metrics, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}

	switch log.GetType() {
	case failedMfaSendVoiceCall:
		increaseCounter(m.failedVoiceCallCounter)
	case successMfaSendVoiceCall:
		increaseCounter(m.successfulVoiceCallCounter)
	default:
		return errors.Annotate(errInvalidLogEvent, "send_voice_call event handler can't handle event")
	}

	return nil
}
