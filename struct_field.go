package metrics

import "github.com/smartystreets/metrics/internal/hdrhistogram"

// Metrics can be used as a struct field and overridden with the Capture
// function in unit test setups to allow assertions on counted and measured
// values. This approach is similar to the one employed by the clock package
// (see github.com/smartystreets/clock).
type Metrics struct {
	counters   map[CounterMetric]int64
	gauges     map[GaugeMetric]int64
	histograms map[HistogramMetric]*hdrhistogram.Histogram
}

func Capture() *Metrics {
	return &Metrics{
		counters:   make(map[CounterMetric]int64),
		gauges:     make(map[GaugeMetric]int64),
		histograms: make(map[HistogramMetric]*hdrhistogram.Histogram),
	}
}

func (this *Metrics) Count(id CounterMetric) bool {
	return this.CountN(id, 1)
}

func (this *Metrics) CountIf(id CounterMetric, condition bool) bool {
	if condition {
		return this.Count(id)
	}
	return condition
}

func (this *Metrics) CountN(id CounterMetric, increment int64) bool {
	if this == nil {
		return standard.CountN(id, increment)
	}
	this.counters[id] += increment
	return true
}

func (this *Metrics) RawCount(id CounterMetric, value int64) bool {
	if this == nil {
		return standard.RawCount(id, value)
	}
	this.counters[id] = value
	return true
}

func (this *Metrics) Measure(id GaugeMetric, value int64) bool {
	if this == nil {
		return standard.Measure(id, value)
	}
	this.gauges[id] = value
	return true
}

func (this *Metrics) Record(id HistogramMetric, value int64) bool {
	if this == nil {
		return standard.Record(id, value)
	}
	histogram := this.histograms[id]
	if histogram == nil {
		histogram = hdrhistogram.New(0, max, resolution)
		this.histograms[id] = histogram
	}
	return histogram.RecordValue(value) == nil
}

// Helper functions for test assertions:
func (this *Metrics) CounterValue(id CounterMetric) int64 {
	return this.counters[id]
}
func (this *Metrics) GaugeValue(id GaugeMetric) int64 {
	return this.gauges[id]
}
func (this *Metrics) HistogramValue(id HistogramMetric) Histogram {
	return this.histograms[id]
}

const (
	// int64 max value causes hdrhistogram to hang. 1 billion is probably high enough for most scenarios:
	max        = 1000000000
	resolution = 5
)
