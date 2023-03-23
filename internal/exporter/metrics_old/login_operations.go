package metrics_old

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"context"

	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	failedLogin                      = "f"
	successfulLogin                  = "s"
	failedLoginWithIncorrectPassword = "fp"
	failedLoginWithIncorrectUsername = "fu"

	TenantSuccessfulLoginsOperations = ctxKey("tenant_successful_login_operations_total")
	TenantFailedLoginsOperations     = ctxKey("tenant_failed_login_operations_total")
)

func NewSuccessfulLoginOperationsMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(TenantSuccessfulLoginsOperations)),
			Help: "The number of successful login operations. (codes: s)",
		}, []string{})
}
func NewFailedLoginOperationsMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(TenantFailedLoginsOperations)),
			Help: "The number of failed login operations. (codes: f,fp,fu)",
		}, []string{"type"})
}

func LoginOperationsEventHandler(ctx context.Context, log *management.Log) error {
	if ctx == nil {
		return errInvalidContext
	}
	if log == nil {
		return errInvalidLogEvent
	}
	success, ok := ctx.Value(TenantSuccessfulLoginsOperations).(*prometheus.CounterVec)
	if !ok {
		return errors.Annotate(errMissingLogEventMetric, "login metric is not in the context")
	}
	failed, ok := ctx.Value(TenantFailedLoginsOperations).(*prometheus.CounterVec)
	if !ok {
		return errors.Annotate(errMissingLogEventMetric, "login metric is not in the context")
	}
	switch log.GetType() {
	case successfulLogin:
		success.WithLabelValues().Inc()
	case failedLogin:
		failed.WithLabelValues(log.GetType()).Inc()
	case failedLoginWithIncorrectPassword:
		failed.WithLabelValues(log.GetType()).Inc()
	case failedLoginWithIncorrectUsername:
		failed.WithLabelValues(log.GetType()).Inc()
	default:
		return errors.Annotate(errInvalidLogEvent, "login event handler can't handle event")
	}
	return nil
}
