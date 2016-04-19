package metrics

import (
	"sync/atomic"
	"time"
)

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
	MetricType int
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

type AtomicMetric struct {
	*ReportingFrequency
	name       string
	value      int64
	metricType int
}

func (this *AtomicMetric) Measure() MetricMeasurement {
	return MetricMeasurement{
		Name:       this.name,
		Value:      atomic.LoadInt64(&this.value),
		MetricType: this.metricType,
	}
}

func (this *AtomicMetric) Add(delta int64) { atomic.AddInt64(&this.value, delta) }
func (this *AtomicMetric) Set(value int64) { atomic.StoreInt64(&this.value, value) }

////////////////////////////////////////////////////////////////////////////
