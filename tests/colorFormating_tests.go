package tests

import (
	"strings"
	"testing"

	"ascii-art/pipeline"
)

func TestColorLines_FullLine(t *testing.T) {
	lines := []string{"Hello", "World"}
	colored, err := pipeline.ColorLines(lines, "red", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, line := range colored {
		if !strings.HasPrefix(line, "\033[31m") || !strings.HasSuffix(line, "\033[0m") {
			t.Errorf("line not correctly colored: %q", line)
		}
	}
}

func TestColorLines_Substring(t *testing.T) {
	lines := []string{"kitten", "a king kitten"}
	colored, err := pipeline.ColorLines(lines, "blue", "kit")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedPrefix := "\033[34m"
	expectedSuffix := "\033[0m"

	for i, line := range colored {
		count := strings.Count(line, expectedPrefix+"kit"+expectedSuffix)
		if count == 0 {
			t.Errorf("substring not colored in line %d: %q", i, line)
		}
	}
}

func TestColorLines_InvalidColor(t *testing.T) {
	lines := []string{"test"}
	_, err := pipeline.ColorLines(lines, "invalidColor", "")
	if err == nil {
		t.Fatalf("expected error for invalid color, got nil")
	}
}

func TestColorLines_Multiline(t *testing.T) {
	lines := []string{"line1", "line2", "line3"}
	colored, err := pipeline.ColorLines(lines, "green", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, line := range colored {
		if !strings.HasPrefix(line, "\033[32m") || !strings.HasSuffix(line, "\033[0m") {
			t.Errorf("multiline line not correctly colored: %q", line)
		}
	}
}
