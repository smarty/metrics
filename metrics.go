package metrics

import (
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/smartystreets/metrics/internal/hdrhistogram"
)

type MetricsTracker struct {
	metrics    []Metric
	counters   map[CounterMetric]*AtomicMetric
	gauges     map[GaugeMetric]*AtomicMetric
	histograms map[HistogramMetric]Histogram
	started    int32
}

func New() *MetricsTracker {
	return &MetricsTracker{
		counters:   make(map[CounterMetric]*AtomicMetric),
		gauges:     make(map[GaugeMetric]*AtomicMetric),
		histograms: make(map[HistogramMetric]Histogram),
	}
}

func (this *MetricsTracker) AddCounter(name string, update time.Duration) CounterMetric {
	if name = cleanName(name); !this.canAddMetric(name, update) {
		return MetricConflict
	}
	metric := NewCounter(name, update)
	id := CounterMetric(len(this.metrics))
	this.metrics = append(this.metrics, metric)
	this.counters[id] = metric
	return id
}
func (this *MetricsTracker) AddGauge(name string, update time.Duration) GaugeMetric {
	if name = cleanName(name); !this.canAddMetric(name, update) {
		return MetricConflict
	}
	metric := NewGauge(name, update)
	id := GaugeMetric(len(this.metrics))
	this.metrics = append(this.metrics, metric)
	this.gauges[id] = metric
	return id
}
func (this *MetricsTracker) AddHistogram(name string, update time.Duration,
	min, max int64, resolution int, quantiles ...float64) HistogramMetric {

	if name = cleanName(name); !this.canAddMetric(name, update) {
		return MetricConflict
	}
	if !histogramOptionsAreValid(min, max, resolution) {
		return MetricConflict
	}
	id, histogram := this.addHistogram(min, max, resolution)
	this.addHistogramMetrics(name, update, histogram, quantiles)
	return id
}
func histogramOptionsAreValid(min, max int64, resolution int) bool {
	return min < max && (resolution >= 1 && resolution <= 5)
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

func cleanName(name string) string {
	return strings.TrimSpace(name)
}

func (this *MetricsTracker) canAddMetric(name string, duration time.Duration) bool {
	if atomic.LoadInt32(&this.started) > 0 {
		return false
	}
	if len(name) == 0 {
		return false
	}
	if this.nameIsTaken(name) {
		return false
	}
	if duration < time.Nanosecond {
		return false
	}
	return true
}

func (this *MetricsTracker) nameIsTaken(name string) bool {
	for _, metric := range this.metrics {
		if metric.Measure().Name == name {
			return true
		}
	}
	return false
}

func (this *MetricsTracker) StartMeasuring() {
	if !this.isRunning() {
		atomic.AddInt32(&this.started, 1)
	}
}
func (this *MetricsTracker) StopMeasuring() {
	atomic.StoreInt32(&this.started, 0)
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
func (this *MetricsTracker) RawCount(id CounterMetric, value int64) bool {
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
		found = histogram.RecordValue(value) == nil
	}
	return found
}

func (this *MetricsTracker) TakeMeasurements(now time.Time) []MetricMeasurement {
	if !this.isRunning() {
		return nil
	}
	measurements := []MetricMeasurement{}
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

func (this *MetricsTracker) isRunning() bool {
	return atomic.LoadInt32(&this.started) > 0
}
