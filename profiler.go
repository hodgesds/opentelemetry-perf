package opentelemetryperf

import (
	"context"

	perf "github.com/hodgesds/perf-utils"
	"go.opentelemetry.io/api/core"
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
) {
	poolable(ctx)
}

// NoopProfilerBuilder implements a ProfilerBuilder.
type NoopProfilerBuilder struct{}

// Build implments the ProfilerBuilder interface.
func (b *NoopProfilerBuilder) Build(int, int) (Profiler, error) {
	return &NoopProfiler{}, nil
}

// HardwareProfiler is a profiler that utilize perf for hardware events.
type HardwareProfiler struct {
	profiler perf.HardwareProfiler
}

// Profile implements the Profiler interface.
func (p *HardwareProfiler) Profile(
	ctx context.Context,
	span apitrace.Span,
	poolable Poolable,
) {
	if !span.IsRecordingEvents() {
		return
	}
	p.profiler.Reset()
	poolable(ctx)
	profile, err := p.profiler.Profile()
	if err != nil {
		// TODO(hodges): Log errors.
		return
	}
	if profile.CPUCycles != nil {
		span.SetAttribute(core.KeyValue{
			Key: core.Key{
				Name: "CPUCycles",
			},
			Value: core.Value{
				Type:   core.UINT64,
				Uint64: *profile.CPUCycles,
			},
		})
	}
	if profile.Instructions != nil {
		span.SetAttribute(core.KeyValue{
			Key: core.Key{
				Name: "Instructions",
			},
			Value: core.Value{
				Type:   core.UINT64,
				Uint64: *profile.Instructions,
			},
		})
	}
	if profile.CacheRefs != nil {
		span.SetAttribute(core.KeyValue{
			Key: core.Key{
				Name: "CacheRefs",
			},
			Value: core.Value{
				Type:   core.UINT64,
				Uint64: *profile.CacheRefs,
			},
		})
	}
	if profile.CacheMisses != nil {
		span.SetAttribute(core.KeyValue{
			Key: core.Key{
				Name: "cacheMisses",
			},
			Value: core.Value{
				Type:   core.UINT64,
				Uint64: *profile.CacheMisses,
			},
		})
	}
	if profile.BranchInstr != nil {
		span.SetAttribute(core.KeyValue{
			Key: core.Key{
				Name: "BranchInstr",
			},
			Value: core.Value{
				Type:   core.UINT64,
				Uint64: *profile.BranchInstr,
			},
		})
	}
	if profile.BranchMisses != nil {
		span.SetAttribute(core.KeyValue{
			Key: core.Key{
				Name: "BranchMisses",
			},
			Value: core.Value{
				Type:   core.UINT64,
				Uint64: *profile.BranchMisses,
			},
		})
	}
	if profile.BusCycles != nil {
		span.SetAttribute(core.KeyValue{
			Key: core.Key{
				Name: "BusCycles",
			},
			Value: core.Value{
				Type:   core.UINT64,
				Uint64: *profile.BusCycles,
			},
		})
	}
	if profile.StalledCyclesFrontend != nil {
		span.SetAttribute(core.KeyValue{
			Key: core.Key{
				Name: "StalledCyclesFrontend",
			},
			Value: core.Value{
				Type:   core.UINT64,
				Uint64: *profile.StalledCyclesFrontend,
			},
		})
	}
	if profile.StalledCyclesBackend != nil {
		span.SetAttribute(core.KeyValue{
			Key: core.Key{
				Name: "StalledCyclesBackend",
			},
			Value: core.Value{
				Type:   core.UINT64,
				Uint64: *profile.StalledCyclesBackend,
			},
		})
	}
	if profile.RefCPUCycles != nil {
		span.SetAttribute(core.KeyValue{
			Key: core.Key{
				Name: "RefCPUCycles",
			},
			Value: core.Value{
				Type:   core.UINT64,
				Uint64: *profile.RefCPUCycles,
			},
		})
	}
}

// HardwareProfilerBuilder is a builder that creates a HardwareProfiler.
type HardwareProfilerBuilder struct{}

// NewHardwareProfilerBuilder returns a HardwareProfilerBuilder.
func NewHardwareProfilerBuilder() (ProfilerBuilder, error) {
	return &HardwareProfilerBuilder{}, nil
}

// Build implments the ProfilerBuilder interface.
func (b *HardwareProfilerBuilder) Build(threadID, cpu int) (Profiler, error) {
	p := perf.NewHardwareProfiler(threadID, cpu)
	return &HardwareProfiler{profiler: p}, p.Start()
}
