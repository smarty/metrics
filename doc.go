package metrics

import (
	"errors"
	"os"
	"time"
)

const MetricConflict = -1

// Add registers a named metric along with the desired reporting frequency.
// The function is meant to be called *only* at application startup.
// A negative return value indicates that the registration was unsuccessful.
func AddCounter(name string, reportingFrequency time.Duration) int {
	return standard.AddCounter(name, reportingFrequency)
}

// Add registers a named metric along with the desired reporting frequency.
// The function is meant to be called *only* at application startup.
// A negative return value indicates that the registration was unsuccessful.
func AddGauge(name string, reportingFrequency time.Duration) int {
	return standard.AddGauge(name, reportingFrequency)
}

// StartMeasuring signals to this library that all
// registrations have been performed.
func StartMeasuring() {
	standard.StartMeasuring()
}

func RegisterChannelDestination(queue chan []Measurement) {
	standard.RegisterChannelDestination(queue)
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

// StartLibrato configures a new metrics instance, specifies a number of Librato writers, and starts measuring.
func StartLibrato(email, key string, queueCapacity, writers int) error {
	if len(email) == 0 || len(key) == 0 || queueCapacity <= 0 || writers <= 0 {
		return libratoConfigurationError
	}

	queue := make(chan []Measurement, queueCapacity)
	RegisterChannelDestination(queue)

	hostname, _ := os.Hostname()
	librato := newLibrato(email, key, hostname, int32(writers))
	go librato.Listen(queue)

	StartMeasuring()

	return nil
}

var libratoConfigurationError = errors.New("You must supply non-empty email address, non-empty key, and positive queueCapacity and positive writers.")
