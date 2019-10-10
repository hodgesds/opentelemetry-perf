package opentelemetryperf

import (
	"context"
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/sdk/trace"
)

func TestPoolManager(t *testing.T) {
	poolable := func(ctx context.Context) {}
	m := NewPoolManager(WithTracer(trace.Register()))
	require.Nil(t, m.Start())
	p := m.Pool(poolable)
	p(context.Background())
	require.Nil(t, m.Stop())
}

func BenchmarkPoolManager(b *testing.B) {
	poolable := func(ctx context.Context) {}
	m := NewPoolManager(WithTracer(trace.Register()))
	m.Start()
	p := m.Pool(poolable)
	ctx := context.Background()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(ctx)
	}
}

func BenchmarkPoolable(b *testing.B) {
	poolable := func(ctx context.Context) {}
	ctx := context.Background()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		poolable(ctx)
	}
}

func TestPerfProfiler(t *testing.T) {
	poolable := func(ctx context.Context) {
		syscall.Getcwd([]byte{})
		syscall.Getcwd([]byte{})
		syscall.Getcwd([]byte{})
		syscall.Getcwd([]byte{})
	}
	builder, err := NewHardwareProfilerBuilder()
	require.Nil(t, err)
	m := NewPoolManager(WithProfilerBuilder(builder), WithTracer(trace.Register()))
	require.Nil(t, m.Start())
	p := m.Pool(poolable)
	p(context.Background())
	require.Nil(t, m.Stop())
}

func BenchmarkPoolManagerPerf(b *testing.B) {
	poolable := func(ctx context.Context) {
	}
	builder, err := NewHardwareProfilerBuilder()
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
