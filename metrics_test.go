package metrics

import (
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type MetricsTrackerFixture struct {
	*gunit.Fixture

	tracker *MetricsTracker
	now     time.Time
}

func (this *MetricsTrackerFixture) Setup() {
	this.tracker = New2()
	this.now = time.Now()
}

func (this *MetricsTrackerFixture) TestMeasuringCounters() {
	counter1 := this.tracker.AddCounter("counter1", time.Nanosecond)
	counter2 := this.tracker.AddCounter("counter2", time.Nanosecond)

	this.tracker.StartMeasuring()

	this.So(this.tracker.Count(counter1), should.BeTrue)
	this.So(this.tracker.CountN(counter1, 9), should.BeTrue)
	this.So(this.tracker.CountRaw(counter2, 42), should.BeTrue)
	this.So(this.tracker.TakeMeasurements(this.now), should.Resemble, []MetricMeasurement{
		{
			Captured:   this.now,
			ID:         0,
			Name:       "counter1",
			MetricType: counterMetricType,
			Value:      10,
		},
		{
			Captured:   this.now,
			ID:         1,
			Name:       "counter2",
			MetricType: counterMetricType,
			Value:      42,
		},
	})
}

func (this *MetricsTrackerFixture) TestCounterIDMustBeWithinRange() {
	counter1 := this.tracker.AddCounter("counter1", time.Nanosecond)

	this.tracker.StartMeasuring()

	this.So(this.tracker.Count(counter1+1), should.BeFalse)
	this.So(this.tracker.TakeMeasurements(this.now), should.Resemble, []MetricMeasurement{
		{
			Captured:   this.now,
			ID:         0,
			Name:       "counter1",
			MetricType: counterMetricType,
			Value:      0,
		},
	})
}

func (this *MetricsTrackerFixture) TestMeasuringGauges() {
	gauge1 := this.tracker.AddGauge("gauge1", time.Nanosecond)
	gauge2 := this.tracker.AddGauge("gauge2", time.Nanosecond)

	this.tracker.StartMeasuring()

	this.So(this.tracker.Measure(gauge1, 1111), should.BeTrue) // 1111 will be overwritten on the next measurement
	this.So(this.tracker.Measure(gauge1, 2222), should.BeTrue)
	this.So(this.tracker.Measure(gauge2, 3333), should.BeTrue)
	this.So(this.tracker.TakeMeasurements(this.now), should.Resemble, []MetricMeasurement{
		{
			Captured:   this.now,
			ID:         0,
			Name:       "gauge1",
			MetricType: gaugeMetricType,
			Value:      2222,
		},
		{
			Captured:   this.now,
			ID:         1,
			Name:       "gauge2",
			MetricType: gaugeMetricType,
			Value:      3333,
		},
	})
}

func (this *MetricsTrackerFixture) TestGaugeIDMustBeWithinRange() {
	gauge1 := this.tracker.AddGauge("gauge1", time.Nanosecond)

	this.tracker.StartMeasuring()

	this.So(this.tracker.Measure(gauge1+1, 42), should.BeFalse)
	this.So(this.tracker.TakeMeasurements(this.now), should.Resemble, []MetricMeasurement{
		{
			Captured:   this.now,
			ID:         0,
			Name:       "gauge1",
			MetricType: gaugeMetricType,
			Value:      0,
		},
	})
}

func (this *MetricsTrackerFixture) TestMeasuringHistograms() {
	histogram1 := this.tracker.AddHistogram("histogram1", time.Nanosecond, 0, 100, 5, 50, 99)
	this.tracker.StartMeasuring()

	for x := 1; x <= 100; x++ {
		this.So(this.tracker.Record(histogram1, int64(x)), should.BeTrue)
	}

	this.So(this.tracker.TakeMeasurements(this.now), should.Resemble, []MetricMeasurement{
		{
			Captured:   this.now,
			ID:         0,
			Value:      1,
			Name:       "histogram1_min",
			MetricType: gaugeMetricType,
		},
		{
			Captured:   this.now,
			ID:         1,
			Value:      100,
			Name:       "histogram1_max",
			MetricType: gaugeMetricType,
		},
		{
			Captured:   this.now,
			ID:         2,
			Value:      50,
			Name:       "histogram1_mean",
			MetricType: gaugeMetricType,
		},
		{
			Captured:   this.now,
			ID:         3,
			Value:      28,
			Name:       "histogram1_stddev",
			MetricType: gaugeMetricType,
		},
		{
			Captured:   this.now,
			ID:         4,
			Value:      100,
			Name:       "histogram1_total",
			MetricType: gaugeMetricType,
		},
		{
			Captured:   this.now,
			ID:         5,
			Value:      50,
			Name:       "histogram1_50.000",
			MetricType: gaugeMetricType,
		},
		{
			Captured:   this.now,
			ID:         6,
			Value:      99,
			Name:       "histogram1_99.000",
			MetricType: gaugeMetricType,
		},
	})
}
