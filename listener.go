package metrics

import "io"

type ListenCloser interface {
	Listen()
	io.Closer
}

func NewListener(key string, options ...option) ListenCloser {
	this := defaultListener{key: key}
	Options.apply(options...)(&this)
	return this
}

type defaultListener struct {
	key      string
	tags     string
	capacity int
	writers  int
}

func (this defaultListener) Listen() {
	InitializeEnvironmentTags(this.tags)
	_ = StartAppOptics(this.key, this.capacity, this.writers)
}

func (this defaultListener) Close() error { return nil }

//////

var Options singleton

type singleton struct{}
type option func(*defaultListener)

func (singleton) Tags(value string) option {
	return func(this *defaultListener) { this.tags = value }
}
func (singleton) QueueCapacity(value int) option {
	return func(this *defaultListener) { this.capacity = value }
}
func (singleton) ConcurrentWriters(value int) option {
	return func(this *defaultListener) { this.writers = value }
}

func (singleton) apply(options ...option) option {
	return func(this *defaultListener) {
		for _, option := range Options.defaults(options...) {
			option(this)
		}
	}
}
func (singleton) defaults(options ...option) []option {
	const defaultTags = ""
	const defaultQueueCapacity = 1024
	const defaultConcurrentWorkers = 4

	return append([]option{
		Options.Tags(defaultTags),
		Options.QueueCapacity(defaultQueueCapacity),
		Options.ConcurrentWriters(defaultConcurrentWorkers),
	}, options...)
}
