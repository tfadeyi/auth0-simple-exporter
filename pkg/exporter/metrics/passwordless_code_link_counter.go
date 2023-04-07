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

func passwordLessSendCodeLinkCounterMetric(namespace, subsystem string, applications []*management.Client) *prometheus.CounterVec {
	m := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, tenantSendCodeLink),
			Help: "The total number of send_code_link operations. (codes: cls,cs)",
		}, []string{"type", "client"})
	for _, client := range applications {
		initCounter(m, "cls", client.GetName())
		initCounter(m, "cs", client.GetName())
	}
	return m
}

func passwordLessSendCodeLink(m *Metrics, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}

	switch log.GetType() {
	case successSendCodeLink:
		increaseCounter(m.passwordLessCodeLinkCounter, log.GetType(), log.GetClientName())
	case successSendCode:
		increaseCounter(m.passwordLessCodeLinkCounter, log.GetType(), log.GetClientName())
	default:
		return errors.Annotate(errInvalidLogEvent, "code_link event handler can't handle event")
	}

	return nil
}
