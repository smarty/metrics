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

func (this *HistogramMinMetric) CalculateMeasurement() MetricMeasurement {
	return MetricMeasurement{
		Name: this.name, MetricType: gaugeMetricType,
		Value: this.histogram.Min(),
	}
}
func (this *HistogramMaxMetric) CalculateMeasurement() MetricMeasurement {
	return MetricMeasurement{
		Name: this.name, MetricType: gaugeMetricType,
		Value: this.histogram.Max(),
	}
}
func (this *HistogramMeanMetric) CalculateMeasurement() MetricMeasurement {
	return MetricMeasurement{
		Name: this.name, MetricType: gaugeMetricType,
		Value: int64(this.histogram.Mean()),
	}
}
func (this *HistogramStandardDeviationMetric) CalculateMeasurement() MetricMeasurement {
	return MetricMeasurement{
		Name: this.name, MetricType: gaugeMetricType,
		Value: int64(this.histogram.StdDev()),
	}
}
func (this *HistogramQuantileMetric) CalculateMeasurement() MetricMeasurement {
	return MetricMeasurement{
		Name: this.name, MetricType: gaugeMetricType,
		Value: this.histogram.ValueAtQuantile(this.quantile),
	}
}
func (this *HistogramTotalCountMetric) CalculateMeasurement() MetricMeasurement {
	return MetricMeasurement{
		Name: this.name, MetricType: gaugeMetricType,
		Value: this.histogram.TotalCount(),
	}
}
