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
	errInvalidLogEvent = errors.New("event handler doesn't accept the event log type")
)

type (
	// See the NewMetrics func for proper descriptions and prometheus names!
	// In case you add a metric here later, make sure to include it in the
	// List method or you'll going to have a bad time.
	Metrics struct {
		loginTotalCounter *prometheus.CounterVec
		loginFailCounter  *prometheus.CounterVec

		logoutTotalCounter *prometheus.CounterVec
		logoutFailCounter  *prometheus.CounterVec

		signupTotalCounter *prometheus.CounterVec
		signupFailCounter  *prometheus.CounterVec

		changePasswordTotalCounter *prometheus.CounterVec
		changePasswordFailCounter  *prometheus.CounterVec

		changeEmailTotalCounter *prometheus.CounterVec
		changeEmailFailCounter  *prometheus.CounterVec

		apiOperationTotalCounter *prometheus.CounterVec
		apiOperationFailCounter  *prometheus.CounterVec

		changePhoneNumberTotalCounter *prometheus.CounterVec
		changePhoneNumberFailCounter  *prometheus.CounterVec

		deleteUserTotalCounter *prometheus.CounterVec
		deleteUserFailCounter  *prometheus.CounterVec

		passwordLessCodeLinkCounter *prometheus.CounterVec

		postChangePasswordHookTotalCounter *prometheus.CounterVec
		postChangePasswordHookFailCounter  *prometheus.CounterVec

		pushNotificationTotalCounter *prometheus.CounterVec
		pushNotificationFailCounter  *prometheus.CounterVec

		successfulSendEmailCounter *prometheus.CounterVec

		sendSMSTotalCounter *prometheus.CounterVec
		sendSMSFailCounter  *prometheus.CounterVec

		voiceCallTotalCounter *prometheus.CounterVec
		voiceCallFailCounter  *prometheus.CounterVec

		changePasswordRequestTotalCounter *prometheus.CounterVec
		changePasswordRequestFailCounter  *prometheus.CounterVec
	}

	LogEventFunc func(m *Metrics, log *management.Log) error
)

// Creates and populates a new Metrics struct
// This is where all the prometheus metrics, names and labels are specified
func New(namespace, subsystem string, applications []*management.Client) *Metrics {
	m := &Metrics{
		loginTotalCounter: loginTotalCounterMetric(namespace, subsystem, applications),
		loginFailCounter:  loginFailCounterMetric(namespace, subsystem, applications),

		logoutTotalCounter: logoutTotalCounterMetric(namespace, subsystem, applications),
		logoutFailCounter:  logoutFailCounterMetric(namespace, subsystem, applications),

		signupTotalCounter: signupTotalCounterMetric(namespace, subsystem, applications),
		signupFailCounter:  signupFailCounterMetric(namespace, subsystem, applications),

		changePasswordTotalCounter: changePasswordTotalCounterMetric(namespace, subsystem, applications),
		changePasswordFailCounter:  changePasswordFailCounterMetric(namespace, subsystem, applications),

		changeEmailTotalCounter: changeEmailTotalCounterMetric(namespace, subsystem, applications),
		changeEmailFailCounter:  changeEmailFailCounterMetric(namespace, subsystem, applications),

		apiOperationTotalCounter: APIOperationTotalCounterMetric(namespace, subsystem, applications),
		apiOperationFailCounter:  APIOperationFailCounterMetric(namespace, subsystem, applications),

		changePhoneNumberTotalCounter: changePhoneNumberTotalCounterMetric(namespace, subsystem, applications),
		changePhoneNumberFailCounter:  changePhoneNumberFailCounterMetric(namespace, subsystem, applications),

		deleteUserTotalCounter: deleteUserTotalCounterMetric(namespace, subsystem, applications),
		deleteUserFailCounter:  deleteUserFailCounterMetric(namespace, subsystem, applications),

		passwordLessCodeLinkCounter: passwordLessSendCodeLinkCounterMetric(namespace, subsystem, applications),

		postChangePasswordHookTotalCounter: postChangePasswordHookTotalCounterMetric(namespace, subsystem, applications),
		postChangePasswordHookFailCounter:  postChangePasswordHookFailCounterMetric(namespace, subsystem, applications),

		pushNotificationTotalCounter: pushNotificationTotalCounterMetric(namespace, subsystem, applications),
		pushNotificationFailCounter:  pushNotificationFailCounterMetric(namespace, subsystem, applications),

		successfulSendEmailCounter: sendEmailCounterMetric(namespace, subsystem, applications),

		sendSMSTotalCounter: sendSMSTotalCounterMetric(namespace, subsystem, applications),
		sendSMSFailCounter:  sendSMSFailCounterMetric(namespace, subsystem, applications),

		voiceCallTotalCounter: sendVoiceCallTotalCounterMetric(namespace, subsystem, applications),
		voiceCallFailCounter:  sendVoiceCallFailCounterMetric(namespace, subsystem, applications),

		changePasswordRequestTotalCounter: changePasswordRequestTotalCounterMetric(namespace, subsystem, applications),
		changePasswordRequestFailCounter:  changePasswordRequestFailCounterMetric(namespace, subsystem, applications),
	}
	return m
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

// List is needed by the server so it can register and collect these metrics
func (m *Metrics) List() []prometheus.Collector {
	return []prometheus.Collector{
		// ADD EVERY METRIC HERE!
		m.loginTotalCounter,
		m.loginFailCounter,

		m.apiOperationTotalCounter,
		m.apiOperationFailCounter,

		m.logoutTotalCounter,
		m.logoutFailCounter,

		m.changeEmailTotalCounter,
		m.changeEmailFailCounter,

		m.changePasswordTotalCounter,
		m.changePasswordFailCounter,

		m.deleteUserTotalCounter,
		m.deleteUserFailCounter,

		m.postChangePasswordHookTotalCounter,
		m.postChangePasswordHookFailCounter,

		m.pushNotificationTotalCounter,
		m.pushNotificationFailCounter,

		m.passwordLessCodeLinkCounter,

		m.successfulSendEmailCounter,

		m.sendSMSTotalCounter,
		m.sendSMSFailCounter,

		m.signupTotalCounter,
		m.signupFailCounter,

		m.voiceCallTotalCounter,
		m.voiceCallFailCounter,

		m.changePasswordRequestTotalCounter,
		m.changePasswordRequestFailCounter,
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

func initCounter(m *prometheus.CounterVec, labels ...string) {
	m.WithLabelValues(labels...)
}

// This will push your metrics object into every request context for later use
func Middleware(next echo.HandlerFunc, applicationClients []*management.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		namespace := c.Get(namespaceCtxKey).(string)
		subsystem := c.Get(subsystemCtxKey).(string)
		c.Set(ListCtxKey, New(namespace, subsystem, applicationClients))
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
