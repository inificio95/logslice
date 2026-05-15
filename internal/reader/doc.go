// Package reader provides sequential and seekable log line reading
// with integrated timestamp extraction.
//
// The Scanner type wraps an io.ReadSeeker and yields log lines one at a
// time, attaching the parsed timestamp and byte offset to each line.
// Byte offsets enable downstream components (e.g., binary search) to
// seek directly to a position in the file without re-reading from the
// beginning, keeping memory usage constant regardless of file size.
//
// Typical usage:
//
//	ext := timestamp.DefaultExtractor
//	s := reader.NewScanner(file, ext)
//	for s.Scan() {
//		line := s.Line()
//		// process line.Raw, line.Timestamp, line.Offset
//	}
//	if err := s.Err(); err != nil {
//		log.Fatal(err)
//	}
package reader
