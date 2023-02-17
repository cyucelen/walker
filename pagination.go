package walker

type CursorPagination struct{}

func (c CursorPagination) StartIndex(batchStart, workerNumber, batchSize int) int {
	return (batchStart + workerNumber) * batchSize
}

func (c CursorPagination) FetchCount(totalCount, start, batchSize int) int {
	return max(0, min(totalCount-start, batchSize))
}

type OffsetPagination struct{}

func (o OffsetPagination) StartIndex(batchStart, workerNumber, batchSize int) int {
	return batchStart + workerNumber
}

func (o OffsetPagination) FetchCount(totalCount, start, batchSize int) int {
	return batchSize
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
