package metrics

import (
	"fmt"
	"net/http"
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
		if item.Type() == "Histogram" {
			_, _ = fmt.Fprintf(response, outputFormatMetric, item.Name(), item.(Histogram).Buckets(), item.Value())
		} else {
			_, _ = fmt.Fprintf(response, outputFormatMetric, item.Name(), item.Labels(), item.Value()) // TODO: Accept multiple label key-pairs
		}
	}
}

const outputFormatHelp = "\n# HELP %s %s\n"
const outputFormatType = "# TYPE %s %s\n"
const outputFormatMetric = "%s%s %d\n"
