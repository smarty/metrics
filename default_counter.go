package metrics

import "sync/atomic"

type defaultCounter struct {
	name        string
	description string
	labels      string
	value       *uint64
}

func NewCounter(name string, options ...option) Counter {
	config := configuration{Name: name}
	Options.apply(options...)(&config)
	var value uint64
	return &defaultCounter{
		name:        config.Name,
		description: config.Description,
		labels:      config.RenderLabels(),
		value:       &value,
	}
}

func (this *defaultCounter) Type() string        { return "counter" }
func (this *defaultCounter) Name() string        { return this.name }
func (this *defaultCounter) Description() string { return this.description }
func (this *defaultCounter) Labels() string      { return this.labels }

func (this *defaultCounter) Increment()              { atomic.AddUint64(this.value, 1) }
func (this *defaultCounter) IncrementN(value uint64) { atomic.AddUint64(this.value, value) }

func (this *defaultCounter) Value() uint64 { return atomic.LoadUint64(this.value) }
