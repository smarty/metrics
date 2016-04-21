package metrics

import (
	"errors"
	"os"
	"time"
)

type (
	CounterMetric   int
	GaugeMetric     int
	HistogramMetric int
)

const MetricConflict = -1

var standard = New()

// AddCounter registers a named counter metric along with the desired reporting frequency.
// The function is meant to be called *only* at application startup and is not thread safe.
// A negative return value indicates that the registration was unsuccessful.
func AddCounter(name string, reportingFrequency time.Duration) CounterMetric {
	return standard.AddCounter(name, reportingFrequency)
}

// AddGauge registers a named gauge metric along with the desired reporting frequency.
// The function is meant to be called *only* at application startup and is not thread safe.
// A negative return value indicates that the registration was unsuccessful.
func AddGauge(name string, reportingFrequency time.Duration) GaugeMetric {
	return standard.AddGauge(name, reportingFrequency)
}

// AddHistogram registers a named histogram metric along with the desired reporting frequency.
// The function is meant to be called *only* at application startup and is not thread safe.
// A negative return value indicates that the registration was unsuccessful.
func AddHistogram(
	name string, reportingFrequency time.Duration,
	min, max int64, resolution int, quantiles ...float64) HistogramMetric {

	return standard.AddHistogram(name, reportingFrequency, min, max, resolution, quantiles...)
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

// Count (automically) increments the metric indicated by one.
// A return value of false indicates the count could not occur.
func Count(id CounterMetric) bool {
	return standard.Count(id)
}

// Count (automically) increments the metric indicated by the number provided.
// A return value of false indicates the count could not occur.
func CountN(id CounterMetric, increment int64) bool {
	return standard.CountN(id, increment)
}

// RawCount is similiar to Measure, except that the backend is still a counter
// but the application is now responsible for tracking the count rather than the metric itself tracking it.
// A return value of false indicates the count could not occur.
func RawCount(id CounterMetric, value int64) bool {
	return standard.RawCount(id, value)
}

// Measure (automically) sets the metric of the metric indicated to the value specified
// A return value of false indicates the count could not occur.
func Measure(id GaugeMetric, value int64) bool {
	return standard.Measure(id, value)
}

///////////////////////////////////////////////////////////////////////////////

// StartLibrato configures a new metrics instance, specifies a number of Librato writers, and starts measuring.
func StartLibrato(email, key string, queueCapacity, writers int) error {
	if len(email) == 0 || len(key) == 0 || queueCapacity <= 0 || writers <= 0 {
		return libratoConfigurationError
	}

	queue := make(chan []MetricMeasurement, queueCapacity)
	hostname, _ := os.Hostname()
	go newLibrato(email, key, hostname, int32(writers)).Listen(queue)
	go sendRegularMeasurements(queue)
	return nil
}

func sendRegularMeasurements(queue chan []MetricMeasurement) {
	for StartMeasuring(); ; time.Sleep(time.Second) {
		queue <- standard.TakeMeasurements(time.Now())
	}
}

var libratoConfigurationError = errors.New("You must supply non-empty email address, non-empty key, and positive queueCapacity and positive writers.")
