package metrics2

import "net/http"

type metric interface {
	Type() string
	Name() string
	Description() string
	Labels() string

	Value() int64
	Increment()
}

type Counter interface {
	metric
	IncrementN(uint64)
}

type Gauge interface {
	metric
	IncrementN(int64)
	Measure(int64)
}

type Exporter interface {
	http.Handler

	Add(...metric)
}
