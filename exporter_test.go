package metrics2

import (
	"log"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestExporter(t *testing.T) {
	recorder := httptest.NewRecorder()
	exporter := NewExporter()

	counter := NewCounter("my_counter",
		Options.Description("counter description"))
	counterWithLabels := NewCounter("my_counter_with_labels",
		Options.Description("counter description"),
		Options.Label("counter_label_key", "counter_label_value"))
	gauge := NewGauge("my_gauge",
		Options.Description("gauge description"))
	gaugeWithLabels := NewGauge("my_gauge_with_labels",
		Options.Description("gauge description"),
		Options.Label("gauge_label_key", "gauge_label_value"))
	exporter.Add(counter, counterWithLabels, gauge, gaugeWithLabels)

	counter.Increment()
	counterWithLabels.IncrementN(2)
	gauge.Increment()
	gauge.IncrementN(2)
	gaugeWithLabels.Measure(4)
	exporter.ServeHTTP(recorder, nil)

	actualBody := recorder.Body.String()
	if strings.TrimSpace(actualBody) != strings.TrimSpace(expectedBody) {
		log.Println(actualBody)
		t.Error("the actual response body didn't match the expected.")
	}
}

const expectedBody = `
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
`
