package metrics2

import (
	"fmt"
	"net/http"
)

type defaultExporter struct {
	metrics []metric
}

func NewExporter() Exporter {
	return &defaultExporter{}
}

func (this *defaultExporter) Add(items ...metric) {
	this.metrics = append(this.metrics, items...)
}

func (this *defaultExporter) ServeHTTP(response http.ResponseWriter, _ *http.Request) {
	response.Header().Set("Content-Type", "text/plain; version=0.0.4")
	for _, item := range this.metrics {
		_, _ = fmt.Fprintf(response, outputFormat,
			item.Name(), item.Description(),
			item.Name(), item.Type(),
			item.Name(), item.Labels(), item.Value())
	}
}

const outputFormat = `
# HELP %s %s
# TYPE %s %s
%s%s %d
`
