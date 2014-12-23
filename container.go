package metrics

import (
	"sync/atomic"
	"time"
)

var standard = New()

type container struct {
	metrics []int64
	meta    []metricInfo
	started int32
	queue   chan []Measurement
}

func New() *container {
	return &container{}
}
func (this *container) Add(name string, reportingFrequency time.Duration) int {
	if atomic.LoadInt32(&this.started) > 0 {
		return MetricConflict
	}

	if int64(reportingFrequency) <= 0 {
		return MetricConflict
	}

	for _, metric := range this.meta {
		if metric.Name == name {
			return MetricConflict
		}
	}

	this.metrics = append(this.metrics, int64(0))
	info := metricInfo{Name: name, ReportingFrequency: reportingFrequency}
	this.meta = append(this.meta, info)
	return len(this.metrics) - 1
}

func (this *container) StartMeasuring() {
	if atomic.AddInt32(&this.started, 1) > 1 {
		return
	}

	durations := map[time.Duration][]int{}
	for i, item := range this.meta {
		indices := durations[item.ReportingFrequency]
		indices = append(indices, i)
		durations[item.ReportingFrequency] = indices
	}

	for d, i := range durations {
		duration := d // save the values for
		indices := i  // the closure below...
		time.AfterFunc(duration, func() { this.report(duration, indices) })
	}
}
func (this *container) StopMeasuring() {
	atomic.SwapInt32(&this.started, 0)
}

func (this *container) registerChannelDestination(queue chan []Measurement) {
	this.queue = queue
}
func (this *container) report(duration time.Duration, indices []int) {
	now := time.Now()
	snapshot := make([]Measurement, len(indices), len(indices))

	for i := 0; i < len(indices); i++ {
		index := indices[i]
		snapshot[i] = Measurement{
			ID:       index,
			Captured: now,
			Value:    atomic.LoadInt64(&this.metrics[index]),
		}
	}

	this.queue <- snapshot

	if atomic.LoadInt32(&this.started) > 0 {
		time.AfterFunc(duration, func() { this.report(duration, indices) })
	}
}

func (this *container) Count(index int) bool {
	if index < 0 || len(this.metrics) <= index {
		return false
	}

	atomic.AddInt64(&this.metrics[index], 1)
	return true
}
func (this *container) Measure(index int, measurement int64) bool {
	if index < 0 || len(this.metrics) <= index {
		return false
	}

	atomic.StoreInt64(&this.metrics[index], measurement)
	return true
}
