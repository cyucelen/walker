package walker

import (
	"context"
	"time"
)

type Option func(*config)

type config struct {
	maxBatchSize  int
	parallelism   int
	pagination    Pagination
	limiter       Limiter
	rateLimit     rateLimit
	context       context.Context
	contextCancel context.CancelFunc
}

func WithMaxBatchSize(size int) Option {
	return func(c *config) {
		c.maxBatchSize = size
	}
}

func WithParallelism(parallelism int) Option {
	return func(c *config) {
		c.parallelism = parallelism
	}
}

func WithLimiter(limiter Limiter) Option {
	return func(c *config) {
		c.limiter = limiter
	}
}

func WithPagination(pagination Pagination) Option {
	return func(c *config) {
		c.pagination = pagination
	}
}

func WithContext(ctx context.Context) Option {
	return func(c *config) {
		c.context, c.contextCancel = context.WithCancel(ctx)
	}
}

func WithRateLimit(count int, per time.Duration) Option {
	return func(c *config) {
		c.rateLimit = rateLimit{count: count, per: per}
	}
}
