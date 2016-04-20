package metrics

import (
	"fmt"
	"time"
)

type (
	HistogramMinMetric struct {
		*ReportingFrequency
		name      string
		histogram Histogram
	}
	HistogramMaxMetric struct {
		*ReportingFrequency
		name      string
		histogram Histogram
	}
	HistogramMeanMetric struct {
		*ReportingFrequency
		name      string
		histogram Histogram
	}
	HistogramStandardDeviationMetric struct {
		*ReportingFrequency
		name      string
		histogram Histogram
	}
	HistogramQuantileMetric struct {
		*ReportingFrequency
		name      string
		quantile  float64
		histogram Histogram
	}
	HistogramTotalCountMetric struct {
		*ReportingFrequency
		name      string
		histogram Histogram
	}
)

func NewHistogramMinMetric(basename string, histogram Histogram, update time.Duration) *HistogramMinMetric {
	return &HistogramMinMetric{
		ReportingFrequency: &ReportingFrequency{interval: update},
		name:               basename + "_min",
		histogram:          histogram,
	}
}
func NewHistogramMaxMetric(basename string, histogram Histogram, update time.Duration) *HistogramMaxMetric {
	return &HistogramMaxMetric{
		ReportingFrequency: &ReportingFrequency{interval: update},
		name:               basename + "_max",
		histogram:          histogram,
	}
}
func NewHistogramMeanMetric(basename string, histogram Histogram, update time.Duration) *HistogramMeanMetric {
	return &HistogramMeanMetric{
		ReportingFrequency: &ReportingFrequency{interval: update},
		name:               basename + "_mean",
		histogram:          histogram,
	}
}
func NewHistogramStandardDeviationMetric(
	basename string, histogram Histogram, update time.Duration) *HistogramStandardDeviationMetric {

	return &HistogramStandardDeviationMetric{
		ReportingFrequency: &ReportingFrequency{interval: update},
		name:               basename + "_stddev",
		histogram:          histogram,
	}
}
func NewHistogramTotalCountMetric(
	basename string, histogram Histogram, update time.Duration) *HistogramTotalCountMetric {

	return &HistogramTotalCountMetric{
		ReportingFrequency: &ReportingFrequency{interval: update},
		name:               basename + "_total",
		histogram:          histogram,
	}
}
func NewHistogramQuantileMetric(
	basename string, quantile float64, histogram Histogram, update time.Duration) *HistogramQuantileMetric {

	return &HistogramQuantileMetric{
		ReportingFrequency: &ReportingFrequency{interval: update},
		name:               fmt.Sprintf("%s_%.3f", basename, quantile),
		histogram:          histogram,
		quantile:           quantile,
	}
}

func (this *HistogramMinMetric) Measure() MetricMeasurement {
	return MetricMeasurement{
		Name: this.name, MetricType: gaugeMetricType,
		Value: this.histogram.Min(),
	}
}
func (this *HistogramMaxMetric) Measure() MetricMeasurement {
	return MetricMeasurement{
		Name: this.name, MetricType: gaugeMetricType,
		Value: this.histogram.Max(),
	}
}
func (this *HistogramMeanMetric) Measure() MetricMeasurement {
	return MetricMeasurement{
		Name: this.name, MetricType: gaugeMetricType,
		Value: int64(this.histogram.Mean()),
	}
}
func (this *HistogramStandardDeviationMetric) Measure() MetricMeasurement {
	return MetricMeasurement{
		Name: this.name, MetricType: gaugeMetricType,
		Value: int64(this.histogram.StdDev()),
	}
}
func (this *HistogramQuantileMetric) Measure() MetricMeasurement {
	return MetricMeasurement{
		Name: this.name, MetricType: gaugeMetricType,
		Value: this.histogram.ValueAtQuantile(this.quantile),
	}
}
func (this *HistogramTotalCountMetric) Measure() MetricMeasurement {
	return MetricMeasurement{
		Name: this.name, MetricType: gaugeMetricType,
		Value: this.histogram.TotalCount(),
	}
}
