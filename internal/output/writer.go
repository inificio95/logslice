// Package output handles writing sliced log segments to various destinations.
package output

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Writer writes log lines to an output destination.
type Writer struct {
	w   *bufio.Writer
	count int
}

// NewWriter creates a Writer that writes to the given io.Writer.
func NewWriter(w io.Writer) *Writer {
	return &Writer{w: bufio.NewWriter(w)}
}

// NewFileWriter opens (or creates) a file at path and returns a Writer for it.
// The caller is responsible for closing the underlying file.
func NewFileWriter(path string) (*Writer, *os.File, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, nil, fmt.Errorf("output: create file %q: %w", path, err)
	}
	return NewWriter(f), f, nil
}

// WriteLine writes a single log line followed by a newline character.
func (w *Writer) WriteLine(line []byte) error {
	if _, err := w.w.Write(line); err != nil {
		return fmt.Errorf("output: write line: %w", err)
	}
	if err := w.w.WriteByte('\n'); err != nil {
		return fmt.Errorf("output: write newline: %w", err)
	}
	w.count++
	return nil
}

// Flush flushes any buffered data to the underlying writer.
func (w *Writer) Flush() error {
	if err := w.w.Flush(); err != nil {
		return fmt.Errorf("output: flush: %w", err)
	}
	return nil
}

// LinesWritten returns the total number of lines written so far.
func (w *Writer) LinesWritten() int {
	return w.count
}
