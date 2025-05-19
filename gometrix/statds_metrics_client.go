package gometrix

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/smira/go-statsd"
)

type StatsdClient struct {
	client *statsd.Client
}

func (s *StatsdMetricsData) getConfigFromInterface(data interface{}) error {
	config, ok := data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid statsd client config type: %T", data)
	}
	if host, ok := config["host"].(string); ok {
		s.ServerHost = host
	} else {
		return fmt.Errorf("host not found in statsd client config")
	}

	if port, ok := config["port"].(int); ok {
		s.ServerPort = int64(port)
	}
	return nil
}

func NewStatsdMetricsClient(data interface{}) (*StatsdClient, error) {
	metricsData := StatsdMetricsData{}
	err := metricsData.getConfigFromInterface(data)
	if err != nil {
		return nil, err
	}

	address := fmt.Sprintf("%s:%d", metricsData.ServerHost, metricsData.ServerPort)

	logrus.Debugf("StatsD metrics client created [%v]", address)

	client := statsd.NewClient(address, statsd.MetricPrefix(metricsData.Prefix))
	return &StatsdClient{client: client}, nil
}

func (s *StatsdClient) Stop() {
	s.client.Close()
}

func (s *StatsdClient) Increment(name string, count int64, tagMap map[string]any) {
	s.client.Incr(name, count, convertToTags(tagMap)...)
}

func (s *StatsdClient) Decrement(name string, count int64, tagMap map[string]any) {
	s.client.Decr(name, count, convertToTags(tagMap)...)
}

func (s *StatsdClient) Count(name string, value int64, tagMap map[string]any) {
	s.client.Gauge(name, value, convertToTags(tagMap)...)
}

func (s *StatsdClient) Gauge(name string, value float64, tagMap map[string]any) {
	s.client.Gauge(name, int64(value), convertToTags(tagMap)...)
}

func (s *StatsdClient) Timing(name string, duration time.Duration, tagMap map[string]any) {
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
