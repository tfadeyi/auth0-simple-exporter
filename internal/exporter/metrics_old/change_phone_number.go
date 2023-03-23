package metrics_old

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"context"
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	failedChangePhoneNumberOperation     = "fcpn"
	successfulChangePhoneNumberOperation = "scpn"

	TenantSuccessChangePhoneNumberOperations = ctxKey("tenant_success_change_phone_number_operations_total")
	TenantFailedChangePhoneNumberOperations  = ctxKey("tenant_failed_change_phone_number_operations_total")
)

func successChangePhoneNumberOperationsMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(TenantSuccessChangeEmailOperations)),
			Help: "The number of successful change phone number operations on the tenant. (codes: scpn)",
		}, []string{})
}

func failedChangePhoneNumberOperationsMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(TenantFailedChangeEmailOperations)),
			Help: "The number of failed change phone number operations on the tenant. (codes: fcpn)",
		}, []string{})
}

func changePhoneNumberOperationsEventHandler(ctx context.Context, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}
	if ctx == nil {
		return errInvalidContext
	}
	scpn, ok := ctx.Value(TenantSuccessChangePhoneNumberOperations).(*prometheus.CounterVec)
	if !ok {
		return errors.Annotate(errMissingLogEventMetric, "change phone number operations metric is not in the context")
	}
	fcpn, ok := ctx.Value(TenantFailedChangePhoneNumberOperations).(*prometheus.CounterVec)
	if !ok {
		return errors.Annotate(errMissingLogEventMetric, "change phone number operations metric is not in the context")
	}
	switch log.GetType() {
	case failedChangePhoneNumberOperation:
		fcpn.WithLabelValues().Inc()
	case successfulChangePhoneNumberOperation:
		scpn.WithLabelValues().Inc()
	default:
		return errors.Annotate(errInvalidLogEvent, "change phone number operations event handler can't handle event")
	}

	return nil
}
