package metrics

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
