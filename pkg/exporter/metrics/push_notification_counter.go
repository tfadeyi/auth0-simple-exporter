package metrics

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// successfulMfaPushNotification Push notification for MFA sent successfully sent.
	successfulMfaPushNotification = "gd_send_pn"
	// failureMfaPushNotification Push notification for MFA failed.
	failureMfaPushNotification = "gd_send_pn_failure"

	tenantSuccessPushNotification = "tenant_success_push_notification_total"
	tenantFailPushNotification    = "tenant_fail_push_notification_total"
)

func successPushNotificationCounterMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(tenantSuccessPushNotification)),
			Help: "The number of successful push_notification operations. (codes: gd_send_pn)",
		}, []string{})
}

func failPushNotificationCounterMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(tenantFailPushNotification)),
			Help: "The number of failed push_notification operations. (codes: gd_send_pn_failure)",
		}, []string{})
}

func pushNotification(m *Metrics, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}

	switch log.GetType() {
	case failureMfaPushNotification:
		increaseCounter(m.failedPushNotificationCounter)
	case successfulMfaPushNotification:
		increaseCounter(m.successfulPushNotificationCounter)
	default:
		return errors.Annotate(errInvalidLogEvent, "push_notification event handler can't handle event")
	}

	return nil
}
