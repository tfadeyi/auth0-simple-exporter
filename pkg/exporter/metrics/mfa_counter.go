package metrics

// https://auth0.com/docs/deploy-monitor/logs/log-event-type-codes

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// successExchangePasswordMFA Successful exchange of Password and MFA Recovery code for Access Token
	//successExchangePasswordMFA = "sercft"
	////failedExchangePasswordMFA Failed Exchange of Password and MFA Recovery code for Access Token
	//failedExchangePasswordMFA = "fercft"
	//
	//// failedMfaAuth Multi-factor authentication failed. This could happen due to a wrong code entered for SMS/Voice/Email/TOTP factors, or a system failure.
	//failedMfaAuth = "gd_auth_failed"
	//// successMfaAuth Multi-factor authentication success.
	//successMfaAuth = "gd_auth_succeed"

	// gd_start_auth 	Second factor started 	Second factor authentication event started for MFA.
	// mfaAuthRejected A user rejected a Multi-factor authentication request via push-notification.
	// mfaAuthRejected = "gd_auth_rejected"
	// mfar 	MFA Required 	A user has been prompted for multi-factor authentication (MFA). When using Adaptive MFA, Auth0 includes details about the risk assessment.

	tenantMfaOperations = "tenant_mfa_operations_total"
)

func NewMfaOperationsMetric(namespace, subsystem string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, string(tenantMfaOperations)),
			Help: "The number of MFA operations. (codes: fs,ss)",
		}, []string{"status", "type"})
}

//
//func MfaOperationsEventHandler(ctx context.Context, log *management.Log) error {
//	if log == nil {
//		return errInvalidLogEvent
//	}
//	op, ok := ctx.Value(tenantMfaOperations).(*prometheus.CounterVec)
//	if !ok {
//		return errors.Annotate(errMissingLogEventMetric, "mfa metric is not in the context")
//	}
//	switch log.GetType() {
//	case failedMfaAuth:
//		op.WithLabelValues(failed, log.GetType()).Inc()
//	case successMfaAuth:
//		op.WithLabelValues(success, log.GetType()).Inc()
//	case failedExchangePasswordMFA:
//		op.WithLabelValues(failed, log.GetType()).Inc()
//	case successExchangePasswordMFA:
//		op.WithLabelValues(success, log.GetType()).Inc()
//	default:
//		return errors.Annotate(errInvalidLogEvent, "mfa event handler can't handle event")
//	}
//
//	return nil
//}
