package rest

import (
	"time"
)

// RateLimiterConfig is the configuration for the rate limiter.
type RateLimiterConfig struct {
	// MaxRetries is the maximum number of retries that can be performed
	// before giving up.
	MaxRetries int
	// CleanupInterval is the interval at which the rate limiter will
	// clean up old buckets.
	CleanupInterval time.Duration
}

// DefaultRateLimiterConfig returns a default rate limiter configuration.
func DefaultRateLimiterConfig() *RateLimiterConfig {
	return &RateLimiterConfig{
		MaxRetries:      5,
		CleanupInterval: 5 * time.Second,
	}
}

// Apply applies the given options to the configuration.
func (c *RateLimiterConfig) Apply(opts []RateLimiterOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

// RateLimiterOpt is an option for the rate limiter.
type RateLimiterOpt func(config *RateLimiterConfig)

// WithMaxRetries sets the maximum number of retries that can be performed
// before giving up.
func WithMaxRetries(maxRetries int) RateLimiterOpt {
	return func(config *RateLimiterConfig) {
		config.MaxRetries = maxRetries
	}
}

// WithCleanupInterval sets the interval at which the rate limiter will
// clean up old buckets.
func WithCleanupInterval(interval time.Duration) RateLimiterOpt {
	return func(config *RateLimiterConfig) {
		config.CleanupInterval = interval
	}
}
