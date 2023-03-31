package metrics

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// successMfaEmailSent Email for MFA successfully sent.
	successMfaEmailSent = "gd_send_email"

	TenantSendEmailOperations = "tenant_send_email_operations_total"
)

func sendEmailCounterMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(TenantSendEmailOperations)),
			Help: "The number of successful send email operations. (codes: gd_send_email)",
		}, []string{"type"})
}

func sendEmail(m *Metrics, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}

	switch log.GetType() {
	case successMfaEmailSent:
		increaseCounter(m.successfulSendEmailCounter, log.GetType())
	default:
		return errors.Annotate(errInvalidLogEvent, "send_email event handler can't handle event")
	}

	return nil
}
