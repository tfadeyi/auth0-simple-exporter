package operations

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"context"

	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	//failedExchangePasswordMFA Failed Exchange of Password and MFA Recovery code for Access Token
	failedExchangePasswordMFA = "fercft"

	//failedMfaAuth Multi-factor authentication failed. This could happen due to a wrong code entered for SMS/Voice/Email/TOTP factors, or a system failure.
	failedMfaAuth = "gd_auth_failed"

	//mfaAuthRejected A user rejected a Multi-factor authentication request via push-notification.
	mfaAuthRejected = "gd_auth_rejected"

	//mfaAuthSuccess Multi-factor authentication success.
	mfaAuthSuccess = "gd_auth_succeed"

	//gd_send_email 	Email Sent 	Email for MFA successfully sent.

	// gd_send_pn 	Push notification sent 	Push notification for MFA sent successfully sent.
	//gd_send_pn_failure 	Push notification sent 	Push notification for MFA failed.

	//gd_send_sms 	SMS sent 	SMS for MFA successfully sent.
	//gd_send_sms_failure 	SMS sent failures 	Attempt to send SMS for MFA failed.

	//gd_send_voice 	Voice call made 	Voice call for MFA successfully made.
	//gd_send_voice_failure 	Voice call failure 	Attempt to make Voice call for MFA failed.

	//gd_start_auth 	Second factor started 	Second factor authentication event started for MFA.

	//mfar 	MFA Required 	A user has been prompted for multi-factor authentication (MFA). When using Adaptive MFA, Auth0 includes details about the risk assessment.
	//sercft 	Success Exchange 	Successful exchange of Password and MFA Recovery code for Access Token

	TenantMfaOperations = CtxKey("tenant_mfa_operations_total")
)

func NewMfaOperationsMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(TenantMfaOperations)),
			Help: "The number of MFA operations. (codes: fs,ss)",
		}, []string{"status", "type"})
}

func MfaOperationsEventHandler(ctx context.Context, log *management.Log) error {
	if ctx == nil {
		return errInvalidContext
	}
	if log == nil {
		return errInvalidLogEvent
	}
	op, ok := ctx.Value(TenantMfaOperations).(*prometheus.CounterVec)
	if !ok {
		return errors.Annotate(errMissingLogEventMetric, "mfa metric is not in the context")
	}
	switch log.GetType() {
	case failedMfaAuth:
		op.WithLabelValues(failed, log.GetType()).Inc()
	case failedMfaAuth:
		op.WithLabelValues(failed, log.GetType()).Inc()
	case failedExchangePasswordMFA:
		op.WithLabelValues(failed, log.GetType()).Inc()
	case mfaAuthSuccess:
		op.WithLabelValues(success, log.GetType()).Inc()
	default:
		return errors.Annotate(errInvalidLogEvent, "mfa event handler can't handle event")
	}

	return nil
}
