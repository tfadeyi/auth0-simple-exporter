package metrics_old

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"context"
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// successMfaSendVoiceCall Voice call for MFA successfully made.
	successMfaSendVoiceCall = "gd_send_voice"
	// failedMfaSendVoiceCall Attempt to make Voice call for MFA failed.
	failedMfaSendVoiceCall = "gd_send_voice_failure"

	TenantSendVoiceCallOperations = ctxKey("tenant_send_voice_call_operations_total")
)

func NewSendVoiceCallOperationsMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(TenantSendVoiceCallOperations)),
			Help: "The number of voice_call operations. (codes: gd_send_voice, gd_send_voice_failure)",
		}, []string{"status"})
}

func SendVoiceCallOperationsEventHandler(ctx context.Context, log *management.Log) error {
	if ctx == nil {
		return errInvalidContext
	}
	if log == nil {
		return errInvalidLogEvent
	}
	op, ok := ctx.Value(TenantSendVoiceCallOperations).(*prometheus.CounterVec)
	if !ok {
		return errors.Annotate(errMissingLogEventMetric, "send_voice_call metric is not in the context")
	}
	switch log.GetType() {
	case failedMfaSendVoiceCall:
		op.WithLabelValues(failed).Inc()
	case successMfaSendVoiceCall:
		op.WithLabelValues(success).Inc()
	default:
		return errors.Annotate(errInvalidLogEvent, "send_voice_call event handler can't handle event")
	}

	return nil
}
