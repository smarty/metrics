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
			_, _ = fmt.Fprintf(response, outputFormatLabels, metric.Name(), metric.Labels(), metric.Value())
		} else {
			bucketKeys := histogram.Buckets()
			bucketValues := histogram.Values()
			for index, bucketKey := range bucketKeys {
				_, _ = fmt.Fprintf(response, outputFormatLabels, metric.Name()+"_bucket",
					formatBucketLabels(bucketKey, metric.Labels()), bucketValues[index])
			}
			// "A histogram must have a bucket with {le="+Inf"}. Its value must be identical to the value of x_count."
			// https://prometheus.io/docs/instrumenting/exposition_formats/#histograms-and-summaries
			_, _ = fmt.Fprintf(response, outputFormatLabels,
				metric.Name()+"_bucket",
				formatBucketLabels(math.MaxUint64, metric.Labels()), metric.(Histogram).Count())

			_, _ = fmt.Fprintf(response, "%s_count%s %d\n", metric.Name(), metric.Labels(), histogram.Count())
			_, _ = fmt.Fprintf(response, "%s_sum%s %d\n", metric.Name(), metric.Labels(), histogram.Sum())
		}
	}
}

const outputFormatHelp = "\n# HELP %s %s\n"
const outputFormatType = "# TYPE %s %s\n"
const outputFormatLabels = "%s%s %d\n"

func formatBucketLabels(bucket uint64, labels string) string {
	var bucketString string
	if bucket == math.MaxUint64 {
		bucketString = `{ le="+Inf"`
	} else {
		bucketString = fmt.Sprintf(`{ le="%5d"`, bucket)
	}
	if labels == "" {
		return fmt.Sprintf(`%s }`, bucketString)
	}
	return fmt.Sprintf("%s, %s", bucketString, strings.Replace(labels, "{ ", "", 1))
}
