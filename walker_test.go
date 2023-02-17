package walker_test

import (
	"sort"
	"sync"
	"testing"

	"github.com/cyucelen/walker"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

type WalkerTestCase struct {
	scenario       string
	shouldStop     func([]int) bool
	sourceFunc     walker.Source[[]int]
	options        []walker.Option
	expectedOutput [][]int
}

func TestWalker(t *testing.T) {
	tests := []WalkerTestCase{
		{
			scenario:   "cursor pagination, limit 100, max batch size 10, parallelism 1",
			sourceFunc: cursorSourceWithUpperbound(100),
			options: []walker.Option{
				walker.WithLimiter(walker.ConstantLimiter(100)),
				walker.WithMaxBatchSize(10),
				walker.WithParallelism(1),
				walker.WithPagination(walker.CursorPagination{}),
			},
			expectedOutput: makeExpectedOutput(100, 10),
		},
		{
			scenario:   "cursor pagination, limit 100, max batch size 12, parallelism 1",
			sourceFunc: cursorSourceWithUpperbound(100),
			options: []walker.Option{
				walker.WithLimiter(walker.ConstantLimiter(100)),
				walker.WithMaxBatchSize(12),
				walker.WithParallelism(1),
				walker.WithPagination(walker.CursorPagination{}),
			},
			expectedOutput: makeExpectedOutput(100, 12),
		},
		{
			scenario:   "cursor pagination, limit 101, max batch size 10, parallelism 1",
			sourceFunc: cursorSourceWithUpperbound(101),
			options: []walker.Option{
				walker.WithLimiter(walker.ConstantLimiter(101)),
				walker.WithMaxBatchSize(10),
				walker.WithParallelism(1),
				walker.WithPagination(walker.CursorPagination{}),
			},
			expectedOutput: makeExpectedOutput(101, 10),
		},
		{
			scenario:   "cursor pagination, limit 99, max batch size 10, parallelism 1",
			sourceFunc: cursorSourceWithUpperbound(99),
			options: []walker.Option{
				walker.WithLimiter(walker.ConstantLimiter(99)),
				walker.WithMaxBatchSize(10),
				walker.WithParallelism(1),
				walker.WithPagination(walker.CursorPagination{}),
			},
			expectedOutput: makeExpectedOutput(99, 10),
		},
		{
			scenario:   "cursor pagination, limit 97, max batch size 23, parallelism 1",
			sourceFunc: cursorSourceWithUpperbound(97),
			options: []walker.Option{
				walker.WithLimiter(walker.ConstantLimiter(97)),
				walker.WithMaxBatchSize(23),
				walker.WithParallelism(1),
				walker.WithPagination(walker.CursorPagination{}),
			},
			expectedOutput: makeExpectedOutput(97, 23),
		},
		{
			scenario:   "cursor pagination, limit 100, max batch size 10, parallelism 10",
			sourceFunc: cursorSourceWithUpperbound(100),
			options: []walker.Option{
				walker.WithLimiter(walker.ConstantLimiter(100)),
				walker.WithMaxBatchSize(10),
				walker.WithParallelism(10),
				walker.WithPagination(walker.CursorPagination{}),
			},
			expectedOutput: makeExpectedOutput(100, 10),
		},
		{
			scenario:   "cursor pagination, limit 100, max batch size 12, parallelism 8",
			sourceFunc: cursorSourceWithUpperbound(100),
			options: []walker.Option{
				walker.WithLimiter(walker.ConstantLimiter(100)),
				walker.WithMaxBatchSize(12),
				walker.WithParallelism(8),
				walker.WithPagination(walker.CursorPagination{}),
			},
			expectedOutput: makeExpectedOutput(100, 12),
		},
		{
			scenario:   "cursor pagination, limit 100, max batch size 1, parallelism 100",
			sourceFunc: cursorSourceWithUpperbound(100),
			options: []walker.Option{
				walker.WithLimiter(walker.ConstantLimiter(100)),
				walker.WithMaxBatchSize(1),
				walker.WithParallelism(100),
				walker.WithPagination(walker.CursorPagination{}),
			},
			expectedOutput: makeExpectedOutput(100, 1),
		},
		{
			scenario:   "cursor pagination, limit 100, max batch size 2, parallelism 100",
			sourceFunc: cursorSourceWithUpperbound(100),
			options: []walker.Option{
				walker.WithLimiter(walker.ConstantLimiter(100)),
				walker.WithMaxBatchSize(2),
				walker.WithParallelism(100),
				walker.WithPagination(walker.CursorPagination{}),
			},
			expectedOutput: makeExpectedOutput(100, 2),
		},
		{
			scenario:   "cursor pagination, limit 2, max batch size 10, parallelism 10, stop called on empty result",
			sourceFunc: cursorSourceWithUpperbound(2),
			options: []walker.Option{
				walker.WithLimiter(walker.ConstantLimiter(2)),
				walker.WithMaxBatchSize(10),
				walker.WithParallelism(10),
				walker.WithPagination(walker.CursorPagination{}),
			},
			expectedOutput: makeExpectedOutput(2, 10),
		},
		{
			scenario:   "cursor pagination, limit infinite, max batch size 10, parallelism 1, stop called on empty result",
			sourceFunc: cursorSourceWithUpperbound(100),
			options: []walker.Option{
				walker.WithLimiter(walker.InfiniteLimiter()),
				walker.WithMaxBatchSize(10),
				walker.WithParallelism(1),
				walker.WithPagination(walker.CursorPagination{}),
			},
			shouldStop:     isEmpty,
			expectedOutput: makeExpectedOutput(100, 10),
		},
		{
			scenario:   "cursor pagination, limit infinite, max batch size 12, parallelism 1, stop called on empty result",
			sourceFunc: cursorSourceWithUpperbound(100),
			options: []walker.Option{
				walker.WithLimiter(walker.InfiniteLimiter()),
				walker.WithMaxBatchSize(12),
				walker.WithParallelism(1),
				walker.WithPagination(walker.CursorPagination{}),
			},
			shouldStop:     isEmpty,
			expectedOutput: makeExpectedOutput(100, 12),
		},
		{
			scenario:   "cursor pagination, limit infinite, max batch size 10, parallelism 10, stop called on empty result",
			sourceFunc: cursorSourceWithUpperbound(100),
			options: []walker.Option{
				walker.WithLimiter(walker.InfiniteLimiter()),
				walker.WithMaxBatchSize(10),
				walker.WithParallelism(10),
				walker.WithPagination(walker.CursorPagination{}),
			},
			shouldStop:     isEmpty,
			expectedOutput: makeExpectedOutput(100, 10),
		},
		{
			scenario:   "cursor pagination, limit infinite, max batch size 12, parallelism 8, stop called on empty result",
			sourceFunc: cursorSourceWithUpperbound(100),
			options: []walker.Option{
				walker.WithLimiter(walker.InfiniteLimiter()),
				walker.WithMaxBatchSize(12),
				walker.WithParallelism(8),
				walker.WithPagination(walker.CursorPagination{}),
			},
			shouldStop:     isEmpty,
			expectedOutput: makeExpectedOutput(100, 12),
		},
		//
		{
			scenario:   "offset pagination, limit 100, max batch size 10, parallelism 1",
			sourceFunc: offsetSourceWithUpperbound(100),
			options: []walker.Option{
				walker.WithLimiter(walker.ConstantLimiter(100)),
				walker.WithMaxBatchSize(10),
				walker.WithParallelism(1),
				walker.WithPagination(walker.OffsetPagination{}),
			},
			expectedOutput: makeExpectedOutput(100, 10),
		},
		{
			scenario:   "offset pagination, limit 100, max batch size 12, parallelism 1",
			sourceFunc: offsetSourceWithUpperbound(100),
			options: []walker.Option{
				walker.WithLimiter(walker.ConstantLimiter(100)),
				walker.WithMaxBatchSize(12),
				walker.WithParallelism(1),
				walker.WithPagination(walker.OffsetPagination{}),
			},
			expectedOutput: makeExpectedOutput(100, 12),
		},
		{
			scenario:   "offset pagination, limit 101, max batch size 10, parallelism 1",
			sourceFunc: offsetSourceWithUpperbound(101),
			options: []walker.Option{
				walker.WithLimiter(walker.ConstantLimiter(101)),
				walker.WithMaxBatchSize(10),
				walker.WithParallelism(1),
				walker.WithPagination(walker.OffsetPagination{}),
			},
			expectedOutput: makeExpectedOutput(101, 10),
		},
		{
			scenario:   "offset pagination, limit 99, max batch size 10, parallelism 1",
			sourceFunc: offsetSourceWithUpperbound(99),
			options: []walker.Option{
				walker.WithLimiter(walker.ConstantLimiter(99)),
				walker.WithMaxBatchSize(10),
				walker.WithParallelism(1),
				walker.WithPagination(walker.OffsetPagination{}),
			},
			expectedOutput: makeExpectedOutput(99, 10),
		},
		{
			scenario:   "offset pagination, limit 97, max batch size 23, parallelism 1",
			sourceFunc: offsetSourceWithUpperbound(97),
			options: []walker.Option{
				walker.WithLimiter(walker.ConstantLimiter(97)),
				walker.WithMaxBatchSize(23),
				walker.WithParallelism(1),
				walker.WithPagination(walker.OffsetPagination{}),
			},
			expectedOutput: makeExpectedOutput(97, 23),
		},
		{
			scenario:   "offset pagination, limit 100, max batch size 10, parallelism 10",
			sourceFunc: offsetSourceWithUpperbound(100),
			options: []walker.Option{
				walker.WithLimiter(walker.ConstantLimiter(100)),
				walker.WithMaxBatchSize(10),
				walker.WithParallelism(10),
				walker.WithPagination(walker.OffsetPagination{}),
			},
			expectedOutput: makeExpectedOutput(100, 10),
		},
		{
			scenario:   "offset pagination, limit 100, max batch size 12, parallelism 8",
			sourceFunc: offsetSourceWithUpperbound(100),
			options: []walker.Option{
				walker.WithLimiter(walker.ConstantLimiter(100)),
				walker.WithMaxBatchSize(12),
				walker.WithParallelism(8),
				walker.WithPagination(walker.OffsetPagination{}),
			},
			expectedOutput: makeExpectedOutput(100, 12),
		},
		{
			scenario:   "offset pagination, limit 100, max batch size 1, parallelism 100",
			sourceFunc: offsetSourceWithUpperbound(100),
			options: []walker.Option{
				walker.WithLimiter(walker.ConstantLimiter(100)),
				walker.WithMaxBatchSize(1),
				walker.WithParallelism(100),
				walker.WithPagination(walker.OffsetPagination{}),
			},
			expectedOutput: makeExpectedOutput(100, 1),
		},
		{
			scenario:   "offset pagination, limit 100, max batch size 2, parallelism 100",
			sourceFunc: offsetSourceWithUpperbound(100),
			options: []walker.Option{
				walker.WithLimiter(walker.ConstantLimiter(100)),
				walker.WithMaxBatchSize(2),
				walker.WithParallelism(100),
				walker.WithPagination(walker.OffsetPagination{}),
			},
			expectedOutput: makeExpectedOutput(100, 2),
		},
		{
			scenario:   "offset pagination, limit infinite, max batch size 10, parallelism 1, stop called on empty result",
			sourceFunc: offsetSourceWithUpperbound(100),
			options: []walker.Option{
				walker.WithLimiter(walker.InfiniteLimiter()),
				walker.WithMaxBatchSize(10),
				walker.WithParallelism(1),
				walker.WithPagination(walker.OffsetPagination{}),
			},
			shouldStop:     isEmpty,
			expectedOutput: makeExpectedOutput(100, 10),
		},
		{
			scenario:   "offset pagination, limit infinite, max batch size 12, parallelism 1, stop called on empty result",
			sourceFunc: offsetSourceWithUpperbound(100),
			options: []walker.Option{
				walker.WithLimiter(walker.InfiniteLimiter()),
				walker.WithMaxBatchSize(12),
				walker.WithParallelism(1),
				walker.WithPagination(walker.OffsetPagination{}),
			},
			shouldStop:     isEmpty,
			expectedOutput: makeExpectedOutput(100, 12),
		},
		{
			scenario:   "offset pagination, limit infinite, max batch size 10, parallelism 10, stop called on empty result",
			sourceFunc: offsetSourceWithUpperbound(100),
			options: []walker.Option{
				walker.WithLimiter(walker.InfiniteLimiter()),
				walker.WithMaxBatchSize(10),
				walker.WithParallelism(10),
				walker.WithPagination(walker.OffsetPagination{}),
			},
			shouldStop:     isEmpty,
			expectedOutput: makeExpectedOutput(100, 10),
		},
		{
			scenario:   "offset pagination, limit infinite, max batch size 12, parallelism 8, stop called on empty result",
			sourceFunc: offsetSourceWithUpperbound(100),
			options: []walker.Option{
				walker.WithLimiter(walker.InfiniteLimiter()),
				walker.WithMaxBatchSize(12),
				walker.WithParallelism(8),
				walker.WithPagination(walker.OffsetPagination{}),
			},
			shouldStop:     isEmpty,
			expectedOutput: makeExpectedOutput(100, 12),
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			mockSink := MockSink{shouldStop: test.shouldStop}
			walker := walker.New(test.sourceFunc, mockSink.sink, test.options...)
			walker.Walk()
			assert.Equal(t, test.expectedOutput, mockSink.sortedResults(), test.scenario)
		})
	}
}

func makeExpectedOutput(limit, batchSize int) [][]int {
	return lo.Chunk(lo.Map(make([]int, limit), func(item, index int) int { return index + 1 }), batchSize)
}

func cursorSource[T []int](limit int) walker.Source[T] {
	return func(start, fetchCount int) (T, error) {
		length := min(limit-start, fetchCount)
		results := make([]int, length)
		for i := range results {
			results[i] = i + start + 1
		}
		return results, nil
	}
}

func cursorSourceWithUpperbound[T []int](limit int) walker.Source[T] {
	return func(start, fetchCount int) (T, error) {
		if start > limit {
			return []int{}, nil
		}

		return cursorSource(limit)(start, fetchCount)
	}
}

func offsetSource[T []int](limit int) walker.Source[T] {
	return func(page, fetchCount int) (T, error) {
		length := min(limit-page*fetchCount, fetchCount)
		results := make([]int, length)
		for i := range results {
			results[i] = fetchCount*page + i + 1
		}
		return results, nil
	}
}

func offsetSourceWithUpperbound[T []int](limit int) walker.Source[T] {
	return func(page, fetchCount int) (T, error) {
		if (page * fetchCount) > limit {
			return []int{}, nil
		}

		return offsetSource(limit)(page, fetchCount)
	}
}

type MockSink struct {
	results    [][]int
	shouldStop func([]int) bool
	sync.Mutex
}

func (m *MockSink) sink(result []int, stop func()) error {
	m.Lock()
	defer m.Unlock()

	if m.results == nil {
		m.results = make([][]int, 0)
	}
	m.results = append(m.results, result)
	if m.shouldStop != nil && m.shouldStop(result) {
		stop()
		return nil

	}
	return nil
}

func (m *MockSink) sortedResults() [][]int {
	copyResults := append([][]int{}, m.results...)
	filtered := lo.Filter(copyResults, func(item []int, index int) bool { return len(item) > 0 })
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i][0] < filtered[j][0]
	})
	return filtered
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func isEmpty(result []int) bool {
	return len(result) == 0
}
