package gometrix

import "time"

type MetricTag struct {
	name  string
	value any
}

type MetricsClient interface {
	Increment(metricName string, count int64, tagMap ...MetricTag)
	Decrement(metricName string, count int64, tagMap ...MetricTag)
	Count(metricName string, value int64, tagMap ...MetricTag)
	Gauge(metricName string, value float64, tagMap ...MetricTag)
	Timing(metricName string, duration time.Duration, tagMap ...MetricTag)

	Stop()
}
