package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/output"
)

func TestWriterSingleLine(t *testing.T) {
	var buf bytes.Buffer
	w := output.NewWriter(&buf)

	if err := w.WriteLine([]byte("hello world")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("flush error: %v", err)
	}

	got := buf.String()
	if got != "hello world\n" {
		t.Errorf("expected %q, got %q", "hello world\n", got)
	}
}

func TestWriterMultipleLines(t *testing.T) {
	lines := []string{
		"2024-01-01T00:00:00Z INFO starting",
		"2024-01-01T00:00:01Z DEBUG tick",
		"2024-01-01T00:00:02Z INFO stopping",
	}

	var buf bytes.Buffer
	w := output.NewWriter(&buf)

	for _, l := range lines {
		if err := w.WriteLine([]byte(l)); err != nil {
			t.Fatalf("write error: %v", err)
		}
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("flush error: %v", err)
	}

	if w.LinesWritten() != len(lines) {
		t.Errorf("expected %d lines written, got %d", len(lines), w.LinesWritten())
	}

	got := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(got) != len(lines) {
		t.Fatalf("expected %d output lines, got %d", len(lines), len(got))
	}
	for i, want := range lines {
		if got[i] != want {
			t.Errorf("line %d: expected %q, got %q", i, want, got[i])
		}
	}
}

func TestWriterEmptyLine(t *testing.T) {
	var buf bytes.Buffer
	w := output.NewWriter(&buf)

	if err := w.WriteLine([]byte{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = w.Flush()

	if buf.String() != "\n" {
		t.Errorf("expected bare newline for empty line, got %q", buf.String())
	}
	if w.LinesWritten() != 1 {
		t.Errorf("expected 1 line written, got %d", w.LinesWritten())
	}
}
