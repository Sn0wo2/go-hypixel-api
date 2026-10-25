package main

import "testing"

func TestGreeting(t *testing.T) {
	t.Parallel()

	if got := greeting(); got != "Hello, world!" {
		t.Errorf("greeting() = %q, want %q", got, "Hello, world!")
	}
}
