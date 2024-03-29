// SPDX-License-Identifier: MIT

package client

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/apex/log"
)

// bucket is one scope of rate limiting.
// It is what is used to determine if a request should be sent or not, and what limit one falls under.
type bucket struct {
	Endpoint  *CompiledEndpoint
	Reset     time.Time
	Mutex     sync.Mutex
	Remaining int
	Limit     int
	Reserved  int
	First     bool
}

// RateLimiter is an interface for handling rate limiting.
type RateLimiter interface {
	// MaxRetries is how many times should a request be retried
	// if it fails due to rate limiting?
	MaxRetries() int
	// Reset the rate limiter.
	Reset()
	// WaitBucket waits for the bucket to be available.
	// Returns true if this is the first time the bucket is used,
	// and so the lock must be held over the request.
	WaitBucket(*CompiledEndpoint) (bool, error)
	// UnlockBucket unlocks the bucket after a request has been sent.
	// Takes in the boolean returned by WaitBucket.
	UnlockBucket(*CompiledEndpoint, *http.Response, bool) error
}

// NewRateLimiter creates a new default rate limiter.
func NewRateLimiter(opts ...RateLimiterOpt) RateLimiter {
	config := DefaultRateLimiterConfig()
	config.Apply(opts)

	rateLimiter := &rateLimiterImpl{
		config:  *config,
		hashes:  map[*Endpoint]string{},
		buckets: map[string]*bucket{},
	}

	go rateLimiter.cleanup()

	return rateLimiter
}

// cleanup continuously clean up old buckets by calling doCleanup.
func (l *rateLimiterImpl) cleanup() {
	ticker := time.NewTicker(l.config.CleanupInterval)
	for range ticker.C {
		l.doCleanup()
	}
}

// doCleanup cleans up old buckets.
func (r *rateLimiterImpl) doCleanup() {
	r.bucketsMutex.Lock()
	defer r.bucketsMutex.Unlock()

	now := time.Now()
	for hash, b := range r.buckets {
		if b.Reset.Before(now) && b.Reserved == 0 {
			delete(r.buckets, hash)
		}
	}
}

// rateLimiterImpl is the default implementation of RateLimiter.
type rateLimiterImpl struct {
	hashes       map[*Endpoint]string
	hashesMutex  sync.Mutex
	buckets      map[string]*bucket
	bucketsMutex sync.Mutex
	config       RateLimiterConfig
}

func (r *rateLimiterImpl) MaxRetries() int {
	return r.config.MaxRetries
}

func (r *rateLimiterImpl) Reset() {
	r.buckets = map[string]*bucket{}
	r.bucketsMutex = sync.Mutex{}
	r.hashes = map[*Endpoint]string{}
	r.hashesMutex = sync.Mutex{}
}

// getRouteHash gets the hash for the current endpoint
func (r *rateLimiterImpl) getRouteHash(endpoint *CompiledEndpoint) string {
	r.hashesMutex.Lock()
	defer r.hashesMutex.Unlock()

	hash, ok := r.hashes[endpoint.Endpoint]
	if !ok {
		// generate routeHash
		hash = endpoint.Endpoint.Method + "+" + endpoint.Endpoint.Route
		r.hashes[endpoint.Endpoint] = hash
	}
	return hash
}

// getBucket gets the bucket for the current endpoint.
// If create is true, the bucket will be created if it does not exist.
// Else, nil will be returned if the bucket could not be found.
func (r *rateLimiterImpl) getBucket(endpoint *CompiledEndpoint, create bool) *bucket {
	hash := r.getRouteHash(endpoint)

	r.bucketsMutex.Lock()
	defer r.bucketsMutex.Unlock()

	b, ok := r.buckets[hash]
	if !ok {
		if !create {
			return nil
		}

		b = &bucket{
			Endpoint:  endpoint,
			Remaining: 1,
			// Limit not known yet.
			Limit: -1,
			// If we can check the rate limits "safely"
			First: true,
		}
		r.buckets[hash] = b
	}
	return b
}

func (r *rateLimiterImpl) WaitBucket(endpoint *CompiledEndpoint) (bool, error) {
	return r.doWaitBucket(endpoint, r.config.MaxRetries, false)
}

// doWaitBucket is the handling for waiting for a bucket.
// This is recursive if the bucket is not ready.
func (r *rateLimiterImpl) doWaitBucket(endpoint *CompiledEndpoint, retries int, first bool) (bool, error) {
	if retries < 0 {
		return false, fmt.Errorf("retries exhausted")
	}

	bucket := r.getBucket(endpoint, true)

	bucket.Mutex.Lock()
	first = bucket.First || first
	bucket.First = false

	now := time.Now()
	var until time.Time
	if (bucket.Remaining - bucket.Reserved) > 0 {
		until = now
	} else if first {
		until = now
	} else {
		until = bucket.Reset
	}

	if until.After(now) {
		// If we are the first request to be limited, we can be the first one to try again
		// and update the rate limit information.
		first := bucket.Reserved == 0
		bucket.Mutex.Unlock()

		log.WithField("bucket", bucket.Endpoint.Path).WithField("wait", until.Sub(now).Seconds()).Debug("Ratelimited")

		<-time.After(until.Sub(now))
		return r.doWaitBucket(endpoint, retries-1, first)
	}

	bucket.Reserved++
	if !first {
		bucket.Mutex.Unlock()
	}

	return first, nil
}

func (l *rateLimiterImpl) UnlockBucket(endpoint *CompiledEndpoint, rs *http.Response, locked bool) error {
	b := l.getBucket(endpoint, false)
	if b == nil {
		return nil
	}

	// no response provided means we can't update anything and just unlock it
	if rs == nil || rs.Header == nil {
		return nil
	}

	if !locked {
		b.Mutex.Lock()
	}

	defer b.Mutex.Unlock()
	b.Reserved--

	resetHeader := rs.Header.Get("X-RateLimit-Reset")
	maxHeader := rs.Header.Get("X-RateLimit-Max")
	lastResetHeader := rs.Header.Get("X-RateLimit-Last-Reset")
	requestCountHeader := rs.Header.Get("X-RateLimit-Request-Count")

	if rs.StatusCode == http.StatusTooManyRequests {
		interval, err := strconv.Atoi(resetHeader)
		if err != nil {
			return fmt.Errorf("invalid retryAfter %s: %s", resetHeader, err)
		}
		lastReset, err := strconv.Atoi(lastResetHeader)
		if err != nil {
			return fmt.Errorf("invalid retryAfter %s: %s", lastResetHeader, err)
		}

		// lastReset is when the bucket resets, interval is time between resets.
		reset_sec := int64(lastReset / 1000)
		reset_nsec := int64(lastReset%1000) * 1000000
		reset := time.Unix(reset_sec, reset_nsec).Add(time.Duration(interval) * time.Millisecond)
		b.Remaining = 0
		b.Reset = reset

		return nil
	}

	if maxHeader != "" {
		limit, err := strconv.Atoi(maxHeader)
		if err != nil {
			return fmt.Errorf("invalid max %s: %s", maxHeader, err)
		}
		b.Limit = limit
	}

	if requestCountHeader != "" {
		requested, err := strconv.Atoi(requestCountHeader)
		if err != nil {
			return fmt.Errorf("invalid request count %s: %s", requestCountHeader, err)
		}
		b.Remaining = b.Limit - requested
	}

	// we prioritize the reset after header over the reset header as it's more accurate due to clock differences
	if resetHeader != "" && lastResetHeader != "" {
		interval, err := strconv.Atoi(resetHeader)
		if err != nil {
			return fmt.Errorf("invalid reset %s: %s", resetHeader, err)
		}
		lastReset, err := strconv.Atoi(lastResetHeader)
		if err != nil {
			return fmt.Errorf("invalid last reset %s: %s", lastResetHeader, err)
		}

		reset_sec := int64(lastReset / 1000)
		reset_nsec := int64(lastReset%1000) * 1000000
		b.Reset = time.Unix(reset_sec, reset_nsec).Add(time.Duration(interval) * time.Millisecond)
	} else {
		return fmt.Errorf("no reset header and last reset header found in response")
	}

	return nil
}
