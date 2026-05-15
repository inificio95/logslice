package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// buildBinary compiles the CLI binary into a temp dir and returns its path.
func buildBinary(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	bin := filepath.Join(dir, "logslice")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	return bin
}

func writeTempLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "test-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	return f.Name()
}

func TestCLINoArgs(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin)
	out, _ := cmd.CombinedOutput()
	if !strings.Contains(string(out), "Usage") {
		t.Errorf("expected usage output, got: %s", out)
	}
}

func TestCLIMissingStartEnd(t *testing.T) {
	bin := buildBinary(t)
	log := writeTempLog(t, []string{"2024-01-01T00:00:00 hello"})
	cmd := exec.Command(bin, log)
	out, _ := cmd.CombinedOutput()
	if !strings.Contains(string(out), "--start") {
		t.Errorf("expected error about --start, got: %s", out)
	}
}

func TestCLIBasicSlice(t *testing.T) {
	bin := buildBinary(t)
	lines := []string{
		"2024-06-01T10:00:00 first line",
		"2024-06-01T10:05:00 second line",
		"2024-06-01T10:10:00 third line",
		"2024-06-01T10:15:00 fourth line",
	}
	log := writeTempLog(t, lines)
	cmd := exec.Command(bin,
		"--start", "2024-06-01T10:04:00",
		"--end", "2024-06-01T10:11:00",
		log,
	)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := string(out)
	if !strings.Contains(result, "second line") {
		t.Errorf("expected 'second line' in output, got: %s", result)
	}
	if !strings.Contains(result, "third line") {
		t.Errorf("expected 'third line' in output, got: %s", result)
	}
	if strings.Contains(result, "first line") {
		t.Errorf("did not expect 'first line' in output, got: %s", result)
	}
	if strings.Contains(result, "fourth line") {
		t.Errorf("did not expect 'fourth line' in output, got: %s", result)
	}
}
