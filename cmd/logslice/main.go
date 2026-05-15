package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/reader"
	"github.com/yourorg/logslice/internal/slicer"
	"github.com/yourorg/logslice/internal/timestamp"
)

const timeLayout = "2006-01-02T15:04:05"

func main() {
	var (
		start     = flag.String("start", "", "Start timestamp (RFC3339 or YYYY-MM-DDTHH:MM:SS)")
		end       = flag.String("end", "", "End timestamp (RFC3339 or YYYY-MM-DDTHH:MM:SS)")
		outFile   = flag.String("out", "", "Output file path (default: stdout)")
		format    = flag.String("format", "", "Explicit timestamp format (optional)")
		verbose   = flag.Bool("verbose", false, "Print summary after slicing")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: logslice [options] <logfile>\n\nOptions:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}
	if *start == "" || *end == "" {
		fmt.Fprintln(os.Stderr, "error: --start and --end are required")
		os.Exit(1)
	}

	parseTime := func(s string) (time.Time, error) {
		t, err := time.Parse(time.RFC3339, s)
		if err == nil {
			return t, nil
		}
		return time.ParseInLocation(timeLayout, s, time.UTC)
	}

	startTime, err := parseTime(*start)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: invalid --start: %v\n", err)
		os.Exit(1)
	}
	endTime, err := parseTime(*end)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: invalid --end: %v\n", err)
		os.Exit(1)
	}

	logPath := flag.Arg(0)
	f, err := os.Open(logPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: cannot open file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	var parser *timestamp.Parser
	if *format != "" {
		parser = timestamp.NewParser(*format)
	} else {
		parser = timestamp.NewParser("")
	}

	extractor := timestamp.NewExtractor(parser)
	scanner := reader.NewScanner(f, extractor)

	var w *output.Writer
	if *outFile != "" {
		w, err = output.NewFileWriter(*outFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: cannot open output file: %v\n", err)
			os.Exit(1)
		}
	} else {
		w = output.NewWriter(os.Stdout)
	}
	defer w.Close()

	s := slicer.New(scanner, w)
	n, err := s.Slice(startTime, endTime)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if *verbose {
		fmt.Fprintf(os.Stderr, "logslice: wrote %d lines\n", n)
	}
}
