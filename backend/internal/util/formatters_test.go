package util

import (
	"testing"
	"time"
)

func TestFormatDate(t *testing.T) {
	got := FormatDate(time.Date(2026, 8, 16, 12, 0, 0, 0, time.UTC))
	if got != "2026-08-16" {
		t.Fatalf("FormatDate=%q want 2026-08-16", got)
	}
}

func TestFormatDateTime(t *testing.T) {
	got := FormatDateTime(time.Date(2026, 8, 16, 12, 0, 0, 0, time.FixedZone("CST", 8*3600)))
	if got == "" {
		t.Fatal("FormatDateTime should not be empty")
	}
}
