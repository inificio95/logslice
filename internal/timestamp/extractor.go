package timestamp

import (
	"fmt"
	"regexp"
	"time"
)

// Extractor finds and parses a timestamp within a raw log line.
type Extractor struct {
	parser  *Parser
	pattern *regexp.Regexp
}

// NewExtractor creates an Extractor using a regex pattern to locate the
// timestamp field within a line. The pattern must contain exactly one
// capturing group that isolates the timestamp string.
func NewExtractor(pattern, format string, loc *time.Location) (*Extractor, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern: %w", err)
	}
	if re.NumSubexp() < 1 {
		return nil, fmt.Errorf("pattern must contain at least one capturing group")
	}
	return &Extractor{
		parser:  NewParser(format, loc),
		pattern: re,
	}, nil
}

// Extract pulls a timestamp out of line using the compiled regex, then parses it.
func (e *Extractor) Extract(line string) (time.Time, error) {
	matches := e.pattern.FindStringSubmatch(line)
	if matches == nil {
		return time.Time{}, fmt.Errorf("no timestamp found in line")
	}
	return e.parser.Parse(matches[1])
}

// DefaultExtractor returns an Extractor that matches common ISO-8601 timestamps
// at the start of a log line (optionally wrapped in brackets).
func DefaultExtractor(loc *time.Location) (*Extractor, error) {
	// Matches: [2024-03-15T10:22:33Z] or 2024-03-15T10:22:33.000Z or similar
	const defaultPattern = `^\[?(\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2}[^\]\s]*)\]?`
	return NewExtractor(defaultPattern, "", loc)
}
