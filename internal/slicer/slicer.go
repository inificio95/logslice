// Package slicer provides functionality to extract time-range segments
// from structured log files using binary search over byte offsets.
package slicer

import (
	"io"
	"time"

	"github.com/yourorg/logslice/internal/reader"
	"github.com/yourorg/logslice/internal/timestamp"
)

// Options configures the slicer behavior.
type Options struct {
	// Format is an optional explicit timestamp format string.
	Format string
	// Location is the time zone to use when parsing timestamps.
	Location *time.Location
}

// Slicer extracts log lines within a time range from a ReadSeeker.
type Slicer struct {
	scanner   *reader.Scanner
	extractor *timestamp.Extractor
	opts      Options
}

// New creates a new Slicer for the given source.
func New(src io.ReadSeeker, size int64, opts Options) (*Slicer, error) {
	ext := timestamp.DefaultExtractor()
	scanner, err := reader.NewScanner(src, size, ext)
	if err != nil {
		return nil, err
	}
	return &Slicer{
		scanner:   scanner,
		extractor: ext,
		opts:      opts,
	}, nil
}

// Slice returns all log lines whose timestamps fall within [from, to].
// Lines without a parseable timestamp are skipped.
func (s *Slicer) Slice(from, to time.Time) ([]string, error) {
	start, err := s.scanner.FindFirst(from)
	if err != nil {
		return nil, err
	}

	lines, err := s.scanner.ReadUntil(start, to)
	if err != nil {
		return nil, err
	}
	return lines, nil
}
