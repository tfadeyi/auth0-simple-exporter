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

	tenantSuccessPostChangePasswordHook = "tenant_success_post_change_password_hook_total"
	tenantFailedPostChangePasswordHook  = "tenant_failed_post_change_password_hook_total"
)

func successPostChangePasswordHookCounterMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(tenantSuccessPostChangePasswordHook)),
			Help: "The number of successful post change user password hook operations on the tenant. (codes: scph)",
		}, []string{})
}

func failPostChangePasswordHookCounterMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(tenantFailedPostChangePasswordHook)),
			Help: "The number of failed post change user password hook operations on the tenant. (codes: fcph)",
		}, []string{})
}

func postChangePasswordHook(m *Metrics, log *management.Log) error {
	if log == nil {
		return errInvalidLogEvent
	}
	switch log.GetType() {
	case failedPostChangePasswordHook:
		increaseCounter(m.failedPostChangePasswordHookCounter)
	case successfulPostChangePasswordHook:
		increaseCounter(m.successfulPostChangePasswordHookCounter)
	default:
		return errors.Annotate(errInvalidLogEvent, "post change user password hook operations event handler can't handle event")
	}

	return nil
}
