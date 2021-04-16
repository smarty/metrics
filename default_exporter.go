package metrics

import (
	"fmt"
	"math"
	"net/http"
	"strings"
)

type defaultExporter struct {
	metrics []Metric
}

func NewExporter() Exporter {
	return &defaultExporter{}
}

func (this *defaultExporter) Add(items ...Metric) {
	this.metrics = append(this.metrics, items...)
}
func (this *defaultExporter) ServeHTTP(response http.ResponseWriter, _ *http.Request) {
	response.Header().Set("Content-Type", "text/plain; version=0.0.4")
	for _, metric := range this.metrics {
		renderMetric(metric, response)
	}
}
func renderMetric(metric Metric, response http.ResponseWriter) {
	_, _ = fmt.Fprintf(response, outputFormatHelp, metric.Name(), metric.Description())
	_, _ = fmt.Fprintf(response, outputFormatType, metric.Name(), metric.Type())

	if counter, isCounter := metric.(Counter); isCounter {
		renderCounter(counter, response)
	} else if gauge, isGauge := metric.(Gauge); isGauge {
		renderGauge(gauge, response)
	} else if histogram, isHistogram := metric.(Histogram); isHistogram {
		renderHistogram(histogram, response)
	}
}
func renderCounter(metric Counter, response http.ResponseWriter) {
	_, _ = fmt.Fprintf(response, outputFormatLabels, metric.Name(), metric.Labels(), metric.Value())
}
func renderGauge(metric Gauge, response http.ResponseWriter) {
	_, _ = fmt.Fprintf(response, outputFormatLabels, metric.Name(), metric.Labels(), metric.Value())
}
func renderHistogram(metric Histogram, response http.ResponseWriter) {
	name, labels := metric.Name(), metric.Labels()
	metricBucketName := name + "_bucket"
	sum, count := metric.Sum(), metric.Count()

	for _, key := range metric.Buckets() {
		_, _ = fmt.Fprintf(response, outputFormatLabels, metricBucketName, formatHistogramBucketLabels(key, labels), metric.Value(key))
	}

	// https://prometheus.io/docs/instrumenting/exposition_formats/#histograms-and-summaries
	// > A histogram must have a bucket with {le="+Inf"}. Its value must be identical to the value of x_count.
	_, _ = fmt.Fprintf(response, outputFormatLabels, metricBucketName, formatHistogramBucketLabels(math.MaxUint64, labels), count)
	_, _ = fmt.Fprintf(response, outputFormatHistogramSum, name, labels, count, name, labels, sum)
}
func formatHistogramBucketLabels(bucket uint64, labels string) string {
	var bucketString string
	if bucket == math.MaxUint64 {
		bucketString = `{ le="+Inf"`
	} else {
		bucketString = fmt.Sprintf(`{ le="%d"`, bucket)
	}

	if len(labels) == 0 {
		return fmt.Sprintf(`%s }`, bucketString)
	}

	return fmt.Sprintf("%s, %s", bucketString, strings.Replace(labels, "{ ", "", 1))
}

const (
	outputFormatHelp         = "\n# HELP %s %s\n"
	outputFormatType         = "# TYPE %s %s\n"
	outputFormatLabels       = "%s%s %d\n"
	outputFormatHistogramSum = "%s_count%s %d\n%s_sum%s %d\n"
)
