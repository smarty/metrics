package metrics

import "testing"

// Labels are stored in a map, so without sorting each metric could render them in a
// different order. Building many metrics makes an unsorted rendering all but certain
// to show up at least once.
func TestLabelsRenderSortedByKey(t *testing.T) {
	for range 100 {
		counter := NewCounter("my_counter",
			Options.Label("stream", "queue-a"),
			Options.Label("namespace", "billing"),
			Options.Label("source", "consumer"),
		)

		assertEqual(t, `{ namespace="billing", source="consumer", stream="queue-a" }`, counter.Labels())
	}
}
