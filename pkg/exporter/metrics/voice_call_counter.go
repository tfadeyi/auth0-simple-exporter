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

	tenantTotalSendVoiceCall  = "tenant_send_voice_call_operations_total"
	tenantFailedSendVoiceCall = "tenant_failed_send_voice_call_operations_total"
)

func sendVoiceCallTotalCounterMetric(namespace, subsystem string, applications []*management.Client) *prometheus.CounterVec {
	m := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, tenantTotalSendVoiceCall),
			Help: "The total number of voice_call operations.",
		}, []string{"client"})
	for _, client := range applications {
		initCounter(m, client.GetName())
	}
	return m
}

func sendVoiceCallFailCounterMetric(namespace, subsystem string, applications []*management.Client) *prometheus.CounterVec {
	m := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, tenantFailedSendVoiceCall),
			Help: "The number of failed voice_call operations. (codes: gd_send_voice_failure)",
		}, []string{"client"})
	for _, client := range applications {
		initCounter(m, client.GetName())
	}
	return m
}

func sendVoiceCall(m *Metrics, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}

	switch log.GetType() {
	case failedMfaSendVoiceCall:
		increaseCounter(m.voiceCallFailCounter, log.GetClientName())
		increaseCounter(m.voiceCallTotalCounter, log.GetClientName())
	case successMfaSendVoiceCall:
		increaseCounter(m.voiceCallTotalCounter, log.GetClientName())
	default:
		return errors.Annotate(errInvalidLogEvent, "send_voice_call event handler can't handle event")
	}

	return nil
}
