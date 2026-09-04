package hypixel

import (
	"math"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type RateLimit struct {
	remaining atomic.Int32
	limit     atomic.Int32
	resetAt   atomic.Value
	mu        sync.Mutex
	waitCh    chan struct{}
}

func NewRateLimit() *RateLimit {
	r := &RateLimit{}
	r.remaining.Store(-1)
	r.limit.Store(-1)
	r.resetAt.Store(time.Time{})
	return r
}

func (r *RateLimit) Reset() {
	r.remaining.Store(-1)
	r.limit.Store(-1)
	r.resetAt.Store(time.Time{})
}

func (r *RateLimit) GetRemaining() int32 {
	return r.remaining.Load()
}

func (r *RateLimit) GetLimit() int32 {
	return r.limit.Load()
}

func (r *RateLimit) GetResetAt() time.Time {
	if val := r.resetAt.Load(); val != nil {
		return val.(time.Time)
	}
	return time.Time{}
}

func (r *RateLimit) String() string {
	return strconv.Itoa(int(r.remaining.Load())) + " remaining until " + r.GetResetAt().Format(time.RFC3339)
}

func (r *RateLimit) waitIfNeeded() {
	for {
		r.mu.Lock()
		reset := r.GetResetAt()

		if r.remaining.Load() > 0 || reset.IsZero() || time.Now().After(reset) {
			r.mu.Unlock()
			return
		}

		if r.waitCh == nil {
			ch := make(chan struct{})
			r.waitCh = ch
			go func(ch chan struct{}, reset time.Time) {
				sleep := time.Until(reset)
				if sleep > 5*time.Minute {
					sleep = 5 * time.Minute
				}
				time.Sleep(sleep)
				r.mu.Lock()
				close(ch)
				r.waitCh = nil
				r.mu.Unlock()
			}(ch, reset)
		}

		ch := r.waitCh
		r.mu.Unlock()

		<-ch
	}
}

func (r *RateLimit) updateFromResponse(resp *http.Response) error {
	if limitStr := resp.Header.Get("RateLimit-Limit"); limitStr != "" {
		lim, err := strconv.Atoi(limitStr)
		if err != nil {
			return err
		}
		if lim > 0 && lim <= math.MaxInt32 {
			r.limit.Store(int32(lim))
		}
	}

	if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
		if secs, err := strconv.Atoi(retryAfter); err == nil {
			r.resetAt.Store(time.Now().Add(time.Duration(secs) * time.Second))
		} else if t, err := http.ParseTime(retryAfter); err == nil {
			r.resetAt.Store(t)
		}
	} else if resetStr := resp.Header.Get("RateLimit-Reset"); resetStr != "" {
		secs, err := strconv.Atoi(resetStr)
		if err != nil {
			return err
		}
		r.resetAt.Store(time.Now().Add(time.Duration(secs) * time.Second))
	} else if resp.StatusCode == http.StatusTooManyRequests {
		r.resetAt.Store(time.Now().Add(time.Second))
	}

	remStr := resp.Header.Get("RateLimit-Remaining")
	if remStr == "" {
		return nil
	}
	rem, err := strconv.Atoi(remStr)
	if err != nil {
		return err
	}
	if rem >= 0 && rem <= math.MaxInt32 {
		r.remaining.Store(int32(rem))
	} else {
		r.remaining.Store(-1)
	}
	return nil
}
