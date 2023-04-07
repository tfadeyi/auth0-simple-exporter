package metrics

import (
	"github.com/auth0/go-auth0/management"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestChangePasswordRequestCounterMetric(t *testing.T) {
	t.Parallel()
	t.Run("the counter is initialise to zero when a new metrics instance is created", func(t *testing.T) {
		var name = "test-app"
		applications := []*management.Client{
			{Name: &name},
		}

		m := New("", "", applications)
		assert.Equal(t, 0, int(getMetricValue(m.changePasswordRequestTotalCounter)))
		assert.Equal(t, 0, int(getMetricValue(m.changePasswordRequestFailCounter)))
	})
	t.Run("the counter is not zero when it is increased", func(t *testing.T) {
		var name = "test-app"
		applications := []*management.Client{
			{Name: &name},
		}

		m := New("", "", applications)
		m.changePasswordRequestFailCounter.WithLabelValues(name).Inc()
		assert.NotEqual(t, 0, int(getMetricValue(m.changePasswordRequestFailCounter)))
	})
	t.Run("the counter errors if the log event is nil", func(t *testing.T) {
		var name = "test-app"
		applications := []*management.Client{
			{Name: &name},
		}
		m := New("", "", applications)

		err := changePasswordRequest(m, nil)
		require.ErrorIs(t, err, errInvalidLogEvent)
	})
	t.Run("the counter errors if the log event cannot be handled", func(t *testing.T) {
		var name = "test-app"
		var code = "crazy-error"
		applications := []*management.Client{
			{Name: &name},
		}
		m := New("", "", applications)

		err := changePasswordRequest(m, &management.Log{ClientName: &name, Type: &code})
		require.Error(t, err)
	})
	t.Run("the counter increases if valid events are passed", func(t *testing.T) {
		var name = "test-app"
		code := failedChangePasswordRequest
		code1 := successfulChangePasswordRequest
		applications := []*management.Client{
			{Name: &name},
		}
		m := New("", "", applications)

		require.NoError(t, changePasswordRequest(m, &management.Log{ClientName: &name, Type: &code}))
		require.NoError(t, changePasswordRequest(m, &management.Log{ClientName: &name, Type: &code1}))
		assert.Equal(t, 1, int(getMetricValue(m.changePasswordRequestFailCounter)))
		assert.Equal(t, 2, int(getMetricValue(m.changePasswordRequestTotalCounter)))
	})
}
