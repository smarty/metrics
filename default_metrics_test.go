package metrics

import (
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func TestMetricsValues(t *testing.T) {
	histogramBucketKeys := []int64{0, 1, 20, 30, 50, 100, 300, 500}
	sample := newSampleMetrics(histogramBucketKeys)

	sample.counter1.Increment()
	sample.counter2.IncrementN(2)
	sample.gauge1.Increment()
	sample.gauge1.IncrementN(2)
	sample.gauge2.Measure(4)

	assertEqual(t, int64(1), sample.counter1.Value(0))
	assertEqual(t, int64(2), sample.counter2.Value(0))
	assertEqual(t, int64(3), sample.gauge1.Value(0))
	assertEqual(t, int64(4), sample.gauge2.Value(0))

	measureHistograms(sample)

	expectedHistogramValues1 := []int64{0, 1, 5, 5, 6, 7, 9, 9}
	for index, key := range sample.histogram1.Keys() {
		assertEqual(t, histogramBucketKeys[index], key)
		assertEqual(t, expectedHistogramValues1[index], sample.histogram1.Value(key))
	}
	assertEqual(t, int64(10), sample.histogram1.Count())
	assertEqual(t, int64(1023), sample.histogram1.Sum())

	expectedHistogramValues2 := []int64{0, 1, 3, 4, 4, 5, 6, 6}
	for index, key := range sample.histogram2.Keys() {
		assertEqual(t, histogramBucketKeys[index], key)
		assertEqual(t, expectedHistogramValues2[index], sample.histogram2.Value(key))
	}
	assertEqual(t, int64(7), sample.histogram2.Count())
	assertEqual(t, int64(1093), sample.histogram2.Sum())
}
func newSampleMetrics(bucketKeys []int64) *sampleMetrics {
	histogramOptions := []option{
		Options.Description("histogram description"),
		Options.Label("histogram_key1", "histogram_value1"),
	}
	for _, bucket := range bucketKeys {
		histogramOptions = append(histogramOptions, Options.Bucket(bucket))
	}

	return &sampleMetrics{
		counter1: NewCounter("my_counter",
			Options.Description("counter description"),
		),
		counter2: NewCounter("my_counter_with_labels",
			Options.Description("counter description"),
			Options.Label("counter_label_key", "counter_label_value"),
		),
		gauge1: NewGauge("my_gauge",
			Options.Description("gauge description"),
		),
		gauge2: NewGauge("my_gauge_with_labels",
			Options.Description("gauge description"),
			Options.Label("gauge_label_key", "gauge_label_value"),
		),
		histogram1: NewHistogram("my_histogram_with_buckets", histogramOptions...),
		histogram2: NewHistogram("my_histogram_with_buckets_and_labels", histogramOptions...),
	}
}
func assertEqual(t *testing.T, expected, actual interface{}) {
	if reflect.DeepEqual(expected, actual) {
		return
	}
	t.Helper()
	t.Errorf("\n"+
		"Expected: [%v]\n"+
		"Actual:   [%v]",
		expected,
		actual,
	)
}
func measureHistograms(metrics *sampleMetrics) {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	for x := int64(1); x < 1000; x = x * 2 {
		wg.Add(1)
		go func(measurement int64) {
			metrics.histogram1.Measure(measurement)
			wg.Done()
		}(x)
	}
	for x := int64(1); x < 1000; x = x * 3 {
		wg.Add(1)
		go func(measurement int64) {
			metrics.histogram2.Measure(measurement)
			wg.Done()
		}(x)
	}
}

func TestMetricsRendering(t *testing.T) {
	histogramBucketKeys := []int64{0, 1, 20, 30, 50, 100, 300, 500}
	metrics := newSampleMetrics(histogramBucketKeys)

	metrics.counter1.IncrementN(1)
	metrics.counter2.IncrementN(2)
	metrics.gauge1.IncrementN(3)
	metrics.gauge2.Measure(4)

	measureHistograms(metrics)

	exporter := NewExporter()
	exporter.Add(
		metrics.counter1,
		metrics.counter2,
		metrics.gauge1,
		metrics.gauge2,
		metrics.histogram1,
		metrics.histogram2,
	)
	recorder := httptest.NewRecorder()

	exporter.ServeHTTP(recorder, nil)

	actualBody := recorder.Body.String()
	assertEqual(t, strings.TrimSpace(expectedExporterBody), strings.TrimSpace(actualBody))
}

const expectedExporterBody = `
# HELP my_counter counter description
# TYPE my_counter counter
my_counter 1

# HELP my_counter_with_labels counter description
# TYPE my_counter_with_labels counter
my_counter_with_labels{ counter_label_key="counter_label_value" } 2

# HELP my_gauge gauge description
# TYPE my_gauge gauge
my_gauge 3

# HELP my_gauge_with_labels gauge description
# TYPE my_gauge_with_labels gauge
my_gauge_with_labels{ gauge_label_key="gauge_label_value" } 4

# HELP my_histogram_with_buckets histogram description
# TYPE my_histogram_with_buckets histogram
my_histogram_with_buckets_bucket{ le="0", histogram_key1="histogram_value1" } 0
my_histogram_with_buckets_bucket{ le="1", histogram_key1="histogram_value1" } 1
my_histogram_with_buckets_bucket{ le="20", histogram_key1="histogram_value1" } 5
my_histogram_with_buckets_bucket{ le="30", histogram_key1="histogram_value1" } 5
my_histogram_with_buckets_bucket{ le="50", histogram_key1="histogram_value1" } 6
my_histogram_with_buckets_bucket{ le="100", histogram_key1="histogram_value1" } 7
my_histogram_with_buckets_bucket{ le="300", histogram_key1="histogram_value1" } 9
my_histogram_with_buckets_bucket{ le="500", histogram_key1="histogram_value1" } 9
my_histogram_with_buckets_bucket{ le="+Inf", histogram_key1="histogram_value1" } 10
my_histogram_with_buckets_count{ histogram_key1="histogram_value1" } 10
my_histogram_with_buckets_sum{ histogram_key1="histogram_value1" } 1023

# HELP my_histogram_with_buckets_and_labels histogram description
# TYPE my_histogram_with_buckets_and_labels histogram
my_histogram_with_buckets_and_labels_bucket{ le="0", histogram_key1="histogram_value1" } 0
my_histogram_with_buckets_and_labels_bucket{ le="1", histogram_key1="histogram_value1" } 1
my_histogram_with_buckets_and_labels_bucket{ le="20", histogram_key1="histogram_value1" } 3
my_histogram_with_buckets_and_labels_bucket{ le="30", histogram_key1="histogram_value1" } 4
my_histogram_with_buckets_and_labels_bucket{ le="50", histogram_key1="histogram_value1" } 4
my_histogram_with_buckets_and_labels_bucket{ le="100", histogram_key1="histogram_value1" } 5
my_histogram_with_buckets_and_labels_bucket{ le="300", histogram_key1="histogram_value1" } 6
my_histogram_with_buckets_and_labels_bucket{ le="500", histogram_key1="histogram_value1" } 6
my_histogram_with_buckets_and_labels_bucket{ le="+Inf", histogram_key1="histogram_value1" } 7
my_histogram_with_buckets_and_labels_count{ histogram_key1="histogram_value1" } 7
my_histogram_with_buckets_and_labels_sum{ histogram_key1="histogram_value1" } 1093
`

type sampleMetrics struct {
	counter1   Counter
	counter2   Counter
	gauge1     Gauge
	gauge2     Gauge
	histogram1 Histogram
	histogram2 Histogram
}
