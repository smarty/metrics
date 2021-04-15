package metrics

import (
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func TestMetricsValues(t *testing.T) {
	metrics := NewTestMetrics()

	metrics.counter1.Increment()
	metrics.counter2.IncrementN(2)
	metrics.gauge1.Increment()
	metrics.gauge1.IncrementN(2)
	metrics.gauge2.Measure(4)

	assertEqual(t, int64(1), metrics.counter1.Value())
	assertEqual(t, int64(2), metrics.counter2.Value())
	assertEqual(t, int64(3), metrics.gauge1.Value())
	assertEqual(t, int64(4), metrics.gauge2.Value())

	measureHistograms(metrics)

	testBuckets1 := []Bucket{
		{key: 0, value: 0},
		{key: 1, value: 1},
		{key: 20, value: 5},
		{key: 30, value: 5},
		{key: 50, value: 6},
		{key: 100, value: 7},
		{key: 300, value: 9},
		{key: 500, value: 9},
	}

	for index, liveBucket := range metrics.histogram1.Buckets() {
		assertEqual(t, testBuckets1[index].Key(), liveBucket.Key())
		assertEqual(t, testBuckets1[index].Value(), liveBucket.Value())
	}
	assertEqual(t, uint64(10), metrics.histogram1.Count())
	assertEqual(t, uint64(1023), metrics.histogram1.Sum())

	testBuckets2 := []Bucket{
		{key: 0, value: 0},
		{key: 1, value: 1},
		{key: 20, value: 3},
		{key: 30, value: 4},
		{key: 50, value: 4},
		{key: 100, value: 5},
		{key: 300, value: 6},
		{key: 500, value: 6},
	}

	for index, liveBucket := range metrics.histogram2.Buckets() {
		assertEqual(t, testBuckets2[index].Key(), liveBucket.Key())
		assertEqual(t, testBuckets2[index].Value(), liveBucket.Value())
	}
	assertEqual(t, uint64(7), metrics.histogram2.Count())
	assertEqual(t, uint64(1093), metrics.histogram2.Sum())
}

func measureHistograms(metrics *TestMetrics) {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	for x := uint64(1); x < 1000; x = x * 2 {
		wg.Add(1)
		go func(measurement uint64) {
			metrics.histogram1.Measure(measurement)
			wg.Done()
		}(x)
	}
	for x := uint64(1); x < 1000; x = x * 3 {
		wg.Add(1)
		go func(measurement uint64) {
			metrics.histogram2.Measure(measurement)
			wg.Done()
		}(x)
	}
}

func TestMetricsRendering(t *testing.T) {
	metrics := NewTestMetrics()

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

	actualBody := strings.TrimSpace(recorder.Body.String())

	assertEqual(t, expectedExporterBody, actualBody)
}

var expectedExporterBody = strings.TrimSpace(`
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
my_histogram_with_buckets_bucket{ le="0" } 0
my_histogram_with_buckets_bucket{ le="1" } 1
my_histogram_with_buckets_bucket{ le="20" } 5
my_histogram_with_buckets_bucket{ le="30" } 5
my_histogram_with_buckets_bucket{ le="50" } 6
my_histogram_with_buckets_bucket{ le="100" } 7
my_histogram_with_buckets_bucket{ le="300" } 9
my_histogram_with_buckets_bucket{ le="500" } 9
my_histogram_with_buckets_bucket{ le="+Inf" } 10
my_histogram_with_buckets_count 10
my_histogram_with_buckets_sum 1023

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
`)

type TestMetrics struct {
	counter1   Counter
	counter2   Counter
	gauge1     Gauge
	gauge2     Gauge
	histogram1 Histogram
	histogram2 Histogram
}

func NewTestMetrics() *TestMetrics {
	counter1 := NewCounter("my_counter",
		Options.Description("counter description"),
	)
	counter2 := NewCounter("my_counter_with_labels",
		Options.Description("counter description"),
		Options.Label("counter_label_key", "counter_label_value"),
	)
	gauge1 := NewGauge("my_gauge",
		Options.Description("gauge description"),
	)
	gauge2 := NewGauge("my_gauge_with_labels",
		Options.Description("gauge description"),
		Options.Label("gauge_label_key", "gauge_label_value"),
	)
	histogram1 := NewHistogram("my_histogram_with_buckets",
		Options.Description("histogram description"),
		Options.Bucket(0),
		Options.Bucket(1),
		Options.Bucket(20),
		Options.Bucket(30),
		Options.Bucket(50),
		Options.Bucket(100),
		Options.Bucket(300),
		Options.Bucket(500),
	)
	histogram2 := NewHistogram("my_histogram_with_buckets_and_labels",
		Options.Description("histogram description"),
		Options.Bucket(0),
		Options.Bucket(1),
		Options.Bucket(20),
		Options.Bucket(30),
		Options.Bucket(50),
		Options.Bucket(100),
		Options.Bucket(300),
		Options.Bucket(500),
		Options.Label("histogram_key1", "histogram_value1"),
	)

	return &TestMetrics{
		counter1:   counter1,
		counter2:   counter2,
		gauge1:     gauge1,
		gauge2:     gauge2,
		histogram1: histogram1,
		histogram2: histogram2,
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
