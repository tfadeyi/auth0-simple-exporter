package metrics

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	failedPostChangePasswordHook     = "fcph"
	successfulPostChangePasswordHook = "scph"

	tenantTotalPostChangePasswordHook  = "tenant_post_change_password_hook_total"
	tenantFailedPostChangePasswordHook = "tenant_failed_post_change_password_hook_total"
)

func postChangePasswordHookTotalCounterMetric(namespace, subsystem string, applications []*management.Client) *prometheus.CounterVec {
	m := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, tenantTotalPostChangePasswordHook),
			Help: "The total number of post change user password hook operations on the tenant. (codes: scph,fcph)",
		}, []string{"client"})
	for _, client := range applications {
		initCounter(m, client.GetName())
	}
	return m
}

func postChangePasswordHookFailCounterMetric(namespace, subsystem string, applications []*management.Client) *prometheus.CounterVec {
	m := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, tenantFailedPostChangePasswordHook),
			Help: "The number of failed post change user password hook operations on the tenant. (codes: fcph)",
		}, []string{"client"})
	for _, client := range applications {
		initCounter(m, client.GetName())
	}
	return m
}

func postChangePasswordHook(m *Metrics, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}

	switch log.GetType() {
	case successfulPostChangePasswordHook:
		increaseCounter(m.postChangePasswordHookTotalCounter, log.GetClientName())
	case failedPostChangePasswordHook:
		increaseCounter(m.postChangePasswordHookFailCounter, log.GetClientName())
		increaseCounter(m.postChangePasswordHookTotalCounter, log.GetClientName())
	default:
		return errors.Annotate(errInvalidLogEvent, "post change user password hook operations event handler can't handle event")
	}

	return nil
}
