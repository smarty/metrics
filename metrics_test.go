package metrics

import (
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestMetricsRendering(t *testing.T) {
	counter := NewCounter("my_counter",
		Options.Description("counter description"),
	)
	counterWithLabels := NewCounter("my_counter_with_labels",
		Options.Description("counter description"),
		Options.Label("counter_label_key", "counter_label_value"),
	)
	gauge := NewGauge("my_gauge",
		Options.Description("gauge description"),
	)
	gaugeWithLabels := NewGauge("my_gauge_with_labels",
		Options.Description("gauge description"),
		Options.Label("gauge_label_key", "gauge_label_value"),
	)

	counter.Increment()
	counterWithLabels.IncrementN(2)
	gauge.Increment()
	gauge.IncrementN(2)
	gaugeWithLabels.Measure(4)

	exporter := NewExporter()
	exporter.Add(counter, counterWithLabels, gauge, gaugeWithLabels)
	recorder := httptest.NewRecorder()

	exporter.ServeHTTP(recorder, nil)

	actualBody := strings.TrimSpace(recorder.Body.String())

	assertEqual(t, expectedBody, actualBody)
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

var expectedBody = strings.TrimSpace(`
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
`)
