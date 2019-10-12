package opentelemetryperf

import (
	"context"

	apitrace "go.opentelemetry.io/api/trace"
)

// Poolable is a function that can be pooled across goroutines.
type Poolable func(context.Context) error

// PoolManager is used to manage a pool of Poolable function.
type PoolManager interface {
	Start() error
	Stop() error
	Pool(Poolable) Poolable
}

// Profiler is used to profile a Poolable function
type Profiler interface {
	Profile(context.Context, apitrace.Span, Poolable) error
}

// ProfilerBuilder is a builder for producing Profilers.
type ProfilerBuilder interface {
	Build(int, int) (Profiler, error)
}
