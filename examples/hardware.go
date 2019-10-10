package main

import (
	"context"
	"log"
	"math/bits"
	"time"

	otperf "github.com/hodgesds/opentelemetry-perf"
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

	builder, err := otperf.NewHardwareProfilerBuilder()
	if err != nil {
		log.Fatal(err)
	}

	m := otperf.NewPoolManager(
		otperf.WithProfilerBuilder(builder),
		otperf.WithTracer(trace.Register()),
	)
	m.Start()

	// Pool a function that calculates the number of leading zero bits from
	// 1-100.
	p := m.Pool(func(ctx context.Context) {
		for i := 1; i <= 100; i++ {
			bits.LeadingZeros(uint(i))
		}
	})

	p(ctx)
	time.Sleep(100 * time.Millisecond)
}
