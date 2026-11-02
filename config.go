package metrics

import (
	"fmt"
	"maps"
	"net/http"
	"slices"
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
func (singleton) Bucket(values ...uint64) option {
	return func(this *configuration) {
		this.Buckets = append(this.Buckets, values...)
		sort.Slice(this.Buckets, func(i, j int) bool { return this.Buckets[i] < this.Buckets[j] })
	}
}

func (singleton) apply(options ...option) option {
	return func(this *configuration) {
		this.Labels = map[string]string{}
		for _, item := range Options.defaults(options...) {
			item(this)
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

	// Sorted keys render every series the same way on every start, so output can be
	// diffed and matched as text; Prometheus itself ignores label order.
	for _, key := range slices.Sorted(maps.Keys(this.Labels)) {
		result += fmt.Sprintf(`%s="%s", `, key, this.Labels[key])
	}
	result = strings.TrimSuffix(result, ", ")
	return fmt.Sprintf("{ %s }", result)
}

///////////////////////////////////////////////////////////////////////////////

type nop struct{}

func (nop) Add(...Metric)                                {}
func (nop) ServeHTTP(http.ResponseWriter, *http.Request) {}
