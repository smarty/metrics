package metrics

import (
	"fmt"
	"strings"
	"sync/atomic"
)

type simpleCounter struct {
	name        string
	description string
	labels      string
	value       *int64
}

func NewCounter(name string, options ...option) Counter {
	config := configuration{Name: name}
	Options.apply(options...)(&config)
	var value int64
	return &simpleCounter{
		name:        config.Name,
		description: config.Description,
		labels:      config.RenderLabels(),
		value:       &value,
	}
}

func (this *simpleCounter) Type() string        { return "counter" }
func (this *simpleCounter) Name() string        { return this.name }
func (this *simpleCounter) Description() string { return this.description }
func (this *simpleCounter) Labels() string      { return this.labels }

func (this *simpleCounter) Keys() []int64       { return defaultKeys }
func (this *simpleCounter) Value(_ int64) int64 { return atomic.LoadInt64(this.value) }

func (this *simpleCounter) Increment()              { atomic.AddInt64(this.value, 1) }
func (this *simpleCounter) IncrementN(value uint64) { atomic.AddInt64(this.value, int64(value)) }

var defaultKeys = []int64{0}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type simpleHistogram struct {
	name        string
	description string
	labels      string
	keys        []int64
	indexes     map[int64]int
	values      []int64
	sum         *int64
	count       *int64
}

func NewHistogram(name string, options ...option) Histogram {
	config := configuration{Name: name}
	Options.apply(options...)(&config)

	indexes := make(map[int64]int, len(config.Keys))
	for index, key := range config.Keys {
		indexes[key] = index
	}

	var sum, count int64
	return &simpleHistogram{
		name:        config.Name,
		description: config.Description,
		labels:      config.RenderLabels(),
		keys:        config.Keys,
		indexes:     indexes,
		values:      make([]int64, len(config.Keys)),
		sum:         &sum,
		count:       &count,
	}
}

func (this *simpleHistogram) Type() string        { return "histogram" }
func (this *simpleHistogram) Name() string        { return this.name }
func (this *simpleHistogram) Description() string { return this.description }
func (this *simpleHistogram) Labels() string      { return this.labels }

func (this *simpleHistogram) Keys() []int64 { return this.keys }
func (this *simpleHistogram) Value(key int64) int64 {
	return atomic.LoadInt64(&this.values[this.indexes[key]])
}

func (this *simpleHistogram) Measure(value int64) {
	for index, key := range this.keys {
		if value <= key {
			atomic.AddInt64(&this.values[index], 1)
		}
	}

	atomic.AddInt64(this.sum, value)
	atomic.AddInt64(this.count, 1)
}
func (this *simpleHistogram) Count() int64 { return atomic.LoadInt64(this.count) }
func (this *simpleHistogram) Sum() int64   { return atomic.LoadInt64(this.sum) }

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var Options singleton

type singleton struct{}
type option func(*configuration)
type configuration struct {
	Name        string
	Description string
	Labels      map[string]string
	Keys        []int64
}

func (singleton) Description(value string) option {
	return func(this *configuration) { this.Description = value }
}
func (singleton) Label(key, value string) option {
	return func(this *configuration) { this.Labels[key] = value }
}
func (singleton) Bucket(value int64) option {
	return func(this *configuration) {
		this.Keys = append(this.Keys, value)
	}
}
func (singleton) apply(options ...option) option {
	return func(this *configuration) {
		this.Labels = map[string]string{}
		for _, option := range Options.defaults(options...) {
			option(this)
		}
	}
}
func (singleton) defaults(options ...option) []option {
	return append([]option{}, options...)
}

func (this configuration) RenderLabels() (result string) {
	if len(this.Labels) == 0 {
		return ""
	}

	for key, value := range this.Labels {
		result += fmt.Sprintf(`%s="%s", `, key, value)
	}
	result = strings.TrimSuffix(result, ", ")
	return fmt.Sprintf("{ %s }", result)
}
