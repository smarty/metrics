package metrics

import (
	"fmt"
	"strings"
)

var Options singleton

type singleton struct{}
type option func(*configuration)
type configuration struct {
	Name        string
	Description string
	Labels      map[string]string
	Buckets     []uint64
}

func (singleton) Description(value string) option {
	return func(this *configuration) { this.Description = value }
}
func (singleton) Label(key, value string) option {
	return func(this *configuration) { this.Labels[key] = value }
}
func (singleton) Bucket(value uint64) option {
	return func(this *configuration) {
		this.Buckets = append(this.Buckets, value)
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
	return append([]option{}, options...)
}

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
