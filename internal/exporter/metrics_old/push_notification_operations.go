package metrics_old

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"context"
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// successMfaPushNotificationSuccess Push notification for MFA sent successfully sent.
	successMfaPushNotificationSuccess = "gd_send_pn"
	// failureMfaPushNotification Push notification for MFA failed.
	failureMfaPushNotification = "gd_send_pn_failure"

	TenantPushNotificationOperations = ctxKey("tenant_push_notification_operations_total")
)

func NewPushNotificationOperationsMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(TenantPushNotificationOperations)),
			Help: "The number of push_notification operations. (codes: gd_send_pn, gd_send_pn_failure)",
		}, []string{"status", "type"})
}

func PushNotificationOperationsEventHandler(ctx context.Context, log *management.Log) error {
	if ctx == nil {
		return errInvalidContext
	}
	if log == nil {
		return errInvalidLogEvent
	}
	op, ok := ctx.Value(TenantPushNotificationOperations).(*prometheus.CounterVec)
	if !ok {
		return errors.Annotate(errMissingLogEventMetric, "push_notification metric is not in the context")
	}
	switch log.GetType() {
	case failureMfaPushNotification:
		op.WithLabelValues(failed, log.GetType()).Inc()
	case successMfaPushNotificationSuccess:
		op.WithLabelValues(success, log.GetType()).Inc()
	default:
		return errors.Annotate(errInvalidLogEvent, "push_notification event handler can't handle event")
	}

	return nil
}
