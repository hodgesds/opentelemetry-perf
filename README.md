# Opentelemetry-Perf
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
to be used to profile the goroutine and annotate traces using the perf
subsystem. Spans can also be annotate with
[kprobes](https://www.kernel.org/doc/html/latest/trace/kprobetrace.html) or
potentially eBPF programs. The `Poolable` type is as follows:

```
type Poolable func(context.Context)
```


# Benchmarks
In order to collect a perf [hardware
profile](https://godoc.org/github.com/hodgesds/perf-utils#HardwareProfile) and
attach it to a span it takes ~33500ns on a Xeon E3-1505M.

```
$ go test -cover -v -bench=.                                                                                                                                    
=== RUN   TestPoolManager
--- PASS: TestPoolManager (0.01s)
=== RUN   TestPerfProfiler
--- PASS: TestPerfProfiler (0.02s)
goos: linux
goarch: amd64
pkg: github.com/hodgesds/opentelemetry-perf
BenchmarkPoolManager-8            200000              9963 ns/op             904 B/op         14 allocs/op
BenchmarkPoolable-8             2000000000               0.28 ns/op            0 B/op          0 allocs/op
BenchmarkPoolManagerPerf-8         50000             33528 ns/op            2872 B/op         56 allocs/op
PASS
coverage: 84.4% of statements
ok      github.com/hodgesds/opentelemetry-perf  4.777s  
```

# Notes
This is still alpha software and rather experimental.
