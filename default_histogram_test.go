package metrics

import (
	"sync"
	"testing"
)

func TestHistogramValues(t *testing.T) {
	assertHistogram(t,
		[]uint64{0, 1, 20, 30, 50, 100, 300, 500},
		[]uint64{0, 1, 5, 5, 6, 7, 9, 9},
		1023,
		10)
}
func assertHistogram(t *testing.T, buckets, values []uint64, sum, count uint64) {
	var options []option
	for _, bucket := range buckets {
		options = append(options, Options.Bucket(bucket))
	}
	histogram := NewHistogram("", options...)
	measureHistogram(histogram)

	for index, key := range histogram.Buckets() {
		assertEqual(t, buckets[index], key)
		assertEqual(t, values[index], histogram.Value(key))
	}
	assertEqual(t, sum, histogram.Sum())
	assertEqual(t, count, histogram.Count())
}
func measureHistogram(metric Histogram) {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	for x := uint64(1); x < 1000; x = x * 2 {
		wg.Add(1)
		go func(measurement uint64) {
			metric.Measure(measurement)
			wg.Done()
		}(x)
	}
}
