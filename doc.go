package metrics

import (
	"errors"
	"os"
	"strings"
	"time"
)

type (
	CounterMetric   int
	GaugeMetric     int
	HistogramMetric int
)

const MetricConflict = -1

var standard = New()

func InitializeEnvironmentTags(tags string) {
	InitializeTags(tags, os.Getenv)
}

func InitializeTags(tags string, reader func(string) string) {
	var pairs []string

	for _, tag := range strings.Split(tags, "|") {
		value := reader(tag)
		pairs = append(pairs, tag, value)
	}

	TagAll(pairs...)
}

// AddCounter registers a named counter metric along with the desired reporting frequency.
// The function is meant to be called *only* at application startup and is not thread safe.
// A negative return value indicates that the registration was unsuccessful.
func AddCounter(name string, reportingFrequency time.Duration) CounterMetric {
	return standard.AddCounter(name, reportingFrequency)
}

// TagAll sets key/value pairs for all metrics.
func TagAll(tagPairs ...string) {
	standard.TagAll(tagPairs...)
}

// TagCounter associates key/value string pairs with the provided metric id.
func TagCounter(id CounterMetric, tagPairs ...string) {
	standard.TagCounter(id, tagPairs...)
}

// TagGauge associates key/value string pairs with the provided metric id.
func TagGauge(id GaugeMetric, tagPairs ...string) {
	standard.TagGauge(id, tagPairs...)
}

// TagHistogram associates key/value string pairs with the provided metric id.
func TagHistogram(id HistogramMetric, tagPairs ...string) {
	standard.TagHistogram(id, tagPairs...)
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

// Count increments the metric indicated by id.
// A return value of false indicates the count could not occur.
func Count(id CounterMetric) bool {
	return standard.Count(id)
}

// CountIf increments the metric indicated by id if the supplied condition is true.
// A return value of false indicates that the count could not occur or
// the supplied condition was false.
func CountIf(id CounterMetric, condition bool) bool {
	if condition {
		return Count(id)
	}
	return condition
}

// Count increments the metric indicated by the number provided.
// A return value of false indicates the count could not occur.
func CountN(id CounterMetric, increment int64) bool {
	return standard.CountN(id, increment)
}

// RawCount is similar to Measure, except that the backend is still a counter
// but the application is now responsible for tracking the count rather than the metric itself tracking it.
// A return value of false indicates the count could not occur.
func RawCount(id CounterMetric, value int64) bool {
	return standard.RawCount(id, value)
}

// Measure sets the metric indicated to the value specified.
// A return value of false indicates the count could not occur.
func Measure(id GaugeMetric, value int64) bool {
	return standard.Measure(id, value)
}

// Record records the value with the histogram indicated.
// A return value of false indicates the count could not occur.
func Record(id HistogramMetric, value int64) bool {
	return standard.Record(id, value)
}

///////////////////////////////////////////////////////////////////////////////

// StartAppOpticsConfig uses the provided configLoader to configure a new metrics instance,
// specify a number of AppOptics writers, and start measuring.
func StartAppOpticsConfig(configLoader AppOpticsConfigLoader, queueCapacity, writers int) error {
	config := configLoader()
	if len(config.Key) == 0 || queueCapacity <= 0 || writers <= 0 {
		return appOpticsConfigurationError
	}

	queue := make(chan []MetricMeasurement, queueCapacity)
	hostname, _ := os.Hostname()
	go newAppOptics(configLoader, hostname, int32(writers)).Listen(queue)
	go sendRegularMeasurements(queue)
	return nil
}

// StartAppOptics configures a new metrics instance, specifies a number of AppOptics writers, and starts measuring.
func StartAppOptics(key string, queueCapacity, writers int) error {
	configLoader := func() AppOpticsConfig {
		return AppOpticsConfig{Key: key}
	}

	return StartAppOpticsConfig(configLoader, queueCapacity, writers)
}

func sendRegularMeasurements(queue chan []MetricMeasurement) {
	for StartMeasuring(); ; time.Sleep(time.Second) {
		queue <- standard.TakeMeasurements(time.Now())
	}
}

var appOpticsConfigurationError = errors.New("You must supply non-empty key, and positive queueCapacity and positive writers.")

///////////////////////////////////////////////////////////////////////////////

// Microseconds converts the duration value to microseconds, a commonly used unit of measurement for latency metrics.
// This return value could easily be used as a gauge or histogram measurement.
func Microseconds(duration time.Duration) int64 { return int64(duration.Nanoseconds() / 1000) }

// Milliseconds converts the duration value to milliseconds, a possible unit of measurement for latency metrics.
// This return value could easily be used as a gauge or histogram measurement.
func Milliseconds(duration time.Duration) int64 { return Microseconds(duration) / 1000 }
