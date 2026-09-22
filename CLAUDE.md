# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```
make test                                   # go fmt, then go test -timeout=1s -race -covermode=atomic ./...
make compile                                # go build ./...
make build                                  # test + compile (this is what CI runs on every push)
go test -race -run TestExporter -count=1 .  # a single test in the root package
go run ./sample                             # demo server exposing one counter on 127.0.0.1:8080
```

The one-second test timeout is deliberate. Tests must stay fast; anything that
sleeps or spins up real network listeners will fail the suite.

## Architecture

Single Go package `metrics` (module `github.com/smarty/metrics/v2`) that
implements a minimal, dependency-free Prometheus client: three metric types and
an HTTP handler that renders them in the Prometheus text exposition format
(version 0.0.4). The `sample/` directory is a runnable demo, not part of the
library.

**Interfaces vs. implementations.** `interfaces.go` defines `Metric`,
`Counter`, `Gauge`, `Histogram`, and `Exporter`. Each has one concrete type in
its own `default_*.go` file, constructed via `NewCounter`, `NewGauge`,
`NewHistogram`, and `NewExporter`. Consumers only ever see the interfaces.

**Functional options through a singleton.** `config.go` exposes
`metrics.Options`, a zero-size struct whose methods return `option` funcs
(`Options.Description`, `Options.Label`, `Options.Bucket`, `Options.Exporter`).
Constructors build a `configuration`, apply the options, then copy the results
into the metric. Two non-obvious consequences:

- `Options.Exporter(e)` registers the metric with `e` at construction time, so
  callers need not call `exporter.Add` themselves. Without it the metric is
  registered with an internal `nop` exporter and is never exported unless added
  manually.
- Labels are rendered to their final string (`{ key="value", ... }`) once, at
  construction, by iterating a map. With more than one label the order is not
  deterministic, which is why the exporter test uses exactly one label per
  metric.

**Histogram semantics.** Buckets are cumulative (`le`), matching Prometheus:
`Measure(v)` increments every bucket whose bound is `>= v`. Buckets are sorted
at option-application time. The mandatory `+Inf` bucket is not stored; the
exporter synthesizes it at render time from `Count()`, using `math.MaxUint64`
as its sentinel key.

**Concurrency model.** Counter, gauge, and histogram values are `sync/atomic`
types and are lock-free. The exporter guards its metric slice with a
`sync.RWMutex` and uses copy-on-write: `Add` publishes a fresh slice, and
`ServeHTTP` reads the current slice header under the read lock, then renders
outside the lock. A published slice is never mutated, so it is safe to iterate
unlocked. Preserve that invariant if you touch `Add`.

## Testing conventions

Standard `testing` package only, no assertion libraries. The shared helper is
`assertEqual` in `default_counter_test.go` (uses `reflect.DeepEqual`).
Concurrency is exercised with `sync.WaitGroup` fan-out under `-race`; see
`measureHistogram` and `TestExporterConcurrentAddAndServe` for the pattern.
The root package is at 100% statement coverage and is expected to stay there.
