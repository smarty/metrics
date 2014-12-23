package metrics

import "time"

type metricInfo struct {
	Name               string
	MetricType         int
	ReportingFrequency time.Duration
}

const CounterMetric = 1
const GaugeMetric = 2
