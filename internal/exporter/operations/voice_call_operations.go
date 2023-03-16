package operations

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"context"
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	failedLoginCode                  = "f"
	successfulLoginCode              = "s"
	failedLoginWithIncorrectPassword = "fp"

	TenantLoginsOperations = CtxKey("tenant_login_operations_total")
)

func NewVoiceCallOperationsMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(TenantLoginsOperations)),
			Help: "The number of login operations. (codes: f,s,fp)",
		}, []string{"status", "type"})
}

func VoiceCallOperationsEventHandler(ctx context.Context, log *management.Log) error {
	if ctx == nil {
		return errInvalidContext
	}
	if log == nil {
		return errInvalidLogEvent
	}
	loginOperations, ok := ctx.Value(TenantLoginsOperations).(*prometheus.CounterVec)
	if !ok {
		return errors.Annotate(errMissingLogEventMetric, "login metric is not in the context")
	}
	switch log.GetType() {
	case failedLoginCode:
		loginOperations.WithLabelValues(failed, log.TypeName()).Inc()
	case successfulLoginCode:
		loginOperations.WithLabelValues(success, log.TypeName()).Inc()
	case failedLoginWithIncorrectPassword:
		loginOperations.WithLabelValues(failed, log.TypeName()).Inc()
	default:
		return errors.Annotate(errInvalidLogEvent, "login event handler can't handle event")
	}

	return nil
}
