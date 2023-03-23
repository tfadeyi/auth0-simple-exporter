package metrics_old

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"

	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
)

const (
	success     = "success"
	failed      = "failed"
	RegistryKey = ctxKey("registry")
	namespace   = "auth0"
	subsystem   = ""
)

type (
	ctxKey     string
	MetricFunc func(ctx context.Context, log *management.Log) error
)

var (
	errInvalidContext        = errors.New("invalid or nil context object")
	errInvalidLogEvent       = errors.New("event handler doesn't accept the event log type")
	errMissingLogEventMetric = errors.New("couldn't find the prometheus metric for the required log event")

	Handlers = []MetricFunc{
		ApiOperationsEventHandler,
		LogoutOperationsEventHandler,
		SendCodeLinkOperationsEventHandler,
		SendVoiceCallOperationsEventHandler,
		SendEmailOperationsEventHandler,
		SendSMSOperationsEventHandler,
		SignupOperationsEventHandler,
		PushNotificationOperationsEventHandler,
		LoginOperationsEventHandler,
		deleteUserOperationsEventHandler,
		changeEmailOperationsEventHandler,
	}
)

func RegistryFromContext(ctx context.Context) *prometheus.Registry {
	registry, ok := ctx.Value(RegistryKey).(*prometheus.Registry)
	if !ok {
		return prometheus.NewRegistry()
	}
	return registry
}

func contextWithRegistry(ctx context.Context, registry *prometheus.Registry) context.Context {
	return context.WithValue(ctx, RegistryKey, registry)
}

func contextWithMetrics(ctx context.Context, mss map[ctxKey]prometheus.Collector) context.Context {
	registry := RegistryFromContext(ctx)

	for name, metric := range mss {
		ctx = context.WithValue(ctx, name, metric)
		registry.MustRegister(metric)
	}
	return ctx
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		mss := map[ctxKey]prometheus.Collector{
			TenantSuccessApiOperations: SapiOperationsMetric(namespace, subsystem),
			TenantFailedApiOperations:  FapiOperationsMetric(namespace, subsystem),

			TenantSuccessfulLoginsOperations: NewSuccessfulLoginOperationsMetric(namespace, subsystem),
			TenantFailedLoginsOperations:     NewFailedLoginOperationsMetric(namespace, subsystem),

			TenantSuccessLogoutOperations: NewSuccessLogoutOperationsMetric(namespace, subsystem),
			TenantFailedLogoutOperations:  NewFailedLogoutOperationsMetric(namespace, subsystem),

			TenantSendCodeLinkOperations:     NewSendCodeLinkOperationsMetric(namespace, subsystem),
			TenantSendEmailOperations:        NewSendEmailOperationsMetric(namespace, subsystem),
			TenantSendSMSOperations:          NewSendSMSOperationsMetric(namespace, subsystem),
			TenantSendVoiceCallOperations:    NewSendVoiceCallOperationsMetric(namespace, subsystem),
			TenantPushNotificationOperations: NewPushNotificationOperationsMetric(namespace, subsystem),
			TenantSignupOperations:           NewSignupOperationsMetric(namespace, subsystem),

			TenantFailedDeleteUserOperations:  failedDeleteUserOperationsMetric(namespace, subsystem),
			TenantSuccessDeleteUserOperations: successDeleteUserOperationsMetric(namespace, subsystem),

			TenantSuccessChangeEmailOperations: successChangeEmailOperationsMetric(namespace, subsystem),
			TenantFailedChangeEmailOperations:  failedChangeEmailOperationsMetric(namespace, subsystem),
		}
		registry := RegistryFromContext(ctx)
		ctx = contextWithRegistry(ctx, registry)
		ctx = contextWithMetrics(ctx, mss)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
