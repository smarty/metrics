package metrics

import (
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestMetricsValues(t *testing.T) {
	metrics := NewTestMetrics()

	metrics.counter1.Increment()
	metrics.counter2.IncrementN(2)
	metrics.gauge1.Increment()
	metrics.gauge1.IncrementN(2)
	metrics.gauge2.Measure(4)
	metrics.histogram1.Increment()
	metrics.histogram2.Observe(6)

	assertEqual(t, int64(1), metrics.counter1.Value())
	assertEqual(t, int64(2), metrics.counter2.Value())
	assertEqual(t, int64(3), metrics.gauge1.Value())
	assertEqual(t, int64(4), metrics.gauge2.Value())
	assertEqual(t, int64(1), metrics.histogram1.Value())
	assertEqual(t, int64(6), metrics.histogram2.Value())
}

func TestMetricsRendering(t *testing.T) {
	metrics := NewTestMetrics()

	metrics.counter1.IncrementN(1)
	metrics.counter2.IncrementN(2)
	metrics.gauge1.IncrementN(3)
	metrics.gauge2.Measure(4)
	metrics.histogram1.Observe(5)
	metrics.histogram2.Observe(6)

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

# HELP my_histogram histogram description
# TYPE my_histogram histogram
my_histogram 5

# HELP my_histogram_with_labels histogram description
# TYPE my_histogram_with_labels histogram
my_histogram_with_labels{ histogram_label_key="histogram_label_value" } 6
`)

type TestMetrics struct {
	counter1 Counter
	counter2 Counter
	gauge1   Gauge
	gauge2   Gauge
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
	histogram1 := NewHistogram("my_histogram",
		Options.Description("histogram description"),
	)
	histogram2 := NewHistogram("my_histogram_with_labels",
		Options.Description("histogram description"),
		Options.Label("histogram_label_key", "histogram_label_value"),
	)

	return &TestMetrics{
		counter1: counter1,
		counter2: counter2,
		gauge1:   gauge1,
		gauge2:   gauge2,
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
