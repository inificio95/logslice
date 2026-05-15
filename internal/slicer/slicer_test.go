package slicer_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/slicer"
)

const sampleLog = `2024-01-10T08:00:00Z INFO  service started
2024-01-10T08:01:00Z DEBUG request received id=1
2024-01-10T08:02:00Z INFO  processed id=1
2024-01-10T08:03:00Z WARN  slow query latency=500ms
2024-01-10T08:04:00Z ERROR disk usage high pct=92
2024-01-10T08:05:00Z INFO  health check ok
`

func newSlicer(t *testing.T, data string) *slicer.Slicer {
	t.Helper()
	rs := strings.NewReader(data)
	s, err := slicer.New(rs, int64(len(data)), slicer.Options{})
	if err != nil {
		t.Fatalf("slicer.New: %v", err)
	}
	return s
}

func mustParse(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestSliceExactRange(t *testing.T) {
	s := newSlicer(t, sampleLog)
	lines, err := s.Slice(mustParse("2024-01-10T08:01:00Z"), mustParse("2024-01-10T08:03:00Z"))
	if err != nil {
		t.Fatalf("Slice: %v", err)
	}
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d: %v", len(lines), lines)
	}
}

func TestSliceFullRange(t *testing.T) {
	s := newSlicer(t, sampleLog)
	lines, err := s.Slice(mustParse("2024-01-10T08:00:00Z"), mustParse("2024-01-10T08:05:00Z"))
	if err != nil {
		t.Fatalf("Slice: %v", err)
	}
	if len(lines) != 6 {
		t.Fatalf("expected 6 lines, got %d", len(lines))
	}
}

func TestSliceNoMatch(t *testing.T) {
	s := newSlicer(t, sampleLog)
	lines, err := s.Slice(mustParse("2024-01-10T09:00:00Z"), mustParse("2024-01-10T10:00:00Z"))
	if err != nil {
		t.Fatalf("Slice: %v", err)
	}
	if len(lines) != 0 {
		t.Fatalf("expected 0 lines, got %d", len(lines))
	}
}

func TestSliceEmptyInput(t *testing.T) {
	s := newSlicer(t, "")
	lines, err := s.Slice(mustParse("2024-01-10T08:00:00Z"), mustParse("2024-01-10T09:00:00Z"))
	if err != nil {
		t.Fatalf("Slice: %v", err)
	}
	if len(lines) != 0 {
		t.Fatalf("expected 0 lines, got %d", len(lines))
	}
}
