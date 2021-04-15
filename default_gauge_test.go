package metrics

import "testing"

func TestGaugeValues(t *testing.T) {
	gauge := NewGauge("")

	gauge.Increment()
	assertEqual(t, int64(1), gauge.Value())

	gauge.IncrementN(2)
	assertEqual(t, int64(3), gauge.Value())

	gauge.Measure(-42)
	assertEqual(t, int64(-42), gauge.Value())
}
