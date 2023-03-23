package metrics_old

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"context"
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// successMfaSendSms SMS for MFA successfully sent.
	successMfaSendSms = "gd_send_sms"
	// failureMfaSendFailure Attempt to send SMS for MFA failed.
	failureMfaSendFailure = "gd_send_sms_failure"

	TenantSendSMSOperations = ctxKey("tenant_send_sms_operations_total")
)

func NewSendSMSOperationsMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(TenantSendSMSOperations)),
			Help: "The number of send_sms operations. (codes: gd_send_sms, gd_send_sms_failure)",
		}, []string{"status"})
}

func SendSMSOperationsEventHandler(ctx context.Context, log *management.Log) error {
	if ctx == nil {
		return errInvalidContext
	}
	if log == nil {
		return errInvalidLogEvent
	}
	op, ok := ctx.Value(TenantSendSMSOperations).(*prometheus.CounterVec)
	if !ok {
		return errors.Annotate(errMissingLogEventMetric, "send_sms metric is not in the context")
	}
	switch log.GetType() {
	case successMfaSendSms:
		op.WithLabelValues(success).Inc()
	case failureMfaSendFailure:
		op.WithLabelValues(failed).Inc()
	default:
		return errors.Annotate(errInvalidLogEvent, "send_sms event handler can't handle event")
	}

	return nil
}
