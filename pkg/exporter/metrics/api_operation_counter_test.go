package metrics

import (
	"github.com/auth0/go-auth0/management"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// getMetricValue, is a test function that returns the values of given metrics
func getMetricValue(col prometheus.Collector) float64 {
	c := make(chan prometheus.Metric, 1) // 1 for metric with no vector
	col.Collect(c)                       // collect current metric value into the channel
	m := dto.Metric{}
	_ = (<-c).Write(&m) // read metric value from the channel
	return *m.Counter.Value
}

func TestAPIOperationCounterMetric(t *testing.T) {
	t.Parallel()
	t.Run("the counter is initialise to zero when a new metrics instance is created", func(t *testing.T) {
		var name = "test-app"
		applications := []*management.Client{
			{Name: &name},
		}

		m := New("", "", applications)
		value := getMetricValue(m.apiOperationFailCounter)
		assert.Equal(t, 0, int(value))
	})
	t.Run("the counter is not zero when it is increased", func(t *testing.T) {
		var name = "test-app"
		applications := []*management.Client{
			{Name: &name},
		}

		m := New("", "", applications)
		m.apiOperationFailCounter.WithLabelValues(name).Inc()
		value := getMetricValue(m.apiOperationFailCounter)
		assert.NotEqual(t, 0, int(value))
	})
	t.Run("the counter errors if the log event is nil", func(t *testing.T) {
		var name = "test-app"
		applications := []*management.Client{
			{Name: &name},
		}
		m := New("", "", applications)

		err := apiOperations(m, nil)
		require.ErrorIs(t, err, errInvalidLogEvent)
	})
	t.Run("the counter errors if the log event cannot be handled", func(t *testing.T) {
		var name = "test-app"
		var code = "invalid-error"
		applications := []*management.Client{
			{Name: &name},
		}
		m := New("", "", applications)

		err := apiOperations(m, &management.Log{ClientName: &name, Type: &code})
		require.Error(t, err)
	})
	t.Run("the counter increases if valid events are passed", func(t *testing.T) {
		var name = "test-app"
		code := "fapi"
		code1 := "sapi"
		applications := []*management.Client{
			{Name: &name},
		}
		m := New("", "", applications)

		require.NoError(t, apiOperations(m, &management.Log{ClientName: &name, Type: &code}))
		require.NoError(t, apiOperations(m, &management.Log{ClientName: &name, Type: &code1}))
		assert.Equal(t, 1, int(getMetricValue(m.apiOperationFailCounter)))
		assert.Equal(t, 2, int(getMetricValue(m.apiOperationTotalCounter)))
	})
}
