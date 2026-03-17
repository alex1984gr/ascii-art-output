package tests

import (
	"strings"  // String utilities for checking newline presence in parsed input
	"testing"  // Go standard testing framework

	"ascii-art/pipeline" // The package under test
)

// TestParseArgs_InputOnly verifies that a single positional argument is parsed as the input
// and that the default font "standard" is used when no banner is specified.
func TestParseArgs_InputOnly(t *testing.T) {
	// Call ParseArgs with only one positional argument
	cfg, err := pipeline.ParseArgs([]string{"hello"})
	// Fail immediately if an unexpected error was returned
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Verify the input field was set to the provided string
	if cfg.Input != "hello" {
		t.Errorf("expected input 'hello', got %q", cfg.Input)
	}
	// Verify the font defaults to "standard" when no banner argument is given
	if cfg.Font != "standard" {
		t.Errorf("expected default font 'standard', got %q", cfg.Font)
	}
}

// TestParseArgs_InputAndBanner verifies that two positional arguments are parsed as
// input and banner font respectively.
func TestParseArgs_InputAndBanner(t *testing.T) {
	// Call ParseArgs with input "hello" and banner "shadow"
	cfg, err := pipeline.ParseArgs([]string{"hello", "shadow"})
	// Fail immediately if an unexpected error was returned
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Verify the input field was set correctly
	if cfg.Input != "hello" {
		t.Errorf("expected input 'hello', got %q", cfg.Input)
	}
	// Verify the font field was set to the provided banner name
	if cfg.Font != "shadow" {
		t.Errorf("expected font 'shadow', got %q", cfg.Font)
	}
}

// TestParseArgs_OutputFlag verifies that the --output=<file> flag is correctly parsed
// into the OutFile field of the config.
func TestParseArgs_OutputFlag(t *testing.T) {
	// Call ParseArgs with the --output flag and one positional argument
	cfg, err := pipeline.ParseArgs([]string{"--output=banner.txt", "hello"})
	// Fail immediately if an unexpected error was returned
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Verify the OutFile field was set to the filename provided in the flag
	if cfg.OutFile != "banner.txt" {
		t.Errorf("expected outFile 'banner.txt', got %q", cfg.OutFile)
	}
	// Verify the input field was set to the positional argument
	if cfg.Input != "hello" {
		t.Errorf("expected input 'hello', got %q", cfg.Input)
	}
}

// TestParseArgs_ColorFlagFullOutput verifies that --color=<color> with a single positional
// argument colors the full output (no substring) and sets ColorName correctly.
func TestParseArgs_ColorFlagFullOutput(t *testing.T) {
	// Call ParseArgs with --color flag and one positional argument (no substring)
	cfg, err := pipeline.ParseArgs([]string{"--color=red", "hello"})
	// Fail immediately if an unexpected error was returned
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Verify the ColorName field was set to "red"
	if cfg.ColorName != "red" {
		t.Errorf("expected colorName 'red', got %q", cfg.ColorName)
	}
	// Verify the input field was set to the positional argument
	if cfg.Input != "hello" {
		t.Errorf("expected input 'hello', got %q", cfg.Input)
	}
	// Verify Substring is empty — meaning the entire output will be colored
	if cfg.Substring != "" {
		t.Errorf("expected empty substring, got %q", cfg.Substring)
	}
}

// TestParseArgs_ColorFlagWithSubstring verifies that --color=<color> with two positional
// arguments treats the first as the substring to color and the second as the full input.
func TestParseArgs_ColorFlagWithSubstring(t *testing.T) {
	// Call ParseArgs with --color flag, a substring "he", and the full input "hello"
	cfg, err := pipeline.ParseArgs([]string{"--color=blue", "he", "hello"})
	// Fail immediately if an unexpected error was returned
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Verify the ColorName field was set to "blue"
	if cfg.ColorName != "blue" {
		t.Errorf("expected colorName 'blue', got %q", cfg.ColorName)
	}
	// Verify the Substring field was set to the first positional argument
	if cfg.Substring != "he" {
		t.Errorf("expected substring 'he', got %q", cfg.Substring)
	}
	// Verify the Input field was set to the second positional argument
	if cfg.Input != "hello" {
		t.Errorf("expected input 'hello', got %q", cfg.Input)
	}
}

// TestParseArgs_ColorFlagWithSubstringAndBanner verifies that three positional arguments
// with --color are parsed as substring, input, and banner font respectively.
func TestParseArgs_ColorFlagWithSubstringAndBanner(t *testing.T) {
	// Call ParseArgs with --color, substring "he", input "hello", and banner "shadow"
	cfg, err := pipeline.ParseArgs([]string{"--color=green", "he", "hello", "shadow"})
	// Fail immediately if an unexpected error was returned
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Verify the Substring field was set to the first positional argument
	if cfg.Substring != "he" {
		t.Errorf("expected substring 'he', got %q", cfg.Substring)
	}
	// Verify the Input field was set to the second positional argument
	if cfg.Input != "hello" {
		t.Errorf("expected input 'hello', got %q", cfg.Input)
	}
	// Verify the Font field was set to the third positional argument
	if cfg.Font != "shadow" {
		t.Errorf("expected font 'shadow', got %q", cfg.Font)
	}
}

// TestParseArgs_OutputAndColorCombined verifies that both --output and --color flags
// can be used together alongside a positional input argument.
func TestParseArgs_OutputAndColorCombined(t *testing.T) {
	// Call ParseArgs with both flags and one positional argument
	cfg, err := pipeline.ParseArgs([]string{"--output=out.txt", "--color=red", "hello"})
	// Fail immediately if an unexpected error was returned
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Verify the OutFile field was set correctly
	if cfg.OutFile != "out.txt" {
		t.Errorf("expected outFile 'out.txt', got %q", cfg.OutFile)
	}
	// Verify the ColorName field was set correctly
	if cfg.ColorName != "red" {
		t.Errorf("expected colorName 'red', got %q", cfg.ColorName)
	}
	// Verify the Input field was set to the positional argument
	if cfg.Input != "hello" {
		t.Errorf("expected input 'hello', got %q", cfg.Input)
	}
}

// TestParseArgs_EscapedNewline verifies that the literal sequence \n in the input string
// is converted to a real newline character during parsing.
func TestParseArgs_EscapedNewline(t *testing.T) {
	// Pass a string containing the literal characters backslash and n
	cfg, err := pipeline.ParseArgs([]string{"hello\\nworld"})
	// Fail immediately if an unexpected error was returned
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Verify the Input field contains a real newline character after parsing
	if !strings.Contains(cfg.Input, "\n") {
		t.Errorf("expected input to contain real newline, got %q", cfg.Input)
	}
}

// TestParseArgs_NoArgs_ReturnsError verifies that passing no arguments returns an error
// because at least one positional argument (the input string) is required.
func TestParseArgs_NoArgs_ReturnsError(t *testing.T) {
	// Call ParseArgs with an empty argument slice
	_, err := pipeline.ParseArgs([]string{})
	// Expect a non-nil error; fail if none was returned
	if err == nil {
		t.Fatal("expected error for empty args, got nil")
	}
}

// TestParseArgs_OutputWithSpace_ReturnsError verifies that passing --output as a separate
// argument (with a space instead of =) is rejected with an error.
func TestParseArgs_OutputWithSpace_ReturnsError(t *testing.T) {
	// Pass --output as a standalone flag followed by the filename as a separate argument
	_, err := pipeline.ParseArgs([]string{"--output", "banner.txt", "hello"})
	// Expect a non-nil error because the correct format is --output=<file>
	if err == nil {
		t.Fatal("expected error for --output with space, got nil")
	}
}

// TestParseArgs_ColorWithSpace_ReturnsError verifies that passing --color as a separate
// argument (with a space instead of =) is rejected with an error.
func TestParseArgs_ColorWithSpace_ReturnsError(t *testing.T) {
	// Pass --color as a standalone flag followed by the color name as a separate argument
	_, err := pipeline.ParseArgs([]string{"--color", "red", "hello"})
	// Expect a non-nil error because the correct format is --color=<color>
	if err == nil {
		t.Fatal("expected error for --color with space, got nil")
	}
}

// TestParseArgs_UnknownFlag_ReturnsError verifies that an unrecognised flag causes an error.
func TestParseArgs_UnknownFlag_ReturnsError(t *testing.T) {
	// Pass a flag that is not registered in the flag set
	_, err := pipeline.ParseArgs([]string{"--unknown=foo", "hello"})
	// Expect a non-nil error because the flag package rejects unknown flags
	if err == nil {
		t.Fatal("expected error for unknown flag, got nil")
	}
}

// TestParseArgs_DefaultFont_IsStandard verifies that when no banner argument is provided,
// the Font field defaults to "standard".
func TestParseArgs_DefaultFont_IsStandard(t *testing.T) {
	// Call ParseArgs with only the input argument and no banner
	cfg, err := pipeline.ParseArgs([]string{"hello"})
	// Fail immediately if an unexpected error was returned
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Verify the Font field is "standard" (the hardcoded default)
	if cfg.Font != "standard" {
		t.Errorf("expected default font 'standard', got %q", cfg.Font)
	}
}
