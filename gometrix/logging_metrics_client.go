package gometrix

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type LoggingMetricsClient struct {
	metrics      map[string]*metricAggregator
	printTimeout time.Duration
	totalTime    float64
	stopChan     chan struct{}
	waitGroup    sync.WaitGroup
	startOnce    sync.Once
	mutex        sync.Mutex
	logger       *logrus.Logger
}

type metricAggregator struct {
	total      float64
	count      float64
	last       float64
	lastCount  float64
	min        float64
	max        float64
	metricType string
}

const (
	defaultTimeout  int = 60 //seconds
	defaultMaxFiles int = 3
	defaultMaxSize  int = 10 //MB
)

func (l *LoggingMetricsData) getConfigFromInterface(data interface{}) error {
	config, ok := data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid logging metrics client config type: %T", data)
	}

	if timeout, ok := config["timeout"].(int); ok {
		l.Timeout = timeout
	} else {
		l.Timeout = defaultTimeout
	}

	if path, ok := config["log_file_path"].(string); ok {
		l.LogFilePath = path
	}

	if maxFiles, ok := config["max_files"].(int); ok {
		l.MaxFiles = maxFiles
	} else {
		l.MaxFiles = defaultMaxFiles
	}

	if maxSize, ok := config["max_file_size"].(int); ok {
		l.MaxFileSize = maxSize
	} else {
		l.MaxFileSize = defaultMaxSize
	}

	return nil
}

func NewLoggingMetricsClient(data interface{}) (*LoggingMetricsClient, error) {
	metricsData := LoggingMetricsData{}
	if err := metricsData.getConfigFromInterface(data); err != nil {
		logrus.Errorf("Error loading metrics configuración [%v]", metricsData)
		return nil, err
	}

	client := LoggingMetricsClient{
		metrics:      make(map[string]*metricAggregator),
		printTimeout: time.Duration(float64(metricsData.Timeout) * float64(time.Second)),
		stopChan:     make(chan struct{}),
		logger:       setupMetricsLogger(metricsData),
	}
	client.start()
	return &client, nil
}

func (c *LoggingMetricsClient) start() {
	c.startOnce.Do(func() {
		c.waitGroup.Add(1)
		go c.run()
	})
}

func (c *LoggingMetricsClient) Stop() {
	close(c.stopChan)
	c.waitGroup.Wait()
}

func (c *LoggingMetricsClient) run() {
	defer c.waitGroup.Done()
	ticker := time.NewTicker(c.printTimeout)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.printMetrics()
			c.resetLast()
		case <-c.stopChan:
			return
		}
	}
}

func (c *LoggingMetricsClient) Increment(metricName string, count int64, _ map[string]any) {
	c.Count(metricName, count, nil)
}

func (c *LoggingMetricsClient) Decrement(metricName string, count int64, _ map[string]any) {
	c.Count(metricName, -count, nil)
}

func (c *LoggingMetricsClient) Count(metricName string, value int64, _ map[string]any) {
	c.updateMetric(metricName, "count", float64(value))
}

func (c *LoggingMetricsClient) Gauge(metricName string, value float64, _ map[string]any) {
	c.updateMetric(metricName, "gauge", value)
}

func (c *LoggingMetricsClient) Timing(metricName string, duration time.Duration, _ map[string]any) {
	c.updateMetric(metricName, "timing", float64(duration.Milliseconds()))
}

func (c *LoggingMetricsClient) updateMetric(name, mType string, value float64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	m, exists := c.metrics[name]
	if !exists {
		m = &metricAggregator{
			metricType: mType,
			min:        value,
			max:        value,
		}
		c.metrics[name] = m
	}

	m.total += value
	m.count++
	m.last += value
	m.lastCount++

	if mType == "gauge" || mType == "timing" {
		if value < m.min {
			m.min = value
		}
		if value > m.max {
			m.max = value
		}
	}
}

func (c *LoggingMetricsClient) resetLast() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, m := range c.metrics {
		m.last = 0
		m.lastCount = 0
	}
}

func (c *LoggingMetricsClient) printMetrics() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if len(c.metrics) == 0 {
		c.logger.Info("[METRICS] NONE")
		return
	}

	var output strings.Builder
	output.WriteString("[METRICS]\n")

	metricsByType := make(map[string][]string)
	for name := range c.metrics {
		metricsByType[c.metrics[name].metricType] = append(metricsByType[c.metrics[name].metricType], name)
	}

	var types []string
	for t := range metricsByType {
		types = append(types, t)
	}
	sort.Strings(types)

	c.totalTime += float64(c.printTimeout.Seconds())
	for _, t := range types {
		output.WriteString(fmt.Sprintf("[%s]\n", strings.ToUpper(t)))
		sort.Strings(metricsByType[t])

		for _, name := range metricsByType[t] {
			m := c.metrics[name]
			switch t {
			case "count":
				output.WriteString(fmt.Sprintf(
					"\t[%s] COUNT[%.0f] AVG_PER_SECOND[%.2f] LAST[%.0f] LAST_PER_SECOND[%.2f]\n",
					name, m.total, m.total/c.totalTime, m.last, m.last/c.totalTime,
				))
			case "gauge", "timing":
				avg := m.total / m.count
				lastAvg := 0.0
				if m.lastCount > 0 {
					lastAvg = m.last / m.lastCount
				}
				output.WriteString(fmt.Sprintf(
					"\t[%s] AVG[%.2f] MIN[%.2f] MAX[%.2f] COUNT[%.0f] LAST_AVG[%.2f] LAST_COUNT[%.0f]\n",
					name, avg, m.min, m.max, m.count, lastAvg, m.lastCount,
				))
			}
		}
	}

	c.logger.Info(output.String())
}
