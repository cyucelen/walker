package main

import (
	"fmt"

	"github.com/cyucelen/walker"
)

func source(start, fetchCount int) ([]int, error) {
	return []int{start, fetchCount}, nil
}

func sink(result []int, stop func()) error {
	fmt.Println(result)
	return nil
}

func main() {
	w := walker.New(
		source,
		sink,
		walker.WithLimiter(walker.ConstantLimiter(100)),
		walker.WithMaxBatchSize(10),
		walker.WithParallelism(1),
		walker.WithPagination(walker.OffsetPagination{}),
	)
	w.Walk()
	// Output (order not guaranteed):
	// [0 10]
	// [1 10]
	// [2 10]
	// [3 10]
	// [4 10]
	// [5 10]
	// [7 10]
	// [8 10]
	// [6 10]
	// [9 10]
}
