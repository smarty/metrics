package metrics

import "sync/atomic"

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
