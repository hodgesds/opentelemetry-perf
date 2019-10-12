package main

import (
	"context"
	"log"
	"syscall"
	"time"

	otperf "github.com/hodgesds/opentelemetry-perf"
	perf "github.com/hodgesds/perf-utils"
	"go.opentelemetry.io/exporter/trace/stdout"
	"go.opentelemetry.io/sdk/trace"
	sdktrace "go.opentelemetry.io/sdk/trace"
)

func initTracer() {
	sdktrace.Register()

	exporter, err := stdout.NewExporter(stdout.Options{PrettyPrint: true})
	if err != nil {
		log.Fatal(err)
	}

	ssp := sdktrace.NewSimpleSpanProcessor(exporter)
	sdktrace.RegisterSpanProcessor(ssp)

	// For the demonstration, use sdktrace.AlwaysSample sampler to sample all traces.
	sdktrace.ApplyConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()})
}

func main() {
	initTracer()

	ctx := context.Background()

	eventAttr, err := perf.TracepointEventAttr("syscalls", "sys_enter_getcwd")
	if err != nil {
		log.Fatal(err)
	}
	builder, err := otperf.NewPerfProfilerBuilder(
		[]otperf.EventAttrConfig{
			{
				EventAttr: *eventAttr,
				SpanKey:   "sys_enter_getcwd",
			}},
		false,
	)
	if err != nil {
		log.Fatal(err)
	}

	m := otperf.NewPoolManager(
		otperf.WithProfilerBuilder(builder),
		otperf.WithTracer(trace.Register()),
	)
	m.Start()

	p := m.Pool(func(ctx context.Context) error {
		syscall.Getcwd([]byte{})
		syscall.Getcwd([]byte{})
		syscall.Getcwd([]byte{})
		_, err := syscall.Getcwd(make([]byte, 100))
		return err
	})

	p(ctx)
	time.Sleep(200 * time.Millisecond)
}
