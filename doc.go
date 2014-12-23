package metrics

import "time"

const MetricConflict = -1

// Add registers a named metric along with the desired reporting frequency.
// The function is meant to be called *only* at application startup.
// A negative return value indicates that the registration was unsuccessful.
func Add(name string, reportingFrequency time.Duration) int {
	return standard.Add(name, reportingFrequency)
}

// StartMeasuring signals to this library that all
// registrations have been performed.
func StartMeasuring() {
	standard.StartMeasuring()
}

// StopMeasuring turns measurement tracking off.
// It can be turned on again by calling StartMeasuring.
func StopMeasuring() {
	standard.StopMeasuring()
}

// Count (automically) increments the metric at index by one.
// A return value of false indicates the count could not occur.
func Count(index int) bool {
	return standard.Count(index)
}

// Measure (automically) sets the metric at the specified index to the specified measurement.
// A return value of false indicates the count could not occur.
func Measure(index int, measurement int64) bool {
	return standard.Measure(index, measurement)
}

///////////////////////////////////////////////////////////////////////////////
