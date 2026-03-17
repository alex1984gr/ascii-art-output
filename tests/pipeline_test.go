package tests

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"ascii-art/pipeline"
)

func TestPipeline_Run_FullFlow(t *testing.T) {
	// Provide arguments as they would be from command line
	args := []string{
		"--color=red",
		"Hello",
		"standard",
	}

	// Use bytes.Buffer instead of stdout for testing
	var out bytes.Buffer

	// Run the pipeline with test arguments
	exitCode := pipeline.Run(args, &out)
	// Verify the pipeline succeeded (exit code 0)
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	// Get the output as a string
	output := out.String()

	// Check that output contains the ANSI code for red color
	if !strings.Contains(output, "\033[31m") {
		t.Errorf("output does not contain red ANSI code: %q", output)
	}

	// Check that output is not empty
	if len(output) == 0 {
		t.Errorf("output is empty")
	}
}

func TestPipeline_Run_InvalidFont(t *testing.T) {
	args := []string{
		"--font=nonexistent",
		"Test",
	}

	var out bytes.Buffer
	exitCode := pipeline.Run(args, &out)
	if exitCode == 0 {
		t.Fatalf("expected non-zero exit code for invalid font")
	}
}

func TestPipeline_Run_NoInput(t *testing.T) {
	args := []string{}

	var out bytes.Buffer
	exitCode := pipeline.Run(args, &out)
	if exitCode == 0 {
		t.Fatalf("expected non-zero exit code for empty input")
	}
}

func TestPipeline_Run_OutputLongFlag(t *testing.T) {
	// filepath.Base is applied by writerFor, so the file lands in the current working directory
	fileName := "long-flag.txt"
	t.Cleanup(func() { os.Remove(fileName) }) // Remove the file after the test

	args := []string{"--output=" + fileName, "Hello"}
	var out bytes.Buffer

	exitCode := pipeline.Run(args, &out)
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("expected output file to be created: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected output file to contain ascii art")
	}
}

func TestPipeline_Run_OutputWithBanner(t *testing.T) {
	// filepath.Base is applied by writerFor, so the file lands in the current working directory
	fileName := "banner-flag.txt"
	t.Cleanup(func() { os.Remove(fileName) }) // Remove the file after the test

	args := []string{"--output=" + fileName, "Hello", "shadow"}
	var out bytes.Buffer

	exitCode := pipeline.Run(args, &out)
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("expected output file to be created: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected output file to contain ascii art")
	}
}

func TestPipeline_Run_EscapedNewlineBecomesMultiline(t *testing.T) {
	var out bytes.Buffer
	exitCode := pipeline.Run([]string{"First\\nTest", "shadow"}, &out)
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	lines := strings.Split(strings.TrimRight(out.String(), "\n"), "\n")
	if len(lines) <= 8 {
		t.Fatalf("expected multiline banner output, got %d lines", len(lines))
	}
}

func TestPipeline_Run_PositionalBannerShadow(t *testing.T) {
	var stdOut bytes.Buffer
	if code := pipeline.Run([]string{"A", "standard"}, &stdOut); code != 0 {
		t.Fatalf("expected standard run to pass, got %d", code)
	}

	var shadowOut bytes.Buffer
	if code := pipeline.Run([]string{"A", "shadow"}, &shadowOut); code != 0 {
		t.Fatalf("expected shadow run to pass, got %d", code)
	}

	if stdOut.String() == shadowOut.String() {
		t.Fatal("expected shadow banner output to differ from standard")
	}
}
