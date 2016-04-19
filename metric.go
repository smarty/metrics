package metrics

import "time"

type Metric interface {
	MeasurementIsOverdue(now time.Time) bool
	Measure() MetricMeasurement
	ScheduleNextMeasurement(now time.Time)
}

////////////////////////////////////////////////////////////////////////////

type MetricMeasurement struct {
	Captured   time.Time
	ID         int
	Name       string
	MetricType uint8
	Value      int64
}

////////////////////////////////////////////////////////////////////////////

type ReportingFrequency struct {
	upcoming time.Time
	interval time.Duration
}

func (this *ReportingFrequency) MeasurementIsOverdue(now time.Time) (overdue bool) {
	return now.After(this.upcoming)
}
func (this *ReportingFrequency) ScheduleNextMeasurement(now time.Time) {
	this.upcoming = now.Add(this.interval)
}

////////////////////////////////////////////////////////////////////////////

type SimpleMetric struct {
	*ReportingFrequency
	name  string
	value int64
}

func (this *SimpleMetric) Measure() MetricMeasurement {
	return MetricMeasurement{}
}

////////////////////////////////////////////////////////////////////////////
