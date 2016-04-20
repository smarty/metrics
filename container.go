package metrics

import (
	"sync/atomic"
	"time"
)

type container struct {
	metrics []int64
	meta    []metricInfo
	started int32
	queue   chan []Measurement
}

func (this *container) AddCounter(name string, reportingFrequency time.Duration) CounterMetric {
	return CounterMetric(this.add(name, counterMetricType, reportingFrequency))
}

func (this *container) AddGauge(name string, reportingFrequency time.Duration) GaugeMetric {
	return GaugeMetric(this.add(name, gaugeMetricType, reportingFrequency))
}

func (this *container) add(name string, metricType uint8, reportingFrequency time.Duration) int {
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
	info := metricInfo{Name: name, MetricType: metricType, ReportingFrequency: reportingFrequency}
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

func (this *container) RegisterChannelDestination(queue chan []Measurement) {
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

func (this *container) Count(id CounterMetric) bool {
	return this.CountN(id, 1)
}
func (this *container) CountN(id CounterMetric, increment int64) bool {
	index := int(id)
	if index < 0 || len(this.metrics) <= index {
		return false
	}

	atomic.AddInt64(&this.metrics[index], increment)
	return true
}
func (this *container) RawCount(id CounterMetric, measurement int64) bool {
	return this.measure(int(id), measurement)
}

func (this *container) Measure(id GaugeMetric, measurement int64) bool {
	return this.measure(int(id), measurement)
}

func (this *container) measure(index int, measurement int64) bool {
	if index < 0 || len(this.metrics) <= index {
		return false
	}

	atomic.StoreInt64(&this.metrics[index], measurement)
	return true
}
