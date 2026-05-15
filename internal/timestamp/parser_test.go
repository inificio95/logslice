package timestamp

import (
	"testing"
	"time"
)

func TestParseRFC3339(t *testing.T) {
	p := NewParser("", time.UTC)
	got, err := p.Parse("2024-03-15T10:22:33Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := time.Date(2024, 3, 15, 10, 22, 33, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseSpaceSeparated(t *testing.T) {
	p := NewParser("", time.UTC)
	got, err := p.Parse("2024-03-15 10:22:33")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := time.Date(2024, 3, 15, 10, 22, 33, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseExplicitFormat(t *testing.T) {
	p := NewParser("2006/01/02 15:04:05", time.UTC)
	got, err := p.Parse("2024/03/15 10:22:33")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := time.Date(2024, 3, 15, 10, 22, 33, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseUnknownFormat(t *testing.T) {
	p := NewParser("", time.UTC)
	_, err := p.Parse("not-a-timestamp")
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}

func TestParseWithLocation(t *testing.T) {
	loc, _ := time.LoadLocation("America/New_York")
	p := NewParser("", loc)
	got, err := p.Parse("2024-03-15 10:22:33")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Location().String() != "America/New_York" {
		t.Errorf("expected America/New_York, got %v", got.Location())
	}
}

func TestAutoDetectNanoseconds(t *testing.T) {
	p := NewParser("", time.UTC)
	got, err := p.Parse("2024-03-15T10:22:33.123456789Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Nanosecond() != 123456789 {
		t.Errorf("expected nanoseconds 123456789, got %d", got.Nanosecond())
	}
}
