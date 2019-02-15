package metrics

import (
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/smartystreets/logging"
	"github.com/smartystreets/metrics/internal/hdrhistogram"
)

type MetricsTracker struct {
	logger *logging.Logger

	metrics map[int]Metric

	counters   map[CounterMetric]*AtomicMetric
	gauges     map[GaugeMetric]*AtomicMetric
	histograms map[HistogramMetric]Histogram

	histogramMetrics map[int]HistogramMetric
	tags             map[int]map[string]string

	started int32
}

func New() *MetricsTracker {
	return &MetricsTracker{
		metrics:          make(map[int]Metric),
		counters:         make(map[CounterMetric]*AtomicMetric),
		gauges:           make(map[GaugeMetric]*AtomicMetric),
		histograms:       make(map[HistogramMetric]Histogram),
		histogramMetrics: make(map[int]HistogramMetric),
		tags:             make(map[int]map[string]string),
	}
}

func (this *MetricsTracker) nextID() int {
	return len(this.metrics) + len(this.histograms)
}

func (this *MetricsTracker) AddCounter(name string, update time.Duration) CounterMetric {
	if name = cleanName(name); !this.canAddMetric(name, update) {
		return MetricConflict
	}
	metric := NewCounter(name, update)
	id := CounterMetric(this.nextID())
	this.metrics[int(id)] = metric
	this.counters[id] = metric
	return id
}
func (this *MetricsTracker) AddGauge(name string, update time.Duration) GaugeMetric {
	if name = cleanName(name); !this.canAddMetric(name, update) {
		return MetricConflict
	}
	metric := NewGauge(name, update)
	id := GaugeMetric(this.nextID())
	this.metrics[int(id)] = metric
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
	this.addHistogramMetrics(id, name, update, histogram, quantiles)
	return id
}
func histogramOptionsAreValid(min, max int64, resolution int) bool {
	return min < max && (resolution >= 1 && resolution <= 5)
}
func (this *MetricsTracker) addHistogram(min, max int64, resolution int) (HistogramMetric, Histogram) {
	mutex := new(sync.RWMutex)
	id := HistogramMetric(this.nextID())
	histogram := hdrhistogram.New(min, max, resolution)
	synchronized := NewSynchronizedHistogram(histogram, mutex.RLocker(), mutex)
	this.histograms[id] = synchronized
	return id, synchronized
}
func (this *MetricsTracker) addHistogramMetrics(id HistogramMetric,
	name string, update time.Duration, histogram Histogram, quantiles []float64) {

	var all = []Metric{
		NewHistogramMinMetric(name, histogram, update),
		NewHistogramMaxMetric(name, histogram, update),
		NewHistogramMeanMetric(name, histogram, update),
		NewHistogramStandardDeviationMetric(name, histogram, update),
		NewHistogramTotalCountMetric(name, histogram, update),
	}
	for _, quantile := range quantiles {
		all = append(all, NewHistogramQuantileMetric(name, quantile, histogram, update))
	}

	for _, metric := range all {
		next := this.nextID()
		this.metrics[next] = metric
		this.histogramMetrics[next] = id
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

func (this *MetricsTracker) TakeMeasurements(now time.Time) (measurements []MetricMeasurement) {
	if !this.isRunning() {
		return nil
	}
	for id, metric := range this.metrics {
		if metric.MeasurementIsOverdue(now) {
			metric.ScheduleNextMeasurement(now)
			measurement := metric.Measure()
			measurement.ID = id
			measurement.Captured = now
			if histogram, ok := this.histogramMetrics[id]; ok {
				measurement.Tags = this.tags[int(histogram)]
			} else {
				measurement.Tags = this.tags[id]
			}
			measurements = append(measurements, measurement)
		}
	}
	return measurements
}

func (this *MetricsTracker) isRunning() bool {
	return atomic.LoadInt32(&this.started) > 0
}

func (this *MetricsTracker) TagCounter(metric CounterMetric, tagPairs ...string) {
	this.addTags(int(metric), tagPairs)
}
func (this *MetricsTracker) TagGauge(metric GaugeMetric, tagPairs ...string) {
	this.addTags(int(metric), tagPairs)
}
func (this *MetricsTracker) TagHistogram(metric HistogramMetric, tagPairs ...string) {
	this.addTags(int(metric), tagPairs)
}
func (this *MetricsTracker) addTags(metric int, tagPairs []string) {
	if len(tagPairs)%2 > 0 {
		this.logger.Printf("[WARN] tags must be submitted as an even number of key/value pairs. You provided %d values.", len(tagPairs))
		return
	}
	this.tags[metric] = map[string]string{}
	for i := 0; i < len(tagPairs); i += 2 {
		this.tags[int(metric)][tagPairs[i+0]] = tagPairs[i+1]
	}
}
