package metrics

import "time"

type MetricsTracker struct {
	metrics    []Metric
	counters   map[CounterMetric]*AtomicMetric
	gauges     map[GaugeMetric]*AtomicMetric
	histograms map[HistogramMetric]Histogram

	started int
}

func New2() *MetricsTracker { // TODO: rename to New() when we finish (get rid of *container)
	return &MetricsTracker{
		counters:   make(map[CounterMetric]*AtomicMetric),
		gauges:     make(map[GaugeMetric]*AtomicMetric),
		histograms: make(map[HistogramMetric]Histogram),
	}
}

func (this *MetricsTracker) AddCounter(name string, update time.Duration) CounterMetric {
	metric := &AtomicMetric{
		ReportingFrequency: &ReportingFrequency{interval: update},
		name:               name,
		metricType:         counterMetricType,
	}
	id := CounterMetric(len(this.metrics))
	this.metrics = append(this.metrics, metric)
	this.counters[id] = metric
	return id
}
func (this *MetricsTracker) AddGauge(name string, update time.Duration) GaugeMetric {
	return -1
}
func (this *MetricsTracker) AddHistogram(name string, update time.Duration,
	min, max int64, resolution int, quantiles ...float64) HistogramMetric {

	return -1
}

func (this *MetricsTracker) StartMeasuring() {}
func (this *MetricsTracker) StopMeasuring()  {}

func (this *MetricsTracker) Count(id CounterMetric) bool {
	// TODO: validate id against len of this.metrics
	this.counters[id].Add(1)
	return true
}
func (this *MetricsTracker) CountN(id CounterMetric, n int64) bool       { return false }
func (this *MetricsTracker) CountRaw(id CounterMetric, raw int64) bool   { return false }
func (this *MetricsTracker) Measure(id GaugeMetric, value int64) bool    { return false }
func (this *MetricsTracker) Record(id HistogramMetric, value int64) bool { return false }

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
