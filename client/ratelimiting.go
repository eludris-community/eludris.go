package client

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type bucket struct {
	Endpoint  *CompiledEndpoint
	Reset     time.Time
	Mutex     sync.Mutex
	Remaining int
	Limit     int
	Reserved  int
	First     bool
	Received  chan struct{}
}

type RateLimiter interface {
	MaxRetries() int
	Reset()
	WaitBucket(*CompiledEndpoint, int) (bool, error)
	UnlockBucket(*CompiledEndpoint, *http.Response, bool) error
}

func NewRateLimiter() RateLimiter {
	rateLimiter := &rateLimiterImpl{
		hashes:  map[*Endpoint]string{},
		buckets: map[string]*bucket{},
	}

	go rateLimiter.cleanup()

	return rateLimiter
}

func (l *rateLimiterImpl) cleanup() {
	// TODO: Make interval configurable.
	ticker := time.NewTicker(time.Second * 10)
	for range ticker.C {
		l.doCleanup()
	}
}

func (r *rateLimiterImpl) doCleanup() {
	r.bucketsMutex.Lock()
	defer r.bucketsMutex.Unlock()

	now := time.Now().UTC()
	for hash, b := range r.buckets {
		if b.Reset.Before(now) && b.Reserved == 0 {
			delete(r.buckets, hash)
		}
	}
}

type rateLimiterImpl struct {
	hashes       map[*Endpoint]string
	hashesMutex  sync.Mutex
	buckets      map[string]*bucket
	bucketsMutex sync.Mutex
}

func (r *rateLimiterImpl) MaxRetries() int {
	return 10 // TODO: Make this configurable.
}

func (r *rateLimiterImpl) Reset() {
	r.buckets = map[string]*bucket{}
	r.bucketsMutex = sync.Mutex{}
	r.hashes = map[*Endpoint]string{}
	r.hashesMutex = sync.Mutex{}
}

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

func (r *rateLimiterImpl) WaitBucket(endpoint *CompiledEndpoint, retries int) (bool, error) {
	return r.doWaitBucket(endpoint, retries, false)
}

func (r *rateLimiterImpl) doWaitBucket(endpoint *CompiledEndpoint, retries int, first bool) (bool, error) {
	if retries < 0 {
		return false, fmt.Errorf("retries exhausted")
	}

	bucket := r.getBucket(endpoint, true)

	bucket.Mutex.Lock()
	first = bucket.First || first
	bucket.First = false

	now := time.Now().UTC()
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

		log.Printf("Ratelimited on bucket %s, waiting %.1fs", bucket.Endpoint.Path, until.Sub(now).Seconds())

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
		reset := time.Unix(reset_sec, reset_nsec).Add(time.Duration(interval) * time.Millisecond).UTC()
		b.Remaining = 0
		b.Reset = reset

		select {
		case <-b.Received:
		default:
		}

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
		b.Reset = time.Unix(reset_sec, reset_nsec).Add(time.Duration(interval) * time.Millisecond).UTC()
	} else {
		return fmt.Errorf("no reset header and last reset header found in response")
	}

	select {
	case <-b.Received:
	default:
	}

	return nil
}
