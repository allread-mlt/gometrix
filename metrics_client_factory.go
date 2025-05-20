package gometrix

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type MetricsClientType string

const (
	metricsTypeDummy   MetricsClientType = "dummy"
	metricsTypeStatsd  MetricsClientType = "statsd"
	metricsTypeLogging MetricsClientType = "logging"
)

type MetricsClientConfig struct {
	Type MetricsClientType `yaml:"type"`
	Data interface{}       `yaml:"data"`
}

func (w *MetricsClientConfig) UnmarshalYAML(value *yaml.Node) error {
	type rawWrapper struct {
		Type MetricsClientType `yaml:"type"`
		Data yaml.Node         `yaml:"data"`
	}
	var temp rawWrapper
	if err := value.Decode(&temp); err != nil {
		return err
	}

	fmt.Printf("Data node kind: %v\n", temp.Data.Kind)
	for i := 0; i < len(temp.Data.Content); i += 2 {
		fmt.Printf("Key: %s, Value: %s\n", temp.Data.Content[i].Value, temp.Data.Content[i+1].Value)
	}

	w.Type = temp.Type
	switch temp.Type {
	case metricsTypeStatsd:
		statsdConfig := &StatsdMetricsData{}
		if err := temp.Data.Decode(statsdConfig); err != nil {
			return err
		}
		w.Data = statsdConfig
	case metricsTypeLogging:
		loggingConfig := &LoggingMetricsData{}
		if err := temp.Data.Decode(loggingConfig); err != nil {
			return err
		}
		w.Data = loggingConfig
	case metricsTypeDummy:
		break
	default:
		return fmt.Errorf("unknown type: %s", temp.Type)
	}
	return nil
}

func NewMetricsClient(config *MetricsClientConfig) (MetricsClient, error) {
	switch MetricsClientType(config.Type) {
	case metricsTypeStatsd:
		logrus.Debugf("Creating statsD metrics client")
		return NewStatsdMetricsClient(config.Data)
	case metricsTypeLogging:
		logrus.Debugf("Creating logging metrics client")
		return NewLoggingMetricsClient(config.Data)
	case metricsTypeDummy:
		logrus.Debugf("Creating dummy metrics client")
		return &DummyMetricsClient{}, nil
	default:
	}
	return nil, fmt.Errorf("unknown metrics client type: %s", config.Type)
}
