package walker

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/alitto/pond"
	"go.uber.org/ratelimit"
)

type Source[T any] func(start, fetchCount int) (T, error)
type Sink[T any] func(result T, stop func()) error
type Limiter func() int

type Pagination interface {
	StartIndex(batchStart, workerNumber, batchSize int) int
	FetchCount(limit, start, batchSize int) int
}

type FailedTask struct {
	Start      int
	FetchCount int
	Err        error
}

type Walker[T any] struct {
	source           Source[T]
	sink             Sink[T]
	isStopped        int32
	rateLimiter      ratelimit.Limiter
	sourcePool       *pond.WorkerPool
	sinkPool         *pond.WorkerPool
	failedTasks      []FailedTask
	failedTasksMutex sync.Mutex
	*config
}

func New[T any](source Source[T], sink Sink[T], options ...Option) *Walker[T] {
	config := &config{
		maxBatchSize: 10,
		parallelism:  runtime.NumCPU(),
		limiter:      InfiniteLimiter(),
		pagination:   OffsetPagination{},
		rateLimit:    defaultRateLimiter,
	}

	for _, option := range options {
		option(config)
	}

	if config.context == nil {
		WithContext(context.Background())(config)
	}

	sourcePoolBuffer := 0
	sinkPoolBuffer := 0
	sourcePool := pond.New(config.parallelism, sourcePoolBuffer, pond.MinWorkers(config.parallelism), pond.Context(config.context))
	sinkPool := pond.New(config.parallelism*2, sinkPoolBuffer, pond.Context(config.context))

	walker := &Walker[T]{
		config:      config,
		source:      source,
		sink:        sink,
		rateLimiter: ratelimit.NewUnlimited(),
		sourcePool:  sourcePool,
		sinkPool:    sinkPool,
		failedTasks: make([]FailedTask, 0),
	}

	if config.rateLimit != defaultRateLimiter {
		walker.rateLimiter = ratelimit.New(config.rateLimit.count, ratelimit.Per(config.rateLimit.per))
	}

	return walker
}

func (w *Walker[T]) Walk() {
	w.submitTasks()
	w.sourcePool.StopAndWait()
	w.sinkPool.StopAndWait()
}

func (w *Walker[T]) submitTasks() {
	limit := w.limiter()
	batch := NewBatch(w.maxBatchSize, limit, w.parallelism)

	for batchIndex := 0; batchIndex < batch.Count; batchIndex++ {
		for workerNumber := 0; workerNumber < w.parallelism; workerNumber++ {
			w.rateLimiter.Take()

			if w.IsStopped() {
				return
			}

			batchStart := w.parallelism * batchIndex
			start := w.pagination.StartIndex(batchStart, workerNumber, batch.Size)
			fetchCount := w.pagination.FetchCount(limit, start, batch.Size)
			w.submitTask(start, fetchCount)
		}
	}
}

func (w *Walker[T]) submitTask(start, fetchCount int) {
	if fetchCount == 0 {
		return
	}

	w.sourcePool.Submit(func() {
		result, err := w.source(start, fetchCount)
		if err != nil {
			w.storeFailedTask(start, fetchCount, err)
		}
		w.sinkPool.Submit(func() {
			err = w.sink(result, w.Stop)
			if err != nil {
				w.storeFailedTask(start, fetchCount, err)
			}
		})
	})
}

func (w *Walker[T]) storeFailedTask(start, fetchCount int, err error) {
	w.failedTasksMutex.Lock()
	defer w.failedTasksMutex.Unlock()
	w.failedTasks = append(w.failedTasks, FailedTask{Start: start, FetchCount: fetchCount, Err: err})
}

func (w *Walker[T]) FailedTasks() []FailedTask {
	return w.failedTasks
}

func (w *Walker[T]) Stop() {
	atomic.StoreInt32(&w.isStopped, 1)
}

func (w *Walker[T]) IsStopped() bool {
	return atomic.LoadInt32(&w.isStopped) == 1
}
