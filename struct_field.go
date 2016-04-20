package metrics

// TODO: add histogram support

// Metrics can be used as a struct field and overridden with the Capture
// function in unit test setups to allow assertions on counted and measured
// values. This approach is similar to the one employed by the clock package
// (see github.com/smartystreets/clock).
type Metrics struct {
	all map[int]int64
}

func Capture() *Metrics {
	return &Metrics{
		all: make(map[int]int64),
	}
}

func (this *Metrics) Count(id CounterMetric) bool {
	if this == nil {
		return standard.Count(id)
	}
	return this.add(int(id), 1)
}

func (this *Metrics) CountN(id CounterMetric, increment int64) bool {
	if this == nil {
		return standard.CountN(id, increment)
	}
	return this.add(int(id), increment)
}

func (this *Metrics) RawCount(id CounterMetric, value int64) bool {
	if this == nil {
		return standard.RawCount(id, value)
	}
	return this.set(int(id), value)
}

func (this *Metrics) Measure(id GaugeMetric, value int64) bool {
	if this == nil {
		return standard.Measure(id, value)
	}
	return this.set(int(id), value)
}

func (this *Metrics) add(id int, increment int64) bool {
	this.all[id] += increment
	return true
}

func (this *Metrics) set(id int, value int64) bool {
	this.all[id] = value
	return true
}

// Helper functions for test assertions:
func (this *Metrics) CounterValue(id CounterMetric) int64 {
	if this == nil {
		return -1
	}
	return this.value(int(id))
}
func (this *Metrics) GaugeValue(id GaugeMetric) int64 {
	if this == nil {
		return -1
	}
	return this.value(int(id))
}
func (this *Metrics) value(id int) int64 { return this.all[id] }
