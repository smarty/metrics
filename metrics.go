package metrics

import (
	"io/ioutil"
	"log"
	"sync/atomic"
	"time"
)

///////////////////////////////////////////////////////////////////////////////

// Add registers a named metric along with the desired reporting frequency.
// The function is meant to be called *only* at application startup.
func Add(name string, reportingFrequency time.Duration) int {
	return standard.Add(name, reportingFrequency)
}

// RegisterChannelDestination assigns the channel on which measurements
// will be published at their respective registered reporting frequencies.
func RegisterChannelDestination(destination chan []Measurement) {
	standard.RegisterChannelDestination(destination)
}

// StartMeasuring signals to this library that all
// registrations have been performed.
func StartMeasuring() bool {
	return standard.StartMeasuring()
}

// StopMeasuring turns measurement tracking off.
// It can be turned on again by calling StartMeasuring.
func StopMeasuring() {
	standard.StopMeasuring()
}

// Count (automically) increments the metric at index by one.
func Count(index int) bool {
	return standard.Count(index)
}

// Measure (automically) sets the metric at the specified index to the specified measurement.
func Measure(index int, measurement int64) bool {
	return standard.Measure(index, measurement)
}

// Measurement is the struct that is sent onto the outbound channel for
// publishing to whatever backend service that happens to be configured.
type Measurement struct {
	Index    int
	Captured time.Time
	Value    int64
}

///////////////////////////////////////////////////////////////////////////////

var standard = New()

type metric struct {
	metrics []int64
	meta    []metricInfo
	started int32
	queue   chan []Measurement
}

type metricInfo struct {
	Name               string
	ReportingFrequency time.Duration
}

func New() *metric {
	return &metric{}
}

func (this *metric) Add(name string, reportingFrequency time.Duration) int {
	if atomic.LoadInt32(&this.started) > 0 {
		return -1
	}

	// TODO: if name already taken; if reporting frequency is zero or negative
	this.metrics = append(this.metrics, int64(0))
	info := metricInfo{Name: name, ReportingFrequency: reportingFrequency}
	this.meta = append(this.meta, info)
	return len(this.metrics) - 1
}

func (this *metric) RegisterChannelDestination(destination chan []Measurement) {
	this.queue = destination
}

func (this *metric) StartMeasuring() bool {
	// if atomic.AddInt32(&this.started, 1) > 1 {
	// 	return false
	// }

	durations := map[time.Duration][]int{}
	for i, item := range this.meta {
		indices := durations[item.ReportingFrequency]
		indices = append(indices, i)
		durations[item.ReportingFrequency] = indices
	}

	for d, i := range durations {
		duration := d // save the values for the closure below...
		indices := i
		time.AfterFunc(duration, func() { this.report(duration, indices) })
	}

	this.started++
	return true
}

func (this *metric) report(duration time.Duration, indices []int) {
	now := time.Now()
	snapshot := make([]Measurement, len(indices), len(indices))

	for i := 0; i < len(indices); i++ {
		index := indices[i]
		snapshot[i] = Measurement{
			Index:    index,
			Captured: now,
			Value:    atomic.LoadInt64(&this.metrics[index]),
		}
	}

	this.queue <- snapshot

	// if this.started > 0 {
	time.AfterFunc(duration, func() { this.report(duration, indices) })
	// }
}

func (this *metric) StopMeasuring() {
	this.started = 0
}

func (this *metric) Count(index int) bool {
	if index < 0 || len(this.metrics) <= index || this.started < 1 {
		return false
	}

	atomic.AddInt64(&this.metrics[index], 1)
	return true
}

func (this *metric) Measure(index int, measurement int64) bool {
	if index < 0 || len(this.metrics) <= index {
		return false
	}

	atomic.StoreInt64(&this.metrics[index], measurement)
	return true
}

///////////////////////////////////////////////////////////////////////////////

func init() {
	log.SetOutput(ioutil.Discard)
}

///////////////////////////////////////////////////////////////////////////////
