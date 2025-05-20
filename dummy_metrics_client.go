package gometrix

import "time"

type DummyMetricsClient struct{}

func (n *DummyMetricsClient) Increment(name string, count int64, tagMap ...MetricTag)         {}
func (n *DummyMetricsClient) Decrement(name string, count int64, tagMap ...MetricTag)         {}
func (n *DummyMetricsClient) Count(name string, value int64, tagMap ...MetricTag)             {}
func (n *DummyMetricsClient) Gauge(name string, value float64, tagMap ...MetricTag)           {}
func (n *DummyMetricsClient) Timing(name string, duration time.Duration, tagMap ...MetricTag) {}

func (n *DummyMetricsClient) Stop() {}
