package reader

import (
	"io"
	"time"
)

// FindFirst performs a binary search over the indexed line offsets and returns
// the byte offset of the first line whose timestamp is >= target.
// Returns -1 if no such line exists.
func (s *Scanner) FindFirst(target time.Time) (int64, error) {
	offsets := s.Offsets()
	if len(offsets) == 0 {
		return -1, nil
	}

	lo, hi := 0, len(offsets)-1
	result := -1

	for lo <= hi {
		mid := (lo + hi) / 2
		t, err := s.TimestampAt(offsets[mid])
		if err != nil {
			// Skip unparseable lines by moving forward.
			lo = mid + 1
			continue
		}
		if !t.Before(target) {
			result = mid
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}

	if result == -1 {
		return -1, nil
	}
	return offsets[result], nil
}

// ReadUntil reads lines sequentially starting at byteOffset, collecting lines
// whose timestamps are <= to. Returns an empty slice if offset is -1.
func (s *Scanner) ReadUntil(offset int64, to time.Time) ([]string, error) {
	if offset < 0 {
		return nil, nil
	}
	if _, err := s.Seek(offset, io.SeekStart); err != nil {
		return nil, err
	}

	var lines []string
	for {
		line, ts, err := s.NextLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if ts.IsZero() {
			continue
		}
		if ts.After(to) {
			break
		}
		lines = append(lines, line)
	}
	return lines, nil
}
