package metrics

import "time"

type metricInfo struct {
	Name               string
	ReportingFrequency time.Duration
}
