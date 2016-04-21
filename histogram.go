package metrics

type Histogram interface {
	RecordValue(v int64) error
	Min() int64
	Max() int64
	Mean() float64
	StdDev() float64
	TotalCount() int64
	ValueAtQuantile(q float64) int64
}
