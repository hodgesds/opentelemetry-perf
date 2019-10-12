package opentelemetryperf

import (
	"context"

	apitrace "go.opentelemetry.io/api/trace"
)

// NoopProfiler is a profiler that runs the Poolable function with no
// profiling.
type NoopProfiler struct{}

// Profile implements the Profiler interface.
func (p *NoopProfiler) Profile(
	ctx context.Context,
	span apitrace.Span,
	poolable Poolable,
) error {
	return poolable(ctx)
}

// NoopProfilerBuilder implements a ProfilerBuilder.
type NoopProfilerBuilder struct{}

// Build implments the ProfilerBuilder interface.
func (b *NoopProfilerBuilder) Build(int, int) (Profiler, error) {
	return &NoopProfiler{}, nil
}
