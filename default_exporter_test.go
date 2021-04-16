package metrics

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestExporter(t *testing.T) {
	counter := NewCounter("my_counter_with_labels",
		Options.Description("counter description"),
		Options.Label("counter_label_key", "counter_label_value"))
	gauge := NewGauge("my_gauge_with_labels",
		Options.Description("gauge description"),
		Options.Label("gauge_label_key", "gauge_label_value"))
	histogram := NewHistogram("my_histogram_with_buckets",
		Options.Description("histogram description"),
		Options.Label("histogram_key1", "histogram_value1"),
		Options.Bucket(0),
		Options.Bucket(1))
	histogram2 := NewHistogram("my_histogram_with_buckets_no_labels",
		Options.Description("histogram description"),
		Options.Bucket(0),
		Options.Bucket(1))
	exporter := NewExporter()
	exporter.Add(counter, gauge, histogram, histogram2)

	counter.Increment()
	gauge.IncrementN(2)
	histogram.Measure(1)
	histogram.Measure(2)
	histogram.Measure(3)
	histogram2.Measure(42)

	recorder := httptest.NewRecorder()
	exporter.ServeHTTP(recorder, nil)

	actualBody := recorder.Body.String()
	assertEqual(t, strings.TrimSpace(expectedBody), strings.TrimSpace(actualBody))
}

const expectedBody = `
# HELP my_counter_with_labels counter description
# TYPE my_counter_with_labels counter
my_counter_with_labels{ counter_label_key="counter_label_value" } 1

# HELP my_gauge_with_labels gauge description
# TYPE my_gauge_with_labels gauge
my_gauge_with_labels{ gauge_label_key="gauge_label_value" } 2

# HELP my_histogram_with_buckets histogram description
# TYPE my_histogram_with_buckets histogram
my_histogram_with_buckets_bucket{ le="0", histogram_key1="histogram_value1" } 0
my_histogram_with_buckets_bucket{ le="1", histogram_key1="histogram_value1" } 1
my_histogram_with_buckets_bucket{ le="+Inf", histogram_key1="histogram_value1" } 3
my_histogram_with_buckets_count{ histogram_key1="histogram_value1" } 3
my_histogram_with_buckets_sum{ histogram_key1="histogram_value1" } 6

# HELP my_histogram_with_buckets_no_labels histogram description
# TYPE my_histogram_with_buckets_no_labels histogram
my_histogram_with_buckets_no_labels_bucket{ le="0" } 0
my_histogram_with_buckets_no_labels_bucket{ le="1" } 0
my_histogram_with_buckets_no_labels_bucket{ le="+Inf" } 1
my_histogram_with_buckets_no_labels_count 1
my_histogram_with_buckets_no_labels_sum 42
`
