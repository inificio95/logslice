package timestamp

import (
	"fmt"
	"time"
)

// Common log timestamp formats to try in order
var knownFormats = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05.999999999",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05.999999999",
	"2006-01-02 15:04:05",
	"02/Jan/2006:15:04:05 -0700",
	"Jan 02 15:04:05",
}

// Parser extracts and parses timestamps from log lines.
type Parser struct {
	format string
	loc    *time.Location
}

// NewParser creates a Parser. If format is empty, auto-detection is used.
func NewParser(format string, loc *time.Location) *Parser {
	if loc == nil {
		loc = time.UTC
	}
	return &Parser{format: format, loc: loc}
}

// Parse attempts to parse a timestamp string into a time.Time.
func (p *Parser) Parse(value string) (time.Time, error) {
	if p.format != "" {
		t, err := time.ParseInLocation(p.format, value, p.loc)
		if err != nil {
			return time.Time{}, fmt.Errorf("parse with format %q: %w", p.format, err)
		}
		return t, nil
	}
	return autoDetect(value, p.loc)
}

// autoDetect tries each known format until one succeeds.
func autoDetect(value string, loc *time.Location) (time.Time, error) {
	for _, f := range knownFormats {
		if t, err := time.ParseInLocation(f, value, loc); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("timestamp %q did not match any known format", value)
}
