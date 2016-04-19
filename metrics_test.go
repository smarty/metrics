package metrics

import (
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type MetricsTrackerFixture struct {
	*gunit.Fixture

	tracker *MetricsTracker
}

func (this *MetricsTrackerFixture) Setup() {
	this.tracker = New2()
}

func (this *MetricsTrackerFixture) TestCounters() {
	counter1 := this.tracker.AddCounter("counter1", time.Millisecond)
	counter2 := this.tracker.AddCounter("counter2", time.Millisecond)

	this.tracker.StartMeasuring()

	this.tracker.Count(counter1)
	this.tracker.Count(counter1)
	this.tracker.Count(counter2)

	now := time.Now()
	measurements := this.tracker.TakeMeasurements(now)
	this.So(measurements, should.Resemble, []MetricMeasurement{
		{
			Captured:   now,
			ID:         0,
			Name:       "counter1",
			MetricType: counterMetricType,
			Value:      2,
		},
		{
			Captured:   now,
			ID:         1,
			Name:       "counter2",
			MetricType: counterMetricType,
			Value:      1,
		},
	})
}
