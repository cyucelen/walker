package walker

import "math"

func InfiniteLimiter() Limiter {
	return func() int {
		return math.MaxInt
	}
}

func ConstantLimiter(limit int) Limiter {
	return func() int {
		return limit
	}
}
