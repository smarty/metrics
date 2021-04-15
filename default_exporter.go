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
		_, _ = fmt.Fprintf(response, outputFormatHelp, metric.Name(), metric.Description())
		_, _ = fmt.Fprintf(response, outputFormatType, metric.Name(), metric.Type())

		histogram, ok := metric.(Histogram)
		if !ok {
			_, _ = fmt.Fprintf(response, outputFormatLabels, metric.Name(), metric.Labels(), metric.Value(0))
		} else {
			metricBucketName := histogram.Name() + "_bucket"

			for _, key := range histogram.Keys() {
				_, _ = fmt.Fprintf(response, outputFormatLabels, metricBucketName, formatBucketLabels(key, histogram.Labels()), histogram.Value(key))
			}

			// "A histogram must have a bucket with {le="+Inf"}. Its value must be identical to the value of x_count."
			// https://prometheus.io/docs/instrumenting/exposition_formats/#histograms-and-summaries
			_, _ = fmt.Fprintf(response, outputFormatLabels, metricBucketName, formatBucketLabels(math.MaxInt64, histogram.Labels()), histogram.(Histogram).Count())

			_, _ = fmt.Fprintf(response, outputFormatHistogramSum,
				histogram.Name(), histogram.Labels(), histogram.Count(),
				histogram.Name(), histogram.Labels(), histogram.Sum())
		}
	}
}

func formatBucketLabels(bucket int64, labels string) string {
	var bucketString string
	if bucket == math.MaxInt64 {
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
