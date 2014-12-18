package metrics

import (
	"sync/atomic"
	"time"
)

///////////////////////////////////////////////////////////////////////////////

var (
	metrics []int64
	meta    []metricInfo
	started int32
	queue   chan []Measurement
)

///////////////////////////////////////////////////////////////////////////////

// Add registers a named metric along with the desired reporting frequency.
// The function is meant to be called *only* at application startup.
func Add(name string, reportingFrequency time.Duration) int {
	if atomic.LoadInt32(&started) > 0 {
		return -1
	}

	// TODO: if name already taken; if reporting frequency is zero or negative
	metrics = append(metrics, int64(0))
	info := metricInfo{Name: name, ReportingFrequency: reportingFrequency}
	meta = append(meta, info)
	return len(metrics) - 1
}

func RegisterChannelDestination(destination chan []Measurement) {
	queue = destination
}

func StartMeasuring() bool {
	if atomic.AddInt32(&started, 1) > 1 {
		return false
	}

	durations := map[time.Duration][]int{}
	for i, item := range meta {
		indices := durations[item.ReportingFrequency]
		indices = append(indices, i)
		durations[item.ReportingFrequency] = indices
	}

	for d, indices := range durations {
		duration := d // save the value for the closure below...
		time.AfterFunc(duration, func() { report(duration, indices) })
	}

	return true
}

func report(duration time.Duration, indices []int) {
	now := time.Now()
	snapshot := make([]Measurement, len(indices), len(indices))

	for i := 0; i < len(indices); i++ {
		index := indices[i]
		snapshot[i] = Measurement{
			Index:    index,
			Captured: now,
			Value:    atomic.LoadInt64(&metrics[index]),
		}
	}

	queue <- snapshot

	if started > 0 {
		time.AfterFunc(duration, func() { report(duration, indices) })
	}
}

///////////////////////////////////////////////////////////////////////////////

func StopMeasuring() {
	started = 0
}

///////////////////////////////////////////////////////////////////////////////

func Count(index int) bool {
	if index < 0 || len(metrics) <= index {
		return false
	}

	atomic.AddInt64(&metrics[index], 1)
	return true
}

func Measure(index int, measurement int64) bool {
	if index < 0 || len(metrics) <= index {
		return false
	}

	atomic.StoreInt64(&metrics[index], measurement)
	return true
}

///////////////////////////////////////////////////////////////////////////////

type metricInfo struct {
	Name               string
	ReportingFrequency time.Duration
}

type Measurement struct {
	Index    int
	Captured time.Time
	Value    int64
}

///////////////////////////////////////////////////////////////////////////////
