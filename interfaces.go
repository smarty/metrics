package metrics

import (
	"net/http"
	"sync/atomic"
)

type Metric interface {
	Type() string
	Name() string
	Description() string
	Labels() string

	Value() int64
}

type Counter interface {
	Metric
	Increment()
	IncrementN(uint64)
}

type Gauge interface {
	Metric
	Increment()
	IncrementN(int64)
	Measure(int64)
}

type Histogram interface {
	Metric
	Measure(uint64)
	Buckets() []Bucket
	Count() uint64
	Sum() uint64
}

type Bucket struct{ key, value uint64 }

func (this Bucket) Key() uint64   { return this.key }
func (this Bucket) Value() uint64 { return atomic.LoadUint64(&this.value) }

type Exporter interface {
	http.Handler

	Add(...Metric)
}
