package metrics

import "net/http"

type Metric interface {
	Type() string
	Name() string
	Description() string
	Labels() string

	Keys() []int64
	Value(key int64) int64
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

	Measure(int64)
	Count() int64
	Sum() int64
}

type Exporter interface {
	http.Handler

	Add(...Metric)
}
