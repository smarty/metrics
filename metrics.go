package metrics

import (
	"sync"
	"time"

	"github.com/smartystreets/metrics/internal/hdrhistogram"
)

type MetricsTracker struct {
	metrics    []Metric
	counters   map[CounterMetric]*AtomicMetric
	gauges     map[GaugeMetric]*AtomicMetric
	histograms map[HistogramMetric]Histogram
}

func New2() *MetricsTracker { // TODO: rename to New() when we finish (get rid of *container)
	return &MetricsTracker{
		counters:   make(map[CounterMetric]*AtomicMetric),
		gauges:     make(map[GaugeMetric]*AtomicMetric),
		histograms: make(map[HistogramMetric]Histogram),
	}
}

// TODO: guard against negative durations?
// TODO: guard against blank names?
// TODO: guard against min > max (histogram)?
// TODO: guard against resolution being below 1 or above 5 (histogram)? (to prevent a panic)

func (this *MetricsTracker) AddCounter(name string, update time.Duration) CounterMetric {
	metric := NewCounter(name, update)
	id := CounterMetric(len(this.metrics))
	this.metrics = append(this.metrics, metric)
	this.counters[id] = metric
	return id
}
func (this *MetricsTracker) AddGauge(name string, update time.Duration) GaugeMetric {
	metric := NewGauge(name, update)
	id := GaugeMetric(len(this.metrics))
	this.metrics = append(this.metrics, metric)
	this.gauges[id] = metric
	return id
}
func (this *MetricsTracker) AddHistogram(name string, update time.Duration,
	min, max int64, resolution int, quantiles ...float64) HistogramMetric {

	id, histogram := this.addHistogram(min, max, resolution)
	this.addHistogramMetrics(name, update, histogram, quantiles)
	return id
}

func (this *MetricsTracker) addHistogram(min, max int64, resolution int) (HistogramMetric, Histogram) {
	mutex := new(sync.RWMutex)
	id := HistogramMetric(len(this.histograms))
	histogram := hdrhistogram.New(min, max, resolution)
	synchronized := NewSynchronizedHistogram(histogram, mutex.RLocker(), mutex)
	this.histograms[id] = synchronized
	return id, synchronized
}
func (this *MetricsTracker) addHistogramMetrics(
	name string, update time.Duration, histogram Histogram, quantiles []float64) {

	this.metrics = append(this.metrics, NewHistogramMinMetric(name, histogram, update))
	this.metrics = append(this.metrics, NewHistogramMaxMetric(name, histogram, update))
	this.metrics = append(this.metrics, NewHistogramMeanMetric(name, histogram, update))
	this.metrics = append(this.metrics, NewHistogramStandardDeviationMetric(name, histogram, update))
	this.metrics = append(this.metrics, NewHistogramTotalCountMetric(name, histogram, update))

	for _, quantile := range quantiles {
		this.metrics = append(this.metrics, NewHistogramQuantileMetric(name, quantile, histogram, update))
	}
}

func (this *MetricsTracker) StartMeasuring() {

}
func (this *MetricsTracker) StopMeasuring() {
	panic("TODO")
}

func (this *MetricsTracker) Count(id CounterMetric) bool {
	return this.CountN(id, 1)
}
func (this *MetricsTracker) CountN(id CounterMetric, n int64) bool {
	counter, found := this.counters[id]
	if found {
		counter.Add(n)
	}
	return found
}
func (this *MetricsTracker) CountRaw(id CounterMetric, value int64) bool {
	counter, found := this.counters[id]
	if found {
		counter.Set(value)
	}
	return found
}
func (this *MetricsTracker) Measure(id GaugeMetric, value int64) bool {
	gauge, found := this.gauges[id]
	if found {
		gauge.Set(value)
	}
	return found
}
func (this *MetricsTracker) Record(id HistogramMetric, value int64) bool {
	histogram, found := this.histograms[id]
	if found {
		histogram.RecordValue(value)
	}
	return found
}

func (this *MetricsTracker) TakeMeasurements(now time.Time) (measurements []MetricMeasurement) {
	for id, metric := range this.metrics {
		if metric.MeasurementIsOverdue(now) {
			metric.ScheduleNextMeasurement(now)
			measurement := metric.Measure()
			measurement.ID = id
			measurement.Captured = now
			measurements = append(measurements, measurement)
		}
	}
	return measurements
}
