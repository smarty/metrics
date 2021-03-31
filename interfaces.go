package metrics

import "net/http"

type Metric interface {
	Type() string
	Name() string
	Description() string
	Labels() string

	Value() int64
	Increment()
}

type Counter interface {
	Metric
	IncrementN(uint64)
}

type Gauge interface {
	Metric
	IncrementN(int64)
	Measure(int64)
}

type Histogram interface {
	Metric
	Buckets() map[float64]uint64
	Count() uint64
	Sum() float64
	Observe(float64)
}

type Exporter interface {
	http.Handler

	Add(...Metric)
}
