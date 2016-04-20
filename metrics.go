package metrics

import "time"

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

	panic("TODO")
}

func (this *MetricsTracker) StartMeasuring() {

}
func (this *MetricsTracker) StopMeasuring()  {
	panic("TODO")
}

func (this *MetricsTracker) Count(id CounterMetric) bool {
	return this.CountN(id, 1)
}
func (this *MetricsTracker) CountN(id CounterMetric, n int64) bool     {
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
	panic("TODO")
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
