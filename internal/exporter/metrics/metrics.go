package metrics

import (
	"github.com/auth0/go-auth0/management"
	"github.com/juju/errors"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespaceCtxKey = "metrics-namespace"
	subsystemCtxKey = "metrics-subsystem"
	ListCtxKey      = "metrics-list"
)

var (
	errInvalidLogEvent       = errors.New("event handler doesn't accept the event log type")
	errMissingLogEventMetric = errors.New("couldn't find the prometheus metric for the required log event")
)

type (
	// See the NewMetrics func for proper descriptions and prometheus names!
	// In case you add a metric here later, make sure to include it in the
	// List method or you'll going to have a bad time.
	Metrics struct {
		successfulLoginCnt *prometheus.CounterVec
		failedLoginCnt     *prometheus.CounterVec

		successfulLogoutCounter *prometheus.CounterVec
		failedLogoutCounter     *prometheus.CounterVec

		successfulChangePasswordCounter *prometheus.CounterVec
		failedChangePasswordCounter     *prometheus.CounterVec

		successfulChangeEmailCounter *prometheus.CounterVec
		failedChangeEmailCounter     *prometheus.CounterVec

		successfulAPIOperationCounter *prometheus.CounterVec
		failedAPIOperationCounter     *prometheus.CounterVec

		successfulChangePhoneNumberCounter *prometheus.CounterVec
		failedChangePhoneNumberCounter     *prometheus.CounterVec

		successfulDeleteUserCounter *prometheus.CounterVec
		failedDeleteUserCounter     *prometheus.CounterVec

		passwordLessCodeLinkCounter *prometheus.CounterVec

		successfulPostChangePasswordHookCounter *prometheus.CounterVec
		failedPostChangePasswordHookCounter     *prometheus.CounterVec

		successfulPushNotificationCounter *prometheus.CounterVec
		failedPushNotificationCounter     *prometheus.CounterVec

		successfulSendEmailCounter *prometheus.CounterVec

		successfulSendSMSCounter *prometheus.CounterVec
		failedSendSMSCounter     *prometheus.CounterVec

		successfulSignupCounter *prometheus.CounterVec
		failedSignupCounter     *prometheus.CounterVec

		successfulVoiceCallCounter *prometheus.CounterVec
		failedVoiceCallCounter     *prometheus.CounterVec

		successfulChangePasswordRequestCounter *prometheus.CounterVec
		failedChangePasswordRequestCounter     *prometheus.CounterVec
	}

	LogEventFunc func(m *Metrics, log *management.Log) error
)

// Creates and populates a new Metrics struct
// This is where all the prometheus metrics, names and labels are specified
func New(namespace, subsystem string) *Metrics {
	return &Metrics{
		successfulLoginCnt: successLoginCounterMetric(namespace, subsystem),
		failedLoginCnt:     failLoginCounterMetric(namespace, subsystem),

		successfulLogoutCounter: successLogoutCounterMetric(namespace, subsystem),
		failedLogoutCounter:     failLogoutCounterMetric(namespace, subsystem),

		successfulChangePasswordCounter: successChangePasswordMetric(namespace, subsystem),
		failedChangePasswordCounter:     failChangePasswordCounterMetric(namespace, subsystem),

		successfulChangeEmailCounter: successChangeEmailCounterMetric(namespace, subsystem),
		failedChangeEmailCounter:     failedChangeEmailCounterMetric(namespace, subsystem),

		successfulAPIOperationCounter: successAPIOperationCounterMetric(namespace, subsystem),
		failedAPIOperationCounter:     failAPIOperationCounterMetric(namespace, subsystem),

		successfulChangePhoneNumberCounter: successChangePhoneNumberMetric(namespace, subsystem),
		failedChangePhoneNumberCounter:     failChangePhoneNumberMetric(namespace, subsystem),

		successfulDeleteUserCounter: successDeleteUserCounterMetric(namespace,subsystem),
		failedDeleteUserCounter:     failDeleteUserCounterMetric(namespace, subsystem),

		passwordLessCodeLinkCounter: passwordLessSendCodeLinkCounterMetric(namespace, subsystem),

		successfulPostChangePasswordHookCounter: successPostChangePasswordHookCounterMetric(namespace, subsystem),
		failedPostChangePasswordHookCounter: failPostChangePasswordHookCounterMetric(namespace, subsystem),

		successfulPushNotificationCounter: successPushNotificationCounterMetric(namespace, subsystem),
		failedPushNotificationCounter:     failPushNotificationCounterMetric(namespace, subsystem),

		successfulSendEmailCounter: sendEmailCounterMetric(namespace, subsystem),

		successfulSendSMSCounter: successSendSMSOperationsMetric(namespace, subsystem),
		failedSendSMSCounter:     failSendSMSOperationsMetric(namespace, subsystem),

		successfulSignupCounter: successSignupCounterMetric(namespace, subsystem),
		failedSignupCounter: failSignupCounterMetric(namespace, subsystem),

		successfulVoiceCallCounter: successSendVoiceCallCounterMetric(namespace, subsystem),
		failedVoiceCallCounter: failSendVoiceCallCounterMetric(namespace, subsystem),

		successfulChangePasswordRequestCounter: successChangePasswordRequestCounterMetric(namespace, subsystem),
		failedChangePasswordRequestCounter:     failChangePasswordRequestCounterMetric(namespace, subsystem),
	}
}

func (m *Metrics) logEventHandlers() []LogEventFunc {
	return []LogEventFunc{
		// ADD EVERY EVENT HANDLER HERE!
		login,
		logout,
		changePassword,
		changeEmail,
		apiOperations,
		changePhoneNumber,
		passwordLessSendCodeLink,
		deleteUser,
		changePasswordRequest,
		postChangePasswordHook,
		pushNotification,
		sendEmail,
		sendSMS,
		signup,
		sendVoiceCall,
	}
}

// Needed by echo-contrib so echo can register and collect these metrics
func (m *Metrics) List() []prometheus.Collector {
	return []prometheus.Collector{
		// ADD EVERY METRIC HERE!
		m.successfulLoginCnt,
		m.failedLoginCnt,

		m.successfulAPIOperationCounter,
		m.failedAPIOperationCounter,

		m.successfulLogoutCounter,
		m.failedLogoutCounter,

		m.successfulChangeEmailCounter,
		m.failedChangeEmailCounter,

		m.successfulChangePasswordCounter,
		m.failedChangePasswordCounter,

		m.successfulDeleteUserCounter,
		m.failedDeleteUserCounter,

		m.successfulPostChangePasswordHookCounter,
		m.failedPostChangePasswordHookCounter,

		m.successfulPushNotificationCounter,
		m.failedPushNotificationCounter,

		m.passwordLessCodeLinkCounter,

		m.successfulSendEmailCounter,

		m.successfulSendSMSCounter ,
		m.failedSendSMSCounter    ,

		m.successfulSignupCounter,
		m.failedSignupCounter,

		m.successfulVoiceCallCounter,
		m.failedVoiceCallCounter,

		m.successfulChangePasswordRequestCounter,
		m.failedChangePasswordRequestCounter,
	}
}

func (m *Metrics) Update(log *management.Log) error {
	for _, fs := range m.logEventHandlers() {
		if err := fs(m, log); err != nil {
			continue
		}
		// success in updating the metrics, goto next event/log
		break
	}
	return nil
}

func increaseCounter(m *prometheus.CounterVec, labels ...string) {
	m.WithLabelValues(labels...).Inc()
}

// This will push your metrics object into every request context for later use
func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		namespace := c.Get(namespaceCtxKey).(string)
		subsystem := c.Get(subsystemCtxKey).(string)
		c.Set(ListCtxKey, New(namespace, subsystem))
		return next(c)
	}
}

func NamespaceMiddleware(next echo.HandlerFunc, namespace string) echo.HandlerFunc {
	// propagate exporter namespace and subsystem
	return func(ctx echo.Context) error {
		ctx.Set(namespaceCtxKey, namespace)
		return next(ctx)
	}
}

func SubsystemMiddleware(next echo.HandlerFunc, subsystem string) echo.HandlerFunc {
	// propagate exporter namespace and subsystem
	return func(ctx echo.Context) error {
		ctx.Set(subsystemCtxKey, subsystem)
		return next(ctx)
	}
}
