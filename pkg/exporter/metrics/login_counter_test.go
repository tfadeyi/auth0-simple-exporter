package metrics

import (
	"github.com/auth0/go-auth0/management"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLoginCounterMetric(t *testing.T) {
	t.Parallel()
	t.Run("the counter is initialise to zero when a new metrics instance is created", func(t *testing.T) {
		var name = "test-app"
		applications := []*management.Client{
			{Name: &name},
		}

		m := New("", "", applications)
		assert.Equal(t, 0, int(getMetricValue(m.loginFailCounter)))
		assert.Equal(t, 0, int(getMetricValue(m.loginTotalCounter)))
	})
	t.Run("the counter is not zero when it is increased", func(t *testing.T) {
		var name = "test-app"
		applications := []*management.Client{
			{Name: &name},
		}

		m := New("", "", applications)
		m.loginFailCounter.WithLabelValues(name).Inc()
		assert.NotEqual(t, 0, int(getMetricValue(m.loginFailCounter)))
	})
	t.Run("the counter errors if the log event is nil", func(t *testing.T) {
		var name = "test-app"
		applications := []*management.Client{
			{Name: &name},
		}
		m := New("", "", applications)

		err := login(m, nil)
		require.ErrorIs(t, err, errInvalidLogEvent)
	})
	t.Run("the counter errors if the log event cannot be handled", func(t *testing.T) {
		var name = "test-app"
		var code = "invalid-error"
		applications := []*management.Client{
			{Name: &name},
		}
		m := New("", "", applications)

		err := login(m, &management.Log{ClientName: &name, Type: &code})
		require.Error(t, err)
	})
	t.Run("the counter increases if valid events are passed", func(t *testing.T) {
		var name = "test-app"
		code := failedLogin
		code1 := failedLoginWithIncorrectPassword
		code2 := failedLoginWithIncorrectUsername
		code3 := successfulLogin
		applications := []*management.Client{
			{Name: &name},
		}
		m := New("", "", applications)

		require.NoError(t, login(m, &management.Log{ClientName: &name, Type: &code}))
		require.NoError(t, login(m, &management.Log{ClientName: &name, Type: &code1}))
		require.NoError(t, login(m, &management.Log{ClientName: &name, Type: &code2}))
		require.NoError(t, login(m, &management.Log{ClientName: &name, Type: &code3}))
		assert.Equal(t, 3, int(getMetricValue(m.loginFailCounter)))
		assert.Equal(t, 4, int(getMetricValue(m.loginTotalCounter)))
	})
}
