package metrics

import (
	"fmt"
	"net/http"
	"sort"
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
			keys := make([]float64, 0)
			for label := range buckets {
				keys = append(keys, label)
			}
			sort.Float64s(keys)
			for _, label := range keys {
				_, _ = fmt.Fprintf(response, outputFormatBuckets, item.Name(), label, *buckets[label]) // TODO: Accept multiple label key-pairs
			}
			fmt.Fprintf(response, "%s_count %d\n", item.Name(), *item.(Histogram).Count())
			fmt.Fprintf(response, "%s_sum %f", item.Name(), *item.(Histogram).Sum())
		} else {
			_, _ = fmt.Fprintf(response, outputFormatLabels, item.Name(), item.Labels(), item.Value()) // TODO: Accept multiple label key-pairs
		}
	}
}

const outputFormatHelp = "\n# HELP %s %s\n"
const outputFormatType = "# TYPE %s %s\n"
const outputFormatLabels = "%s%s %d\n"
const outputFormatBuckets = "%s{le=\"%5.3f\"} %d\n"
