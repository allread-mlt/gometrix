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

func (s *statsdClient) Increment(name string, count int64, tagMap map[string]any) {
	s.client.Incr(name, count, convertToTags(tagMap)...)
}

func (s *statsdClient) Decrement(name string, count int64, tagMap map[string]any) {
	s.client.Decr(name, count, convertToTags(tagMap)...)
}

func (s *statsdClient) Count(name string, value int64, tagMap map[string]any) {
	s.client.Gauge(name, value, convertToTags(tagMap)...)
}

func (s *statsdClient) Gauge(name string, value float64, tagMap map[string]any) {
	s.client.Gauge(name, int64(value), convertToTags(tagMap)...)
}

func (s *statsdClient) Timing(name string, duration time.Duration, tagMap map[string]any) {
	s.client.PrecisionTiming(name, duration, convertToTags(tagMap)...)
}

func convertToTags(tagMap map[string]any) []statsd.Tag {
	if len(tagMap) == 0 {
		return nil
	}

	tags := make([]statsd.Tag, 0, len(tagMap))

	for k, v := range tagMap {
		switch val := v.(type) {
		case string:
			tags = append(tags, statsd.StringTag(k, val))
		case int:
			tags = append(tags, statsd.IntTag(k, val))
		case int64:
			tags = append(tags, statsd.Int64Tag(k, val))
		default:
			tags = append(tags, statsd.StringTag(k, fmt.Sprintf("%v", val)))
		}
	}

	return tags
}
