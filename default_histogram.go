package metrics

import "sync/atomic"

type simpleHistogram struct {
	name        string
	description string
	labels      string
	buckets     []uint64
	indexes     map[uint64]int
	values      []uint64
	sum         *atomic.Uint64
	count       *atomic.Uint64
}

func NewHistogram(name string, options ...option) Histogram {
	config := configuration{Name: name}
	Options.apply(options...)(&config)

	indexes := make(map[uint64]int, len(config.Buckets))
	for index, bucket := range config.Buckets {
		indexes[bucket] = index // bucket
	}

	this := &simpleHistogram{
		name:        config.Name,
		description: config.Description,
		labels:      config.RenderLabels(),
		buckets:     config.Buckets,
		indexes:     indexes,
		values:      make([]uint64, len(config.Buckets)),
		sum:         new(atomic.Uint64),
		count:       new(atomic.Uint64),
	}

	config.Exporter.Add(this)
	return this
}

func (this *simpleHistogram) Type() string        { return "histogram" }
func (this *simpleHistogram) Name() string        { return this.name }
func (this *simpleHistogram) Description() string { return this.description }
func (this *simpleHistogram) Labels() string      { return this.labels }

func (this *simpleHistogram) Measure(value uint64) {
	this.count.Add(1)
	this.sum.Add(value)

	for index, key := range this.buckets {
		if value <= key {
			atomic.AddUint64(&this.values[index], 1)
		}
	}
}

func (this *simpleHistogram) Buckets() []uint64 { return this.buckets }
func (this *simpleHistogram) Value(key uint64) uint64 {
	return atomic.LoadUint64(&this.values[this.indexes[key]])
}
func (this *simpleHistogram) Count() uint64 { return this.count.Load() }
func (this *simpleHistogram) Sum() uint64   { return this.sum.Load() }
