package metrics

import "sync/atomic"

type simpleGauge struct {
	name        string
	description string
	labels      string
	value       *atomic.Int64
}

func NewGauge(name string, options ...option) Gauge {
	config := configuration{Name: name}
	Options.apply(options...)(&config)

	this := &simpleGauge{
		name:        config.Name,
		description: config.Description,
		labels:      config.RenderLabels(),
		value:       new(atomic.Int64),
	}

	config.Exporter.Add(this)
	return this
}

func (this *simpleGauge) Type() string        { return "gauge" }
func (this *simpleGauge) Name() string        { return this.name }
func (this *simpleGauge) Description() string { return this.description }
func (this *simpleGauge) Labels() string      { return this.labels }

func (this *simpleGauge) Increment()             { this.value.Add(1) }
func (this *simpleGauge) IncrementN(value int64) { this.value.Add(value) }
func (this *simpleGauge) Measure(value int64)    { this.value.Store(value) }

func (this *simpleGauge) Value() int64 { return this.value.Load() }
