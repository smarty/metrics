package metrics

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
)

var Options singleton

type singleton struct{}
type option func(*configuration)
type configuration struct {
	Name        string
	Description string
	Labels      map[string]string
	Exporter    Exporter
	Buckets     []uint64
}

func (singleton) Description(value string) option {
	return func(this *configuration) { this.Description = value }
}
func (singleton) Label(key, value string) option {
	return func(this *configuration) { this.Labels[key] = value }
}
func (singleton) Exporter(value Exporter) option {
	return func(this *configuration) { this.Exporter = value }
}
func (singleton) Bucket(value uint64) option {
	return func(this *configuration) {
		this.Buckets = append(this.Buckets, value)
		sort.Slice(this.Buckets, func(i, j int) bool { return this.Buckets[i] < this.Buckets[j] })
	}
}

func (singleton) apply(options ...option) option {
	return func(this *configuration) {
		this.Labels = map[string]string{}
		for _, option := range Options.defaults(options...) {
			option(this)
		}
	}
}
func (singleton) defaults(options ...option) []option {
	return append([]option{
		Options.Exporter(nop{}),
	}, options...)
}

///////////////////////////////////////////////////////////////////////////////

func (this configuration) RenderLabels() (result string) {
	if len(this.Labels) == 0 {
		return ""
	}

	for key, value := range this.Labels {
		result += fmt.Sprintf(`%s="%s", `, key, value)
	}
	result = strings.TrimSuffix(result, ", ")
	return fmt.Sprintf("{ %s }", result)
}

///////////////////////////////////////////////////////////////////////////////

type nop struct{}

func (nop) Add(...Metric)                                {}
func (nop) ServeHTTP(http.ResponseWriter, *http.Request) {}
