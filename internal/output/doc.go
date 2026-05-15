// Package output provides buffered writers for emitting sliced log segments
// to files or arbitrary io.Writer destinations.
//
// # Usage
//
// Use [NewWriter] to wrap any io.Writer:
//
//	w := output.NewWriter(os.Stdout)
//	w.WriteLine([]byte("2024-01-01T00:00:00Z INFO hello"))
//	w.Flush()
//
// Use [NewFileWriter] to write directly to a file path:
//
//	w, f, err := output.NewFileWriter("/tmp/slice.log")
//	if err != nil { ... }
//	defer f.Close()
//	// ... write lines ...
//	w.Flush()
//
// Call [Writer.LinesWritten] to retrieve a count of emitted lines, useful
// for progress reporting or validation.
package output
