// Package slicer implements the core log-slicing logic for logslice.
//
// It combines the reader.Scanner (which provides indexed access to log lines
// and binary-search helpers) with the timestamp extraction and parsing
// pipeline to efficiently locate the first line whose timestamp is >= a
// given "from" value, then streams forward collecting lines until the
// timestamp exceeds the "to" value.
//
// Typical usage:
//
//	f, _ := os.Open("app.log")
//	info, _ := f.Stat()
//	s, _ := slicer.New(f, info.Size(), slicer.Options{})
//	lines, _ := s.Slice(from, to)
//	for _, l := range lines {
//		fmt.Println(l)
//	}
package slicer
