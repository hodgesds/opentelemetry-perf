package opentelemetryperf

import (
	"context"
	"syscall"
	"testing"

	perf "github.com/hodgesds/perf-utils"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/sdk/trace"
)

func TestPerfProfiler(t *testing.T) {
	poolable := func(ctx context.Context) error {
		syscall.Getcwd([]byte{})
		syscall.Getcwd([]byte{})
		syscall.Getcwd([]byte{})
		syscall.Getcwd([]byte{})
		return nil
	}
	eventAttr, err := perf.TracepointEventAttr("syscalls", "sys_enter_getcwd")
	require.Nil(t, err)
	builder, err := NewPerfProfilerBuilder(
		[]EventAttrConfig{
			{
				EventAttr: *eventAttr,
				SpanKey:   "sys_enter_getcwd",
			}},
		false,
	)
	require.Nil(t, err)
	m := NewPoolManager(WithProfilerBuilder(builder), WithTracer(trace.Register()))
	require.Nil(t, m.Start())
	p := m.Pool(poolable)
	p(context.Background())
	require.Nil(t, m.Stop())
}

func BenchmarkPoolManagerPerfProfiler(b *testing.B) {
	poolable := func(ctx context.Context) error {
		return nil
	}
	eventAttr, err := perf.TracepointEventAttr("syscalls", "sys_enter_getcwd")
	require.Nil(b, err)
	builder, err := NewPerfProfilerBuilder(
		[]EventAttrConfig{
			{
				EventAttr: *eventAttr,
				SpanKey:   "sys_enter_getcwd",
			},
			{
				EventAttr: *eventAttr,
				SpanKey:   "sys_enter_fork",
			},
			{
				EventAttr: *eventAttr,
				SpanKey:   "sys_enter_gettid",
			},
			{
				EventAttr: *eventAttr,
				SpanKey:   "sys_enter_brk",
			},
			{
				EventAttr: *eventAttr,
				SpanKey:   "sys_enter_getcpu",
			},
		},
		false,
	)
	if err != nil {
		b.Fatal(err)
	}
	m := NewPoolManager(WithProfilerBuilder(builder), WithTracer(trace.Register()))
	m.Start()
	p := m.Pool(poolable)
	ctx := context.Background()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(ctx)
	}
}
