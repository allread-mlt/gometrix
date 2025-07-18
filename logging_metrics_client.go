package gometrix

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type loggingMetricsClient struct {
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

func NewLoggingMetricsClient(data any) (MetricsClient, error) {
	metricsData, ok := data.(*LoggingMetricsData)
	if !ok {
		return nil, fmt.Errorf("expected LoggingMetricsData, got %T", data)
	}

	client := loggingMetricsClient{
		metrics:      make(map[string]*metricAggregator),
		printTimeout: time.Duration(float64(metricsData.Timeout) * float64(time.Second)),
		stopChan:     make(chan struct{}),
		logger:       setupMetricsLogger(*metricsData),
	}
	client.start()
	return &client, nil
}

func (c *loggingMetricsClient) start() {
	c.startOnce.Do(func() {
		c.waitGroup.Add(1)
		go c.run()
	})
}

func (c *loggingMetricsClient) Stop() {
	close(c.stopChan)
	c.waitGroup.Wait()
}

func (c *loggingMetricsClient) run() {
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

func (c *loggingMetricsClient) Increment(name string, count int64, tagMap ...MetricTag) {
	c.Count(joinTagsToName(name, tagMap), count)
}

func (c *loggingMetricsClient) Decrement(name string, count int64, tagMap ...MetricTag) {
	c.Count(joinTagsToName(name, tagMap), -count)
}

func (c *loggingMetricsClient) Count(name string, value int64, tagMap ...MetricTag) {
	c.updateMetric(joinTagsToName(name, tagMap), "count", float64(value))
}

func (c *loggingMetricsClient) Gauge(name string, value float64, tagMap ...MetricTag) {
	c.updateMetric(joinTagsToName(name, tagMap), "gauge", value)
}

func (c *loggingMetricsClient) Timing(name string, duration time.Duration, tagMap ...MetricTag) {
	c.updateMetric(joinTagsToName(name, tagMap), "timing", float64(duration.Milliseconds()))
}

func (c *loggingMetricsClient) updateMetric(name, mType string, value float64) {
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

func (c *loggingMetricsClient) resetLast() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, m := range c.metrics {
		m.last = 0
		m.lastCount = 0
	}
}

func (c *loggingMetricsClient) printMetrics() {
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
