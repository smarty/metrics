package metrics

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestHistogramMetricsFixture(t *testing.T) {
	gunit.Run(new(HistogramMetricsFixture), t)
}

type HistogramMetricsFixture struct {
	*gunit.Fixture

	histogram *FakeHistogram
}

func (this *HistogramMetricsFixture) Setup() {
	this.histogram = &FakeHistogram{}
}

func (this *HistogramMetricsFixture) TestHistogramMetricsSatisfyMetricInterface() {
	var metric *Metric
	this.So(new(HistogramMinMetric), should.Implement, metric)
	this.So(new(HistogramMaxMetric), should.Implement, metric)
	this.So(new(HistogramMeanMetric), should.Implement, metric)
	this.So(new(HistogramStandardDeviationMetric), should.Implement, metric)
	this.So(new(HistogramQuantileMetric), should.Implement, metric)
	this.So(new(HistogramTotalCountMetric), should.Implement, metric)
}

func (this *HistogramMetricsFixture) TestMinCalculation() {
	metric := &HistogramMinMetric{histogram: this.histogram, name: "hello"}
	this.So(metric.Measure(), should.Resemble, MetricMeasurement{
		MetricType: gaugeMetricType,
		Name:       "hello",
		Value:      12345,
	})
}
func (this *HistogramMetricsFixture) TestMaxCalculation() {
	metric := &HistogramMaxMetric{histogram: this.histogram, name: "hello"}
	this.So(metric.Measure(), should.Resemble, MetricMeasurement{
		MetricType: gaugeMetricType,
		Name:       "hello",
		Value:      54321,
	})
}
func (this *HistogramMetricsFixture) TestMeanCalculation() {
	metric := &HistogramMeanMetric{histogram: this.histogram, name: "hello"}
	this.So(metric.Measure(), should.Resemble, MetricMeasurement{
		MetricType: gaugeMetricType,
		Name:       "hello",
		Value:      123,
	})
}
func (this *HistogramMetricsFixture) TestStandardDeviationCalculation() {
	metric := &HistogramStandardDeviationMetric{histogram: this.histogram, name: "hello"}
	this.So(metric.Measure(), should.Resemble, MetricMeasurement{
		MetricType: gaugeMetricType,
		Name:       "hello",
		Value:      54,
	})
}
func (this *HistogramMetricsFixture) TestQuantileCalculation() {
	metric := &HistogramQuantileMetric{histogram: this.histogram, name: "hello", quantile: 99.9}
	this.So(metric.Measure(), should.Resemble, MetricMeasurement{
		MetricType: gaugeMetricType,
		Name:       "hello",
		Value:      99,
	})
}
func (this *HistogramMetricsFixture) TestTotalCountCalculation() {
	metric := &HistogramTotalCountMetric{histogram: this.histogram, name: "hello"}
	this.So(metric.Measure(), should.Resemble, MetricMeasurement{
		MetricType: gaugeMetricType,
		Name:       "hello",
		Value:      99999,
	})
}
