package opentelemetryperf

import (
	"context"

	perf "github.com/hodgesds/perf-utils"
	"go.opentelemetry.io/api/core"
	apitrace "go.opentelemetry.io/api/trace"
)

// HardwareProfiler is a profiler that utilize perf for hardware events.
type HardwareProfiler struct {
	profiler    perf.HardwareProfiler
	returnOnErr bool
}

// Profile implements the Profiler interface.
func (p *HardwareProfiler) Profile(
	ctx context.Context,
	span apitrace.Span,
	poolable Poolable,
) error {
	var err error
	if !span.IsRecordingEvents() {
		return poolable(ctx)
	}
	p.profiler.Reset()
	err = poolable(ctx)
	if err != nil && p.returnOnErr {
		return err
	}

	profile, err := p.profiler.Profile()
	if err != nil {
		return err
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
	return nil
}

// HardwareProfilerBuilder is a builder that creates a HardwareProfiler.
type HardwareProfilerBuilder struct {
	returnOnErr bool
}

// NewHardwareProfilerBuilder returns a HardwareProfilerBuilder.
func NewHardwareProfilerBuilder(returnOnErr bool) (ProfilerBuilder, error) {
	return &HardwareProfilerBuilder{
		returnOnErr: returnOnErr,
	}, nil
}

// Build implments the ProfilerBuilder interface.
func (b *HardwareProfilerBuilder) Build(threadID, cpu int) (Profiler, error) {
	p := perf.NewHardwareProfiler(threadID, cpu)
	return &HardwareProfiler{
		profiler:    p,
		returnOnErr: b.returnOnErr,
	}, p.Start()
}
