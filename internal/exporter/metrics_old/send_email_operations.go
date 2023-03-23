package metrics_old

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"context"
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// successMfaEmailSent Email for MFA successfully sent.
	successMfaEmailSent = "gd_send_email"

	TenantSendEmailOperations = ctxKey("tenant_send_email_operations_total")
)

func NewSendEmailOperationsMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(TenantSendEmailOperations)),
			Help: "The number of successful MFA send email operations. (codes: gd_send_email)",
		}, []string{})
}

func SendEmailOperationsEventHandler(ctx context.Context, log *management.Log) error {
	if ctx == nil {
		return errInvalidContext
	}
	if log == nil {
		return errInvalidLogEvent
	}
	op, ok := ctx.Value(TenantSendEmailOperations).(*prometheus.CounterVec)
	if !ok {
		return errors.Annotate(errMissingLogEventMetric, "send_email metric is not in the context")
	}
	switch log.GetType() {
	case successMfaEmailSent:
		op.WithLabelValues().Inc()
	default:
		return errors.Annotate(errInvalidLogEvent, "send_email event handler can't handle event")
	}

	return nil
}
