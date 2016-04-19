package metrics

import "time"

type metricInfo struct {
	Name               string
	MetricType         uint8
	ReportingFrequency time.Duration
}
