/*
Package ratelimit helps you limit the transfer rate using Token-Bucket algorithm.
It is rewrote from and inspired by http://github.com/juju/ratelimit
*/
package ratelimit

import (
	"sync"
	"time"
)

// Bucket is a thread-safe rate limiter.
// It uses Token-Bucket algorithm to limit the transfer rate.
type Bucket struct {
	lastTime     time.Time
	capacity     int64
	fillInterval time.Duration
	avail        int64
	lock         sync.Mutex
	transferUnit int64
}

func (b *Bucket) fill() int64 {
	now := time.Now()
	dur := now.Sub(b.lastTime)
	tokens := int64(dur) / int64(b.fillInterval)
	b.lastTime = now
	if tokens < b.transferUnit {
		time.Sleep(time.Duration(b.transferUnit - tokens) * b.fillInterval)
		tokens = b.transferUnit
		b.lastTime = time.Now()
	}

	b.avail += tokens
	if b.avail > b.capacity {
		b.avail = b.capacity
	}

	return b.avail
}

// Take will accquire at most n tokens from bucket.
// It returns the number of tokens accquired, not more than n or capacity.
//
// Take will block until (at least) a number of tokens (transferUnit) available,
// even if n < transferUnit.
func (b *Bucket) Take(n int64) int64 {
	b.lock.Lock()
	defer b.lock.Unlock()

	if n > b.capacity {
		n = b.capacity
	}

	if n > b.avail {
		if n > b.fill() {
			n = b.avail
		}
	}

	b.avail -= n
	return n
}

// Return releases n unused tokens.
func (b *Bucket) Return(n int64) {
	if n <= 0 {
		return
	}
	b.lock.Lock()
	defer b.lock.Unlock()

	b.avail += n
	if b.avail > b.capacity {
		b.avail = b.capacity
	}
}

// Capacity returns capacity of this bucket.
func (b *Bucket) Capacity() int64 {
	return b.capacity
}

// New creates a Bucket by specifying intervals to fill a token.
func New(fillInterval time.Duration, capacity int64, transferUnit int64) *Bucket {
	if capacity < 2 {
		capacity = 2
	}
	if transferUnit <= 0 || transferUnit > capacity/2 {
		transferUnit = capacity/2
		if transferUnit<1 {
			transferUnit = 1
		}
	}
	return &Bucket{
		lastTime:     time.Now(),
		capacity:     capacity,
		fillInterval: fillInterval,
		avail:        0,
		lock:         sync.Mutex{},
		transferUnit: transferUnit,
	}
}

// NewFromRate creates a Bucket by specifying transfer rate in bytes per second.
func NewFromRate(rate float64, capacity int64, transferUnit int64) *Bucket {
	return New(time.Second/time.Duration(rate), capacity, capacity)
}
