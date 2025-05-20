package gometrix

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/smira/go-statsd"
)

type statsdClient struct {
	client *statsd.Client
}

func NewStatsdMetricsClient(data any) (MetricsClient, error) {
	metricsData, ok := data.(*StatsdMetricsData)
	if !ok {
		return nil, fmt.Errorf("expected StatsdMetricsData, got %T", data)
	}

	address := fmt.Sprintf("%s:%d", metricsData.Host, metricsData.Port)

	logrus.Debugf("StatsD metrics client created [%v]", address)

	client := statsd.NewClient(address, statsd.MetricPrefix(metricsData.Prefix))
	return &statsdClient{client: client}, nil
}

func (s *statsdClient) Stop() {
	s.client.Close()
}

func (s *statsdClient) Increment(name string, count int64, tagMap ...MetricTag) {
	s.client.Incr(joinTagsToName(name, tagMap), count)
}

func (s *statsdClient) Decrement(name string, count int64, tagMap ...MetricTag) {
	s.client.Decr(joinTagsToName(name, tagMap), count)
}

func (s *statsdClient) Count(name string, value int64, tagMap ...MetricTag) {
	s.client.Gauge(joinTagsToName(name, tagMap), value)
}

func (s *statsdClient) Gauge(name string, value float64, tagMap ...MetricTag) {
	s.client.Gauge(joinTagsToName(name, tagMap), int64(value))
}

func (s *statsdClient) Timing(name string, duration time.Duration, tagMap ...MetricTag) {
	s.client.PrecisionTiming(joinTagsToName(name, tagMap), duration)
}
