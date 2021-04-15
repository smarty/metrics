package metrics

import "net/http"

type Metric interface {
	Type() string
	Name() string
	Description() string
	Labels() string
}

type Counter interface {
	Metric

	Increment()
	IncrementN(uint64)

	Value() uint64
}

type Gauge interface {
	Metric

	Increment()
	IncrementN(int64)
	Measure(int64)

	Value() int64
}

type Histogram interface {
	Metric

	Measure(uint64)
	Count() uint64
	Sum() uint64

	Buckets() []uint64
	Value(bucket uint64) uint64
}

type Exporter interface {
	http.Handler

	Add(...Metric)
}
