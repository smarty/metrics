package metrics

import "sync/atomic"

type defaultCounter struct {
	name        string
	description string
	labels      string
	value       *int64
}

func NewCounter(name string, options ...option) Counter {
	config := configuration{Name: name}
	Options.apply(options...)(&config)
	var value int64
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

func (this *defaultCounter) Keys() []int64       { return defaultKeys }
func (this *defaultCounter) Value(_ int64) int64 { return atomic.LoadInt64(this.value) }

func (this *defaultCounter) Increment()              { atomic.AddInt64(this.value, 1) }
func (this *defaultCounter) IncrementN(value uint64) { atomic.AddInt64(this.value, int64(value)) }

var defaultKeys = []int64{0}
