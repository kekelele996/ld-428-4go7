package util

import (
	"strings"
	"testing"
)

func TestNewID(t *testing.T) {
	a := NewID("art")
	b := NewID("art")
	if !strings.HasPrefix(a, "art-") || len(a) < 5 {
		t.Fatalf("unexpected id: %s", a)
	}
	if a == b {
		t.Fatalf("ids should differ: %s", a)
	}
}
