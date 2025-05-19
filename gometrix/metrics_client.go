package gometrix

import "time"

type MetricsClient interface {
	Increment(metricName string, count int64, tagMap map[string]any)
	Decrement(metricName string, count int64, tagMap map[string]any)
	Count(metricName string, value int64, tagMap map[string]any)
	Gauge(metricName string, value float64, tagMap map[string]any)
	Timing(metricName string, duration time.Duration, tagMap map[string]any)

	Stop()
}
