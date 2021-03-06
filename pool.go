package opentelemetryperf

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	apitrace "go.opentelemetry.io/api/trace"
	"golang.org/x/sys/unix"
)

// errPool is for pooling error channels to reduce allocations.
var errPool = sync.Pool{
	New: func() interface{} {
		return make(chan error)
	},
}

type managerOpt func(*manager)

// WithTracer sets the tracer on a PoolManager.
func WithTracer(tracer apitrace.Tracer) managerOpt {
	return func(m *manager) {
		for _, worker := range m.workers {
			worker.tracer = tracer
		}
	}
}

// WithProfilerBuilder sets the Profiler on a PoolManager.
func WithProfilerBuilder(builder ProfilerBuilder) managerOpt {
	return func(m *manager) {
		for _, worker := range m.workers {
			worker.builder = builder
		}
	}
}

// NewPoolManager returns a PoolManager.
func NewPoolManager(opts ...managerOpt) PoolManager {
	ncpus := runtime.NumCPU()
	workers := make([]*poolWorker, ncpus)
	for i := 0; i < ncpus; i++ {
		workers[i] = &poolWorker{
			cpu:      i,
			tracer:   apitrace.NoopTracer{},
			builder:  &NoopProfilerBuilder{},
			stopChan: make(chan struct{}),
		}
	}

	m := &manager{
		workers: workers,
		work:    make(chan func() (context.Context, chan error, Poolable)),
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

type manager struct {
	workersMu sync.RWMutex
	workers   []*poolWorker
	work      chan func() (context.Context, chan error, Poolable)
}

// Start implements the PoolManager interface.
func (m *manager) Start() error {
	m.workersMu.RLock()
	defer m.workersMu.RUnlock()
	for _, worker := range m.workers {
		go worker.run(m.work)
	}
	return nil
}

// Stop implements the PoolManager interface.
func (m *manager) Stop() error {
	m.workersMu.Lock()
	defer m.workersMu.Unlock()
	for _, worker := range m.workers {
		worker.stop()
	}
	return nil
}

// Pool is used to pool a Poolable function.
func (m *manager) Pool(p Poolable) Poolable {
	return func(ctx context.Context) error {
		err := errPool.Get().(chan error)
		m.work <- func() (context.Context, chan error, Poolable) {
			return ctx, err, p
		}
		return <-err
	}
}

type poolWorker struct {
	tracer   apitrace.Tracer
	builder  ProfilerBuilder
	cpu      int
	stopChan chan struct{}
}

func (w *poolWorker) run(
	work <-chan func() (context.Context, chan error, Poolable),
) error {
	runtime.LockOSThread()
	threadID := unix.Gettid()
	cpuSet := unix.CPUSet{}
	cpuSet.Set(w.cpu)
	unix.SchedSetaffinity(0, &cpuSet)
	profiler, err := w.builder.Build(threadID, w.cpu)
	if err != nil {
		// TODO(hodges): Log errors.
		return err
	}

	spanName := fmt.Sprintf("pool-%d", threadID)
	for {
		select {
		case item := <-work:
			ctx, err, fn := item()
			ctx, span := w.tracer.Start(ctx, spanName, apitrace.WithRecordEvents())
			err <- profiler.Profile(ctx, span, fn)
			span.End()
			errPool.Put(err)
		case <-w.stopChan:
			// TODO(hodgesds): Unset thread affinity?
			runtime.UnlockOSThread()
			return nil
		}
	}
}

func (w *poolWorker) stop() {
	w.stopChan <- struct{}{}
}
