package gometrix

import (
	"github.com/sirupsen/logrus"
)

type MetricsClientType string

const (
	metricsTypeDummy   MetricsClientType = "dummy"
	metricsTypeStatsd  MetricsClientType = "statsd"
	metricsTypeLogging MetricsClientType = "logging"
)

type MetricsClientData any

type MetricsClientConfig struct {
	ClientType MetricsClientType `yaml:"type"`
	ClientData interface{}       `yaml:"data"`
}

func NewMetricsClient(config *MetricsClientConfig) (MetricsClient, error) {
	switch MetricsClientType(config.ClientType) {
	case metricsTypeStatsd:
		logrus.Debugf("Creating statsD metrics client")
		return NewStatsdMetricsClient(config.ClientData)
	case metricsTypeLogging:
		logrus.Debugf("Creating logging metrics client")
		return NewLoggingMetricsClient(config.ClientData)
	case metricsTypeDummy:
		logrus.Debugf("Creating dummy metrics client")
		return &DummyMetricsClient{}, nil
	default:
	}
	logrus.Warnf("Unknown metrics client type: %s. Creating dummy client", config.ClientType)
	return &DummyMetricsClient{}, nil
}
