package metrics

import "time"

type metricInfo struct {
	Name               string
	MetricType         int
	ReportingFrequency time.Duration
}

const GaugeMetric = 1
const CounterMetric = 1
