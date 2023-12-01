package metrics

import (
	"time"

	"github.com/auth0/go-auth0/management"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	tenantTotalMonthlyActiveUsers = "tenant_total_monthly_active_users"
)

func monthlyActiveUsersGaugeMetric(namespace, subsystem string, applications []*management.Client) *prometheus.Gauge {
	m := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, tenantTotalMonthlyActiveUsers),
			Help: "The total number of monthly active users on the tenant.",
		})
	return &m
}

func processMonthlyActiveUsers(m *Metrics, users []*management.User) error {
	if len(users) == 0 {
		return nil
	}

	currentTime := time.Now()
	startTime := currentTime.AddDate(0, 0, -30)
	count := 0.0

	for _, user := range users {
		if user.LastLogin != nil && (*user.LastLogin).After(startTime) {
			count += 1
		}
	}

	(*m.monthlyActiveUsersGaugeMetric).Set(count)

	return nil
}
