package util

import "testing"

func TestNewRequestID(t *testing.T) {
	a, b := NewRequestID(), NewRequestID()
	if a == "" || b == "" || a == b {
		t.Fatalf("request ids should be non-empty and unique, got %q %q", a, b)
	}
}
