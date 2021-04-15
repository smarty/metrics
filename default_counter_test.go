package metrics

import (
	"reflect"
	"testing"
)

func TestCounterValues(t *testing.T) {
	counter := NewCounter("")

	counter.Increment()
	assertEqual(t, uint64(1), counter.Value())

	counter.IncrementN(2)
	assertEqual(t, uint64(3), counter.Value())
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
