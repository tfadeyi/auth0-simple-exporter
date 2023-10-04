package metrics

import (
	"time"

	"github.com/auth0/go-auth0/management"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	tenantTotalMonthlyActiveUsers = "tenant_total_monthly_active_users"
)

func monthlyActiveUsersCounterMetric(namespace, subsystem string, applications []*management.Client) *prometheus.Counter {
	m := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, tenantTotalMonthlyActiveUsers),
			Help: "The total number of monthly active users on the tenant.",
		})
	return &m
}

func processMonthlyActiveUsers(m *Metrics, users []*management.User) error {
	currentTime := time.Now()
	startTime := currentTime.AddDate(0, 0, -30)

	for _, user := range users {
		if user.LastLogin != nil && (*user.LastLogin).After(startTime) {
			(*m.monthlyActiveUsersCounterMetric).Inc()
		}
	}
	return nil
}
