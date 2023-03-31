package metrics

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// successSendCodeLink Passwordless login code/link has been sent.
	successSendCodeLink = "cls"
	// successSendCode Passwordless login code has been sent.
	successSendCode = "cs"

	tenantSendCodeLink = "tenant_send_code_link_total"
)

func passwordLessSendCodeLinkCounterMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(tenantSendCodeLink)),
			Help: "The number of send_code_link operations. (codes: cls,cs)",
		}, []string{"type"})
}

func passwordLessSendCodeLink(m *Metrics, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}

	switch log.GetType() {
	case successSendCodeLink:
		increaseCounter(m.passwordLessCodeLinkCounter, log.GetType())
	case successSendCode:
		increaseCounter(m.passwordLessCodeLinkCounter, log.GetType())
	default:
		return errors.Annotate(errInvalidLogEvent, "code_link event handler can't handle event")
	}

	return nil
}
