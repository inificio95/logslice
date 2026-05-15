package reader

import (
	"bufio"
	"io"
	"time"

	"github.com/logslice/logslice/internal/timestamp"
)

// Line represents a single log line with its parsed timestamp.
type Line struct {
	Raw       string
	Timestamp time.Time
	Offset    int64
}

// Scanner reads log lines sequentially and extracts timestamps.
type Scanner struct {
	reader    io.ReadSeeker
	bufReader *bufio.Reader
	extractor *timestamp.Extractor
	offset    int64
	current   Line
	err       error
}

// NewScanner creates a Scanner wrapping the given ReadSeeker.
func NewScanner(r io.ReadSeeker, ext *timestamp.Extractor) *Scanner {
	return &Scanner{
		reader:    r,
		bufReader: bufio.NewReaderSize(r, 64*1024),
		extractor: ext,
	}
}

// Scan advances to the next line. Returns true while lines are available.
func (s *Scanner) Scan() bool {
	lineStart := s.offset
	raw, err := s.bufReader.ReadString('\n')
	if len(raw) == 0 {
		s.err = err
		return false
	}
	s.offset += int64(len(raw))

	// Trim trailing newline for cleaner processing.
	trimmed := raw
	if len(trimmed) > 0 && trimmed[len(trimmed)-1] == '\n' {
		trimmed = trimmed[:len(trimmed)-1]
	}

	ts, _ := s.extractor.Extract(trimmed)
	s.current = Line{
		Raw:       trimmed,
		Timestamp: ts,
		Offset:    lineStart,
	}
	return true
}

// Line returns the current line.
func (s *Scanner) Line() Line {
	return s.current
}

// Err returns any non-EOF error encountered.
func (s *Scanner) Err() error {
	if s.err == io.EOF {
		return nil
	}
	return s.err
}

// SeekOffset moves the underlying reader to the given byte offset.
func (s *Scanner) SeekOffset(offset int64) error {
	_, err := s.reader.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}
	s.bufReader.Reset(s.reader)
	s.offset = offset
	return nil
}
