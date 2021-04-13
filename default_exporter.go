package metrics

import (
	"fmt"
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
	for _, item := range this.metrics {
		_, _ = fmt.Fprintf(response, outputFormatHelp, item.Name(), item.Description())
		_, _ = fmt.Fprintf(response, outputFormatType, item.Name(), item.Type())
		if item.Type() == "histogram" {
			buckets := item.(Histogram).Buckets()
			for _, bucket := range buckets {
				_, _ = fmt.Fprintf(response, outputFormatLabels, item.Name()+"_bucket",
					formatBucketLabels(bucket.key, item.Labels()), *bucket.value)
			}
			// "A histogram must have a bucket with {le="+Inf"}. Its value must be identical to the value of x_count."
			// https://prometheus.io/docs/instrumenting/exposition_formats/#histograms-and-summaries
			_, _ = fmt.Fprintf(response, outputFormatLabels, item.Name()+"_bucket",
				formatBucketLabels(-1, item.Labels()), *item.(Histogram).Count())
			_, _ = fmt.Fprintf(response, "%s_count%s %d\n", item.Name(), item.Labels(), *item.(Histogram).Count())
			_, _ = fmt.Fprintf(response, "%s_sum%s %f\n", item.Name(), item.Labels(), *item.(Histogram).Sum())
		} else {
			_, _ = fmt.Fprintf(response, outputFormatLabels, item.Name(), item.Labels(), item.Value())
		}
	}
}

const outputFormatHelp = "\n# HELP %s %s\n"
const outputFormatType = "# TYPE %s %s\n"
const outputFormatLabels = "%s%s %d\n"

func formatBucketLabels(bucket float64, labels string) string {
	var bucketString string
	if bucket == -1 {
		bucketString = `{ le="+Inf"`
	} else {
		bucketString = fmt.Sprintf(`{ le="%5.3f"`, bucket)
	}
	if labels == "" {
		return fmt.Sprintf(`%s }`, bucketString)
	}
	return fmt.Sprintf("%s, %s", bucketString, strings.Replace(labels, "{ ", "", 1))
}
