package metrics

import "sync/atomic"

type defaultCounter struct {
	name        string
	description string
	labels      string
	value       *atomic.Uint64
}

func NewCounter(name string, options ...option) Counter {
	config := configuration{Name: name}
	Options.apply(options...)(&config)

	this := &defaultCounter{
		name:        config.Name,
		description: config.Description,
		labels:      config.RenderLabels(),
		value:       new(atomic.Uint64),
	}

	config.Exporter.Add(this)
	return this
}

func (this *defaultCounter) Type() string        { return "counter" }
func (this *defaultCounter) Name() string        { return this.name }
func (this *defaultCounter) Description() string { return this.description }
func (this *defaultCounter) Labels() string      { return this.labels }

func (this *defaultCounter) Increment()              { this.value.Add(1) }
func (this *defaultCounter) IncrementN(value uint64) { this.value.Add(value) }

func (this *defaultCounter) Value() uint64 { return this.value.Load() }
