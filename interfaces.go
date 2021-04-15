package metrics

import "net/http"

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
	Buckets() []uint64
	Values() []uint64
	Count() uint64
	Sum() uint64
}

type Exporter interface {
	http.Handler

	Add(...Metric)
}
