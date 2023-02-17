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
		walker.WithPagination(walker.CursorPagination{}),
	)
	w.Walk()
	// Output (order not guaranteed):
	// [0 10]
	// [20 10]
	// [30 10]
	// [40 10]
	// [50 10]
	// [60 10]
	// [70 10]
	// [10 10]
	// [80 10]
	// [90 10]
}
