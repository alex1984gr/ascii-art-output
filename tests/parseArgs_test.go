package tests

import (
	"strings"
	"testing"

	"ascii-art/pipeline"
)

func TestParseArgs_InputOnly(t *testing.T) {
	cfg, err := pipeline.ParseArgs([]string{"hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Input != "hello" {
		t.Errorf("expected input 'hello', got %q", cfg.Input)
	}
	if cfg.Font != "standard" {
		t.Errorf("expected default font 'standard', got %q", cfg.Font)
	}
}

func TestParseArgs_InputAndBanner(t *testing.T) {
	cfg, err := pipeline.ParseArgs([]string{"hello", "shadow"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Input != "hello" {
		t.Errorf("expected input 'hello', got %q", cfg.Input)
	}
	if cfg.Font != "shadow" {
		t.Errorf("expected font 'shadow', got %q", cfg.Font)
	}
}

func TestParseArgs_OutputFlag(t *testing.T) {
	cfg, err := pipeline.ParseArgs([]string{"--output=banner.txt", "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.OutFile != "banner.txt" {
		t.Errorf("expected outFile 'banner.txt', got %q", cfg.OutFile)
	}
	if cfg.Input != "hello" {
		t.Errorf("expected input 'hello', got %q", cfg.Input)
	}
}

func TestParseArgs_ColorFlagFullOutput(t *testing.T) {
	cfg, err := pipeline.ParseArgs([]string{"--color=red", "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ColorName != "red" {
		t.Errorf("expected colorName 'red', got %q", cfg.ColorName)
	}
	if cfg.Input != "hello" {
		t.Errorf("expected input 'hello', got %q", cfg.Input)
	}
	if cfg.Substring != "" {
		t.Errorf("expected empty substring, got %q", cfg.Substring)
	}
}

func TestParseArgs_ColorFlagWithSubstring(t *testing.T) {
	cfg, err := pipeline.ParseArgs([]string{"--color=blue", "he", "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ColorName != "blue" {
		t.Errorf("expected colorName 'blue', got %q", cfg.ColorName)
	}
	if cfg.Substring != "he" {
		t.Errorf("expected substring 'he', got %q", cfg.Substring)
	}
	if cfg.Input != "hello" {
		t.Errorf("expected input 'hello', got %q", cfg.Input)
	}
}

func TestParseArgs_ColorFlagWithSubstringAndBanner(t *testing.T) {
	cfg, err := pipeline.ParseArgs([]string{"--color=green", "he", "hello", "shadow"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Substring != "he" {
		t.Errorf("expected substring 'he', got %q", cfg.Substring)
	}
	if cfg.Input != "hello" {
		t.Errorf("expected input 'hello', got %q", cfg.Input)
	}
	if cfg.Font != "shadow" {
		t.Errorf("expected font 'shadow', got %q", cfg.Font)
	}
}

func TestParseArgs_OutputAndColorCombined(t *testing.T) {
	cfg, err := pipeline.ParseArgs([]string{"--output=out.txt", "--color=red", "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.OutFile != "out.txt" {
		t.Errorf("expected outFile 'out.txt', got %q", cfg.OutFile)
	}
	if cfg.ColorName != "red" {
		t.Errorf("expected colorName 'red', got %q", cfg.ColorName)
	}
	if cfg.Input != "hello" {
		t.Errorf("expected input 'hello', got %q", cfg.Input)
	}
}

func TestParseArgs_EscapedNewline(t *testing.T) {
	cfg, err := pipeline.ParseArgs([]string{"hello\\nworld"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(cfg.Input, "\n") {
		t.Errorf("expected input to contain real newline, got %q", cfg.Input)
	}
}

func TestParseArgs_NoArgs_ReturnsError(t *testing.T) {
	_, err := pipeline.ParseArgs([]string{})
	if err == nil {
		t.Fatal("expected error for empty args, got nil")
	}
}

func TestParseArgs_OutputWithSpace_ReturnsError(t *testing.T) {
	_, err := pipeline.ParseArgs([]string{"--output", "banner.txt", "hello"})
	if err == nil {
		t.Fatal("expected error for --output with space, got nil")
	}
}

func TestParseArgs_ColorWithSpace_ReturnsError(t *testing.T) {
	_, err := pipeline.ParseArgs([]string{"--color", "red", "hello"})
	if err == nil {
		t.Fatal("expected error for --color with space, got nil")
	}
}

func TestParseArgs_UnknownFlag_ReturnsError(t *testing.T) {
	_, err := pipeline.ParseArgs([]string{"--unknown=foo", "hello"})
	if err == nil {
		t.Fatal("expected error for unknown flag, got nil")
	}
}

func TestParseArgs_DefaultFont_IsStandard(t *testing.T) {
	cfg, err := pipeline.ParseArgs([]string{"hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Font != "standard" {
		t.Errorf("expected default font 'standard', got %q", cfg.Font)
	}
}
