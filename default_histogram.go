package metrics

import "sync/atomic"

type simpleHistogram struct {
	name        string
	description string
	labels      string
	buckets     []uint64
	indexes     map[uint64]int
	values      []uint64
	sum         *uint64
	count       *uint64
}

func NewHistogram(name string, options ...option) Histogram {
	config := configuration{Name: name}
	Options.apply(options...)(&config)

	indexes := make(map[uint64]int, len(config.Buckets))
	for index, bucket := range config.Buckets {
		indexes[bucket] = index // bucket
	}

	var sum, count uint64
	return &simpleHistogram{
		name:        config.Name,
		description: config.Description,
		labels:      config.RenderLabels(),
		buckets:     config.Buckets,
		indexes:     indexes,
		values:      make([]uint64, len(config.Buckets)),
		sum:         &sum,
		count:       &count,
	}
}

func (this *simpleHistogram) Type() string        { return "histogram" }
func (this *simpleHistogram) Name() string        { return this.name }
func (this *simpleHistogram) Description() string { return this.description }
func (this *simpleHistogram) Labels() string      { return this.labels }

func (this *simpleHistogram) Measure(value uint64) {
	atomic.AddUint64(this.count, 1)
	atomic.AddUint64(this.sum, value)

	for index, key := range this.buckets {
		if value <= key {
			atomic.AddUint64(&this.values[index], 1)
		}
	}
}
func (this *simpleHistogram) Count() uint64 { return atomic.LoadUint64(this.count) }
func (this *simpleHistogram) Sum() uint64   { return atomic.LoadUint64(this.sum) }

func (this *simpleHistogram) Buckets() []uint64 { return this.buckets }
func (this *simpleHistogram) Value(key uint64) uint64 {
	return atomic.LoadUint64(&this.values[this.indexes[key]])
}
