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

	tenantTotalPushNotification = "tenant_push_notification_total"
	tenantFailPushNotification  = "tenant_failed_push_notification_total"
)

func pushNotificationTotalCounterMetric(namespace, subsystem string, applications []*management.Client) *prometheus.CounterVec {
	m := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, tenantTotalPushNotification),
			Help: "The total number of push_notification operations. (codes: gd_send_pn,gd_send_pn_failure)",
		}, []string{"client"})
	for _, client := range applications {
		initCounter(m, client.GetName())
	}
	return m
}

func pushNotificationFailCounterMetric(namespace, subsystem string, applications []*management.Client) *prometheus.CounterVec {
	m := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, tenantFailPushNotification),
			Help: "The number of failed push_notification operations. (codes: gd_send_pn_failure)",
		}, []string{"client"})
	for _, client := range applications {
		initCounter(m, client.GetName())
	}
	return m
}

func pushNotification(m *Metrics, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}

	switch log.GetType() {
	case successfulMfaPushNotification:
		increaseCounter(m.pushNotificationTotalCounter, log.GetClientName())
	case failureMfaPushNotification:
		increaseCounter(m.pushNotificationFailCounter, log.GetClientName())
		increaseCounter(m.pushNotificationTotalCounter, log.GetClientName())
	default:
		return errors.Annotate(errInvalidLogEvent, "push_notification event handler can't handle event")
	}

	return nil
}
