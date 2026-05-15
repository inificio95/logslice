package reader

import (
	"strings"
	"testing"
	"time"

	"github.com/logslice/logslice/internal/timestamp"
)

func defaultExtractor(t *testing.T) *timestamp.Extractor {
	t.Helper()
	ext, err := timestamp.NewExtractor("", nil)
	if err != nil {
		t.Fatalf("failed to create extractor: %v", err)
	}
	return ext
}

func TestScannerBasic(t *testing.T) {
	log := "2024-01-15T10:00:00Z level=info msg=start\n" +
		"2024-01-15T10:01:00Z level=info msg=middle\n" +
		"2024-01-15T10:02:00Z level=info msg=end\n"

	s := NewScanner(strings.NewReader(log), defaultExtractor(t))

	var lines []Line
	for s.Scan() {
		lines = append(lines, s.Line())
	}
	if err := s.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0].Timestamp.IsZero() {
		t.Error("expected non-zero timestamp on first line")
	}
}

func TestScannerOffsets(t *testing.T) {
	log := "2024-01-15T10:00:00Z first\n2024-01-15T10:01:00Z second\n"

	s := NewScanner(strings.NewReader(log), defaultExtractor(t))

	s.Scan()
	firstLine := s.Line()
	if firstLine.Offset != 0 {
		t.Errorf("expected offset 0, got %d", firstLine.Offset)
	}

	s.Scan()
	secondLine := s.Line()
	if secondLine.Offset != int64(len("2024-01-15T10:00:00Z first\n")) {
		t.Errorf("unexpected offset for second line: %d", secondLine.Offset)
	}
}

func TestScannerTimestampParsed(t *testing.T) {
	log := "2024-03-01T08:30:00Z event=login user=alice\n"

	s := NewScanner(strings.NewReader(log), defaultExtractor(t))
	s.Scan()
	line := s.Line()

	expected := time.Date(2024, 3, 1, 8, 30, 0, 0, time.UTC)
	if !line.Timestamp.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, line.Timestamp)
	}
}

func TestScannerEmptyInput(t *testing.T) {
	s := NewScanner(strings.NewReader(""), defaultExtractor(t))
	if s.Scan() {
		t.Error("expected Scan to return false on empty input")
	}
	if s.Err() != nil {
		t.Errorf("unexpected error: %v", s.Err())
	}
}
