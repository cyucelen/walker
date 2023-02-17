package walker

import "time"

type rateLimit struct {
	count     int
	per       time.Duration
	unlimited bool
}

var defaultRateLimiter = rateLimit{unlimited: true}
