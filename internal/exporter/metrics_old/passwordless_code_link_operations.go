package metrics_old

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"context"
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// successSendCodeLink Passwordless login code/link has been sent.
	successSendCodeLink = "cls"
	// successSendCode Passwordless login code has been sent.
	successSendCode = "cs"

	TenantSendCodeLinkOperations = ctxKey("tenant_send_code_link_operations_total")
)

func NewSendCodeLinkOperationsMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(TenantSendCodeLinkOperations)),
			Help: "The number of send_code_link operations. (codes: cls,cs)",
		}, []string{"type"})
}

func SendCodeLinkOperationsEventHandler(ctx context.Context, log *management.Log) error {
	if ctx == nil {
		return errInvalidContext
	}
	if log == nil {
		return errInvalidLogEvent
	}
	op, ok := ctx.Value(TenantSendCodeLinkOperations).(*prometheus.CounterVec)
	if !ok {
		return errors.Annotate(errMissingLogEventMetric, "code_link metric is not in the context")
	}
	switch log.GetType() {
	case successSendCodeLink:
		op.WithLabelValues(log.GetType()).Inc()
	case successSendCode:
		op.WithLabelValues(log.GetType()).Inc()
	default:
		return errors.Annotate(errInvalidLogEvent, "code_link event handler can't handle event")
	}

	return nil
}
