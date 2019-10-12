# Opentelemetry-Perf
[![GoDoc](https://godoc.org/github.com/hodgesds/opentelemetry-perf?status.svg)](https://godoc.org/github.com/hodgesds/opentelemetry-perf)

This is a Go package that allows for the use of the Linux perf subsystem with
[OpenTelemetry](https://opentelemetry.io/). It allows for annotating
OpenTelemetry [spans](https://godoc.org/go.opentelemetry.io/api/trace#Span)
with data collected from the perf subsystem.

## System Configuration
See the [perf-utils](https://github.com/hodgesds/perf-utils#setup) setup
instructions for system configuration.

## Pooling functions
This library handles functions that are `Poolable` and does a few special
things. When a function is pooled by a PoolManager it returns a `Poolable`
function that when called will execute on a worker goroutine. The worker
goroutines are locked to an OS thread and bound to a CPU. This allows for the
[`perf_event_open`](http://www.man7.org/linux/man-pages/man2/perf_event_open.2.html)
to be used to profile the thread that the goroutine is executing on and
annotate traces using the perf subsystem. Spans can also be annotate with
[kprobes](https://www.kernel.org/doc/html/latest/trace/kprobetrace.html) or
potentially eBPF programs. The `Poolable` type is as follows:

```
type Poolable func(context.Context) error
```

A `Poolable` function should **not** launch goroutines.

# Benchmarks
In order to collect a perf profile using 5 kprobes and attach it to a span it
takes <17000ns on a Xeon E3-1505M. Note that the number of events being traced
will cause more overhead for span allocations.

```
$ go test -cover -v -bench=.
goos: linux
goarch: amd64
pkg: github.com/hodgesds/opentelemetry-perf
BenchmarkPoolManagerPerfProfiler-8                100000             16677 ns/op            2057 B/op         37 allocs/op
BenchmarkPoolManager-8                            200000              9638 ns/op             808 B/op         13 allocs/op
BenchmarkPoolable-8                             2000000000               0.31 ns/op            0 B/op          0 allocs/op
BenchmarkPoolManagerHardwareProfiler-8             50000             32662 ns/op            2777 B/op         55 allocs/op
PASS
coverage: 90.2% of statements
ok      github.com/hodgesds/opentelemetry-perf  6.610s

```

# Notes
This is still alpha software and rather experimental.
