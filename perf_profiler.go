package opentelemetryperf

import (
	"context"

	perf "github.com/hodgesds/perf-utils"
	"go.opentelemetry.io/api/core"
	apitrace "go.opentelemetry.io/api/trace"
	"golang.org/x/sys/unix"
)

// EventAttrConfig is a configuration for a perf event attr.
type EventAttrConfig struct {
	EventAttr unix.PerfEventAttr
	SpanKey   string
}

func eventAttrsFromConfigs(configs []EventAttrConfig) []unix.PerfEventAttr {
	eventAttrs := make([]unix.PerfEventAttr, len(configs))
	for i, config := range configs {
		eventAttrs[i] = config.EventAttr
	}
	return eventAttrs
}

// PerfProfiler is a profiler that utilize perf for annotating spans.
type PerfProfiler struct {
	profiler    perf.GroupProfiler
	configs     []EventAttrConfig
	returnOnErr bool
}

// Profile implements the Profiler interface.
func (p *PerfProfiler) Profile(
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
	if err != nil {
		if p.returnOnErr {
			return err
		}
	}

	profile, err2 := p.profiler.Profile()
	if err2 != nil {
		return err2
	}

	for i, config := range p.configs {
		// The only type that should be supported is a uint64.
		span.SetAttribute(core.KeyValue{
			Key: core.Key{Name: config.SpanKey},
			Value: core.Value{
				Type:   core.UINT64,
				Uint64: profile.Values[i],
			},
		})
	}
	return err
}

// PerfProfilerBuilder is a builder that creates a PerfProfiler.
type PerfProfilerBuilder struct {
	configs     []EventAttrConfig
	returnOnErr bool
}

// NewPerfProfilerBuilder returns a PerfProfilerBuilder.
func NewPerfProfilerBuilder(
	configs []EventAttrConfig,
	returnOnErr bool,
) (ProfilerBuilder, error) {
	return &PerfProfilerBuilder{
		configs:     configs,
		returnOnErr: returnOnErr,
	}, nil
}

// Build implments the ProfilerBuilder interface.
func (b *PerfProfilerBuilder) Build(threadID, cpu int) (Profiler, error) {
	configs := eventAttrsFromConfigs(b.configs)
	p, err := perf.NewGroupProfiler(threadID, cpu, 0, configs...)
	if err != nil {
		return nil, err
	}
	return &PerfProfiler{
		profiler:    p,
		configs:     b.configs,
		returnOnErr: b.returnOnErr,
	}, p.Start()
}
