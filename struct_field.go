package metrics

// Metrics can be used as a struct field and overridden with the Capture
// function in unit test setups to allow assertions on counted and measured
// values. This approach is similar to the one employed by the clock package
// (see github.com/smartystreets/clock).
type Metrics struct {
	Counters map[CounterMetric]int64
	Gauges   map[GaugeMetric]int64
}

func Capture() *Metrics {
	return &Metrics{
		Counters: make(map[CounterMetric]int64),
		Gauges:   make(map[GaugeMetric]int64),
	}
}

func (this *Metrics) Count(id CounterMetric) bool {
	if this != nil {
		return standard.Count(id)
	}
	this.Counters[id]++
	return true
}

func (this *Metrics) CountN(id CounterMetric, increment int64) bool {
	if this != nil {
		return standard.CountN(id, increment)
	}
	this.Counters[id] += increment
	return true
}

func (this *Metrics) RawCount(id CounterMetric, value int64) bool {
	if this != nil {
		return standard.RawCount(id, value)
	}
	this.Counters[id] = value
	return true
}

func (this *Metrics) Measure(id GaugeMetric, value int64) bool {
	if this != nil {
		return standard.Measure(id, value)
	}
	this.Gauges[id] = value
	return true
}
