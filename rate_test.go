package hypixel

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func makeResp(body, rem, reset string, status int) *http.Response {
	header := http.Header{}
	if rem != "" {
		header.Set("RateLimit-Remaining", rem)
	}
	if reset != "" {
		header.Set("RateLimit-Reset", reset)
	}
	return &http.Response{
		Header:     header,
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestUpdateFromResponse_Throttle(t *testing.T) {
	r := NewRateLimit()
	resp := makeResp(`{"throttle":true}`, "", "", 429)
	if err := r.updateFromResponse(resp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := r.remaining.Load(); got != -1 {
		t.Errorf("got remaining=%d; want -1", got)
	}
}

func TestUpdateFromResponse_ResetHeader(t *testing.T) {
	r := NewRateLimit()
	resp := makeResp(`{"throttle":false}`, "", "2", 200)
	start := time.Now()
	if err := r.updateFromResponse(resp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	reset := r.resetAt.Load().(time.Time)
	want := start.Add(2 * time.Second)
	if reset.Before(start) || reset.After(want.Add(50*time.Millisecond)) {
		t.Errorf("resetAt=%v; want between %v and %v", reset, start, want.Add(50*time.Millisecond))
	}
}

func TestUpdateFromResponse_RemainingHeader(t *testing.T) {
	r := NewRateLimit()
	resp := makeResp(`{"throttle":false}`, "5", "", 200)
	if err := r.updateFromResponse(resp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := r.remaining.Load(); got != 5 {
		t.Errorf("got remaining=%d; want 5", got)
	}
}

func TestReset(t *testing.T) {
	r := NewRateLimit()
	r.remaining.Store(10)
	r.resetAt.Store(time.Now().Add(10 * time.Minute))
	r.Reset()
	if got := r.remaining.Load(); got != -1 {
		t.Errorf("after Reset, remaining=%d; want -1", got)
	}
	if reset := r.resetAt.Load().(time.Time); !reset.IsZero() {
		t.Errorf("after Reset, resetAt=%v; want zero", reset)
	}
}

func TestString(t *testing.T) {
	r := NewRateLimit()
	r.remaining.Store(3)
	tm := time.Date(2025, 5, 18, 12, 0, 0, 0, time.UTC)
	r.resetAt.Store(tm)
	s := r.String()
	if !strings.Contains(s, "3 remaining") || !strings.Contains(s, tm.Format(time.RFC3339)) {
		t.Errorf("String() = %q; want to contain remaining and timestamp", s)
	}
}

func TestWaitIfNeeded_NoWaitOnPositive(t *testing.T) {
	r := NewRateLimit()
	r.remaining.Store(1)
	before := time.Now()
	r.waitIfNeeded()
	if time.Since(before) > 50*time.Millisecond {
		t.Errorf("waitIfNeeded slept unexpectedly")
	}
}

func TestWaitIfNeeded_NoWaitOnPastReset(t *testing.T) {
	r := NewRateLimit()
	r.remaining.Store(0)
	r.resetAt.Store(time.Now().Add(-time.Minute))
	before := time.Now()
	r.waitIfNeeded()
	if time.Since(before) > 50*time.Millisecond {
		t.Errorf("waitIfNeeded slept for past reset")
	}
}
