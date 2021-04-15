package metrics

import "sync/atomic"

type simpleGauge struct {
	name        string
	description string
	labels      string
	value       *int64
}

func NewGauge(name string, options ...option) Gauge {
	config := configuration{Name: name}
	Options.apply(options...)(&config)
	var value int64
	return &simpleGauge{
		name:        config.Name,
		description: config.Description,
		labels:      config.RenderLabels(),
		value:       &value,
	}
}

func (this *simpleGauge) Type() string        { return "gauge" }
func (this *simpleGauge) Name() string        { return this.name }
func (this *simpleGauge) Description() string { return this.description }
func (this *simpleGauge) Labels() string      { return this.labels }

func (this *simpleGauge) Keys() []int64       { return defaultKeys }
func (this *simpleGauge) Value(_ int64) int64 { return atomic.LoadInt64(this.value) }

func (this *simpleGauge) Increment()             { atomic.AddInt64(this.value, 1) }
func (this *simpleGauge) IncrementN(value int64) { atomic.AddInt64(this.value, value) }
func (this *simpleGauge) Measure(value int64)    { atomic.StoreInt64(this.value, value) }
