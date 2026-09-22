#### SMARTY DISCLAIMER: Subject to the terms of the associated license agreement, this software is freely available for your use. This software is FREE, AS IN PUPPIES, and is a gift. Enjoy your new responsibility. This means that while we may consider enhancement requests, we may or may not choose to entertain requests at our sole and absolute discretion.

[![Build](https://github.com/smarty/metrics/actions/workflows/build.yml/badge.svg)](https://github.com/smarty/metrics/actions/workflows/build.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/smarty/metrics/v2.svg)](https://pkg.go.dev/github.com/smarty/metrics/v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/smarty/metrics/v2)](https://goreportcard.com/report/github.com/smarty/metrics/v2)

# metrics

A minimal, dependency-free Prometheus client for Go. It provides counters,
gauges, and histograms backed by atomic operations, plus an `http.Handler`
that renders them in the Prometheus text exposition format.

## Installation

```
go get github.com/smarty/metrics/v2
```

## Usage

```go
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/smarty/metrics/v2"
)

func main() {
	counter := metrics.NewCounter("my_counter",
		metrics.Options.Description("this is a description"),
		metrics.Options.Label("label_key", "label_value"),
	)

	exporter := metrics.NewExporter()
	exporter.Add(counter)

	go func() {
		for {
			counter.Increment()
			time.Sleep(time.Millisecond * 10)
		}
	}()

	server := &http.Server{Addr: "127.0.0.1:8080", Handler: exporter}
	log.Printf("[INFO] Listening for HTTP traffic on [%s]", server.Addr)
	_ = server.ListenAndServe()
}
```

Run this example with `go run ./sample`, then `curl 127.0.0.1:8080`.

## Metric types

| Constructor    | Interface   | Operations                                    |
|----------------|-------------|-----------------------------------------------|
| `NewCounter`   | `Counter`   | `Increment`, `IncrementN`, `Value`            |
| `NewGauge`     | `Gauge`     | `Increment`, `IncrementN`, `Measure`, `Value` |
| `NewHistogram` | `Histogram` | `Measure`, `Buckets`, `Value`, `Count`, `Sum` |

All constructors take a name followed by any number of options.

## Options

| Option                          | Effect                                                              |
|---------------------------------|---------------------------------------------------------------------|
| `Options.Description(text)`     | Sets the `# HELP` text.                                             |
| `Options.Label(key, value)`     | Adds a label. May be repeated.                                      |
| `Options.Bucket(bounds...)`     | Adds histogram bucket upper bounds. Sorted automatically.           |
| `Options.Exporter(exporter)`    | Registers the metric with the exporter at construction time.        |

Passing `Options.Exporter` removes the need to call `exporter.Add` separately:

```go
exporter := metrics.NewExporter()
requests := metrics.NewCounter("http_requests_total",
	metrics.Options.Exporter(exporter),
)
```

## Histograms

Buckets follow Prometheus cumulative semantics: a measurement increments every
bucket whose upper bound is greater than or equal to the value. The `+Inf`
bucket, `_count`, and `_sum` series are emitted automatically.

```go
latency := metrics.NewHistogram("request_duration_ms",
	metrics.Options.Description("request latency in milliseconds"),
	metrics.Options.Bucket(10, 50, 100, 500),
)
latency.Measure(42)
```

Renders as:

```
# HELP request_duration_ms request latency in milliseconds
# TYPE request_duration_ms histogram
request_duration_ms_bucket{ le="10" } 0
request_duration_ms_bucket{ le="50" } 1
request_duration_ms_bucket{ le="100" } 1
request_duration_ms_bucket{ le="500" } 1
request_duration_ms_bucket{ le="+Inf" } 1
request_duration_ms_count 1
request_duration_ms_sum 42
```

## Concurrency

Every type in this package is safe for concurrent use. Metric values are
updated with `sync/atomic`. The exporter may have metrics added to it from any
goroutine, including while it is serving a scrape.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

[MIT](LICENSE.md)
