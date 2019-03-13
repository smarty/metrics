package metrics

import (
	"sort"
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/logging"
)

func TestMetricsTrackerFixture(t *testing.T) {
	gunit.Run(new(MetricsTrackerFixture), t)
}

type MetricsTrackerFixture struct {
	*gunit.Fixture

	tracker *MetricsTracker
	now     time.Time
}

func (this *MetricsTrackerFixture) Setup() {
	this.tracker = New()
	this.tracker.logger = logging.Capture()
	this.now = time.Now()
}

func (this *MetricsTrackerFixture) measure() []MetricMeasurement {
	measurements := this.tracker.TakeMeasurements(this.now)
	sort.Slice(measurements, func(i, j int) bool {
		return measurements[i].ID < measurements[j].ID
	})
	return measurements
}

func (this *MetricsTrackerFixture) TestMetricNamesMustBeUnique() {
	counter1 := this.tracker.AddCounter("1", time.Nanosecond)
	counter1a := this.tracker.AddCounter("1", time.Nanosecond)
	gauge1 := this.tracker.AddGauge("1", time.Nanosecond)
	histogram1 := this.tracker.AddHistogram("1", time.Nanosecond, 0, 10, 5, 50)
	this.So(counter1, should.Equal, 0)
	this.So(counter1a, should.Equal, MetricConflict)
	this.So(gauge1, should.Equal, MetricConflict)
	this.So(histogram1, should.Equal, MetricConflict)

	this.tracker.StartMeasuring()
	this.So(this.measure(), should.HaveLength, 1)
}
func (this *MetricsTrackerFixture) TestMetricsCannotBeAddedWhenTheTrackerHasBeenStarted() {
	this.tracker.StartMeasuring()
	this.So(this.tracker.AddCounter("1", time.Nanosecond), should.Equal, MetricConflict)
	this.So(this.tracker.AddGauge("2", time.Nanosecond), should.Equal, MetricConflict)
	this.So(this.tracker.AddHistogram("3", time.Nanosecond, 0, 10, 5, 50), should.Equal, MetricConflict)
}
func (this *MetricsTrackerFixture) TestMetricsCanBeAddedAfterStoppingTheTracker() {
	this.tracker.StartMeasuring()
	this.tracker.StopMeasuring()
	this.So(this.tracker.AddCounter("1", time.Nanosecond), should.NotEqual, MetricConflict)
	this.So(this.tracker.AddGauge("2", time.Nanosecond), should.NotEqual, MetricConflict)
	this.So(this.tracker.AddHistogram("3", time.Nanosecond, 0, 10, 5, 50), should.NotEqual, MetricConflict)
}
func (this *MetricsTrackerFixture) TestReportingFrequencyMustBePositive() {
	var duration time.Duration
	this.So(this.tracker.AddCounter("1", duration), should.Equal, MetricConflict)
	this.So(this.tracker.AddGauge("2", duration), should.Equal, MetricConflict)
	this.So(this.tracker.AddHistogram("3", duration, 0, 10, 5, 50), should.Equal, MetricConflict)
}
func (this *MetricsTrackerFixture) TestMeasurementsOnlyAvailableWhenTrackerIsStarted() {
	this.tracker.AddCounter("counter1", time.Nanosecond)
	this.So(this.measure(), should.BeEmpty)

	this.tracker.StartMeasuring()
	this.So(this.measure(), should.NotBeEmpty)

	this.tracker.StopMeasuring()
	this.So(this.tracker.TakeMeasurements(this.now.Add(time.Second)), should.BeEmpty)

	this.tracker.StartMeasuring()
	this.So(this.tracker.TakeMeasurements(this.now.Add(time.Second)), should.NotBeEmpty)
}

func (this *MetricsTrackerFixture) TestMeasuringCounters() {
	counter1 := this.tracker.AddCounter("counter1", time.Nanosecond)
	counter2 := this.tracker.AddCounter("counter2", time.Nanosecond)

	this.tracker.TagCounter(counter1, "color", "plaid", "size", "XL")
	this.tracker.TagCounter(counter2, "color", "fuscia")

	this.tracker.StartMeasuring()

	this.So(this.tracker.Count(counter1), should.BeTrue)
	this.So(this.tracker.CountN(counter1, 9), should.BeTrue)
	this.So(this.tracker.RawCount(counter2, 42), should.BeTrue)
	this.So(this.measure(), should.Resemble, []MetricMeasurement{
		{
			Captured:   this.now,
			ID:         0,
			Name:       "counter1",
			MetricType: counterMetricType,
			Value:      10,
			Tags:       map[string]string{"color": "plaid", "size": "XL"},
		},
		{
			Captured:   this.now,
			ID:         1,
			Name:       "counter2",
			MetricType: counterMetricType,
			Value:      42,
			Tags:       map[string]string{"color": "fuscia"},
		},
	})
}
func (this *MetricsTrackerFixture) TestCounterIDMustBeValid() {
	counter1 := this.tracker.AddCounter("counter1", time.Nanosecond)

	this.tracker.StartMeasuring()

	this.So(this.tracker.Count(counter1+1), should.BeFalse)
	this.So(this.measure(), should.Resemble, []MetricMeasurement{
		{
			Captured:   this.now,
			ID:         0,
			Name:       "counter1",
			MetricType: counterMetricType,
			Value:      0,
		},
	})
}
func (this *MetricsTrackerFixture) TestCounterNameMustHaveSubstance() {
	counter1 := this.tracker.AddCounter("\t\t\n\n     \t", time.Nanosecond)
	this.So(counter1, should.Equal, MetricConflict)
	this.So(this.measure(), should.BeEmpty)
}

func (this *MetricsTrackerFixture) TestMeasuringGauges() {
	gauge1 := this.tracker.AddGauge("gauge1", time.Nanosecond)
	gauge2 := this.tracker.AddGauge("gauge2", time.Nanosecond)
	this.tracker.TagGauge(gauge1, "color", "blue")
	this.tracker.TagGauge(gauge2, "color", "red")

	this.tracker.StartMeasuring()

	this.So(this.tracker.Measure(gauge1, 1111), should.BeTrue) // 1111 will be overwritten on the next measurement
	this.So(this.tracker.Measure(gauge1, 2222), should.BeTrue)
	this.So(this.tracker.Measure(gauge2, 3333), should.BeTrue)
	this.So(this.measure(), should.Resemble, []MetricMeasurement{
		{
			Captured:   this.now,
			ID:         0,
			Name:       "gauge1",
			MetricType: gaugeMetricType,
			Value:      2222,
			Tags:       map[string]string{"color": "blue"},
		},
		{
			Captured:   this.now,
			ID:         1,
			Name:       "gauge2",
			MetricType: gaugeMetricType,
			Value:      3333,
			Tags:       map[string]string{"color": "red"},
		},
	})
}
func (this *MetricsTrackerFixture) TestGaugeIDMustBeValid() {
	gauge1 := this.tracker.AddGauge("gauge1", time.Nanosecond)

	this.tracker.StartMeasuring()

	this.So(this.tracker.Measure(gauge1+1, 42), should.BeFalse)
	this.So(this.measure(), should.Resemble, []MetricMeasurement{
		{
			Captured:   this.now,
			ID:         0,
			Name:       "gauge1",
			MetricType: gaugeMetricType,
			Value:      0,
		},
	})
}
func (this *MetricsTrackerFixture) TestGaugeNameMustHaveSubstance() {
	gauge1 := this.tracker.AddGauge("\t\t\n\n     \t", time.Nanosecond)
	this.So(gauge1, should.Equal, MetricConflict)
	this.So(this.measure(), should.BeEmpty)
}

func (this *MetricsTrackerFixture) TestMeasuringHistograms() {
	histogram1 := this.tracker.AddHistogram("histogram1", time.Nanosecond, 0, 100, 5, 50, 99)
	this.tracker.TagHistogram(histogram1, "color", "plaid")
	this.tracker.StartMeasuring()

	for x := 1; x <= 100; x++ {
		this.So(this.tracker.Record(histogram1, int64(x)), should.BeTrue)
	}

	this.So(this.tracker.Record(histogram1, 1000000), should.BeFalse) // out of range (value too large)

	this.So(this.measure(), should.Resemble, []MetricMeasurement{
		{
			Captured:   this.now,
			ID:         1,
			Value:      1,
			Name:       "histogram1_min",
			MetricType: gaugeMetricType,
			Tags:       map[string]string{"color": "plaid"},
		},
		{
			Captured:   this.now,
			ID:         2,
			Value:      100,
			Name:       "histogram1_max",
			MetricType: gaugeMetricType,
			Tags:       map[string]string{"color": "plaid"},
		},
		{
			Captured:   this.now,
			ID:         3,
			Value:      50,
			Name:       "histogram1_mean",
			MetricType: gaugeMetricType,
			Tags:       map[string]string{"color": "plaid"},
		},
		{
			Captured:   this.now,
			ID:         4,
			Value:      28,
			Name:       "histogram1_stddev",
			MetricType: gaugeMetricType,
			Tags:       map[string]string{"color": "plaid"},
		},
		{
			Captured:   this.now,
			ID:         5,
			Value:      100,
			Name:       "histogram1_total",
			MetricType: gaugeMetricType,
			Tags:       map[string]string{"color": "plaid"},
		},
		{
			Captured:   this.now,
			ID:         6,
			Value:      50,
			Name:       "histogram1_50.000",
			MetricType: gaugeMetricType,
			Tags:       map[string]string{"color": "plaid"},
		},
		{
			Captured:   this.now,
			ID:         7,
			Value:      99,
			Name:       "histogram1_99.000",
			MetricType: gaugeMetricType,
			Tags:       map[string]string{"color": "plaid"},
		},
	})
}
func (this *MetricsTrackerFixture) TestHistogramIDMustBeValid() {
	histogram := this.tracker.AddHistogram("histogram1", time.Nanosecond, 0, 10, 3, 90)

	this.tracker.StartMeasuring()

	this.So(this.tracker.Record(histogram+1, 42), should.BeFalse)
	this.So(this.measure(), should.Resemble, []MetricMeasurement{
		{
			Captured:   this.now,
			ID:         1,
			Value:      0,
			Name:       "histogram1_min",
			MetricType: gaugeMetricType,
		},
		{
			Captured:   this.now,
			ID:         2,
			Value:      0,
			Name:       "histogram1_max",
			MetricType: gaugeMetricType,
		},
		{
			Captured:   this.now,
			ID:         3,
			Value:      0,
			Name:       "histogram1_mean",
			MetricType: gaugeMetricType,
		},
		{
			Captured:   this.now,
			ID:         4,
			Value:      0,
			Name:       "histogram1_stddev",
			MetricType: gaugeMetricType,
		},
		{
			Captured:   this.now,
			ID:         5,
			Value:      0,
			Name:       "histogram1_total",
			MetricType: gaugeMetricType,
		},
		{
			Captured:   this.now,
			ID:         6,
			Value:      0,
			Name:       "histogram1_90.000",
			MetricType: gaugeMetricType,
		},
	})
}
func (this *MetricsTrackerFixture) TestHistogramNameMustHaveSubstance() {
	histogram1 := this.tracker.AddHistogram("\t\t\n\n     \t", time.Nanosecond, 0, 10, 5, 50)
	this.So(histogram1, should.Equal, MetricConflict)
	this.So(this.measure(), should.BeEmpty)
}
func (this *MetricsTrackerFixture) TestHistogramMinAndMaxMustBeAscending() {
	this.So(this.tracker.AddHistogram("h", time.Nanosecond, 10, 9, 3, 50), should.Equal, MetricConflict)
}
func (this *MetricsTrackerFixture) TestHistogramResolutionMustBeWithinTolerance() {
	// invalid range:                   resolution column: | ----->##<-----|
	this.So(this.tracker.AddHistogram("h", time.Nanosecond, 1, 10, -1, 50), should.Equal, MetricConflict)
	this.So(this.tracker.AddHistogram("h0", time.Nanosecond, 1, 10, 0, 50), should.Equal, MetricConflict)
	// valid range:
	this.So(this.tracker.AddHistogram("h1", time.Nanosecond, 1, 10, 1, 50), should.NotEqual, MetricConflict)
	this.So(this.tracker.AddHistogram("h2", time.Nanosecond, 1, 10, 2, 50), should.NotEqual, MetricConflict)
	this.So(this.tracker.AddHistogram("h3", time.Nanosecond, 1, 10, 3, 50), should.NotEqual, MetricConflict)
	this.So(this.tracker.AddHistogram("h4", time.Nanosecond, 1, 10, 4, 50), should.NotEqual, MetricConflict)
	this.So(this.tracker.AddHistogram("h5", time.Nanosecond, 1, 10, 5, 50), should.NotEqual, MetricConflict)
	// invalid range:
	this.So(this.tracker.AddHistogram("h6", time.Nanosecond, 1, 10, 6, 50), should.Equal, MetricConflict)
	this.So(this.tracker.AddHistogram("h7", time.Nanosecond, 1, 10, 7, 50), should.Equal, MetricConflict)
}

func (this *MetricsTrackerFixture) TestOddNumberOfTagsIgnored() {
	counter := this.tracker.AddCounter("counter", time.Nanosecond)
	this.tracker.TagCounter(counter, "1", "2", "extra")
	this.tracker.StartMeasuring()
	this.So(this.measure(), should.Resemble, []MetricMeasurement{
		{
			Captured:   this.now,
			ID:         0,
			Name:       "counter",
			MetricType: counterMetricType,
			Value:      0,
			Tags:       nil,
		},
	})
	this.So(this.tracker.logger.Log.String(), should.ContainSubstring, "[WARN] tags must be submitted as an even number of key/value pairs")
}

func (this *MetricsTrackerFixture) TestTagAll() {
	this.tracker.AddCounter("counter", time.Nanosecond)
	this.tracker.AddHistogram("histogram", time.Nanosecond, 0, 10, 3, 90)
	this.tracker.AddGauge("gauge", time.Nanosecond)
	tags := map[string]string{"testing": "true", "valid": "false"}
	this.tracker.TagAll("testing", tags["testing"], "valid", tags["valid"])
	this.tracker.StartMeasuring()
	for _, measure := range this.measure() {
		actual := map[string]map[string]string{measure.Name: measure.Tags}
		expected := map[string]map[string]string{measure.Name: tags}
		this.So(actual, should.Resemble, expected)
	}
}
