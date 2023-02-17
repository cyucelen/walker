package walker

import "math"

type Batch struct {
	Size  int
	Count int
}

func NewBatch(batchSize, limit, parallelism int) Batch {
	return Batch{
		Size:  batchSize,
		Count: int(math.Ceil((float64(limit) / float64(batchSize)) / float64(parallelism))),
	}
}
