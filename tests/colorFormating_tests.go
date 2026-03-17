package tests

import (
	"strings"  // String utilities for checking substrings in colored output
	"testing"  // Go standard testing framework

	"ascii-art/pipeline" // The package under test
)

// TestColorLines_FullLine verifies that ColorLines wraps every line with the correct ANSI color codes
// when no substring is provided (i.e., the entire line should be colored).
func TestColorLines_FullLine(t *testing.T) {
	// Define two plain text lines to be colored
	lines := []string{"Hello", "World"}
	// Call ColorLines with color "red" and no substring — should color entire lines
	colored, err := pipeline.ColorLines(lines, "red", "")
	// Fail immediately if an unexpected error occurred
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify each colored line starts with the red ANSI code and ends with the reset code
	for _, line := range colored {
		// \033[31m is the ANSI escape for red; \033[0m resets the color
		if !strings.HasPrefix(line, "\033[31m") || !strings.HasSuffix(line, "\033[0m") {
			t.Errorf("line not correctly colored: %q", line)
		}
	}
}

// TestColorLines_Substring verifies that ColorLines only colors occurrences of the given substring,
// leaving the rest of each line uncolored.
func TestColorLines_Substring(t *testing.T) {
	// Define lines that both contain the substring "kit"
	lines := []string{"kitten", "a king kitten"}
	// Call ColorLines with color "blue" and substring "kit"
	colored, err := pipeline.ColorLines(lines, "blue", "kit")
	// Fail immediately if an unexpected error occurred
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// The expected colored substring pattern: blue code + "kit" + reset code
	expectedPrefix := "\033[34m" // ANSI blue
	expectedSuffix := "\033[0m"  // ANSI reset

	// Check that each line contains at least one colored occurrence of "kit"
	for i, line := range colored {
		// Count how many times the colored "kit" pattern appears in this line
		count := strings.Count(line, expectedPrefix+"kit"+expectedSuffix)
		// If the pattern is absent, the substring was not colored correctly
		if count == 0 {
			t.Errorf("substring not colored in line %d: %q", i, line)
		}
	}
}

// TestColorLines_InvalidColor verifies that ColorLines returns an error for an unrecognised color name.
func TestColorLines_InvalidColor(t *testing.T) {
	// Define a single test line
	lines := []string{"test"}
	// Call ColorLines with a color name that does not exist in the ansiColors map
	_, err := pipeline.ColorLines(lines, "invalidColor", "")
	// Expect a non-nil error; fail if none was returned
	if err == nil {
		t.Fatalf("expected error for invalid color, got nil")
	}
}

// TestColorLines_Multiline verifies that ColorLines correctly colors every line
// in a multi-line slice when no substring is provided.
func TestColorLines_Multiline(t *testing.T) {
	// Define three plain text lines
	lines := []string{"line1", "line2", "line3"}
	// Call ColorLines with color "green" and no substring
	colored, err := pipeline.ColorLines(lines, "green", "")
	// Fail immediately if an unexpected error occurred
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify every line is wrapped with the green ANSI code and the reset code
	for _, line := range colored {
		// \033[32m is the ANSI escape for green
		if !strings.HasPrefix(line, "\033[32m") || !strings.HasSuffix(line, "\033[0m") {
			t.Errorf("multiline line not correctly colored: %q", line)
		}
	}
}
