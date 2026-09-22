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

	Buckets() []uint64
	Value(bucket uint64) uint64
	Count() uint64
	Sum() uint64
}

// Exporter renders registered metrics in the Prometheus text exposition format.
// Implementations are safe for concurrent use: Add may be called from multiple
// goroutines, and concurrently with ServeHTTP.
type Exporter interface {
	http.Handler

	Add(...Metric)
}
