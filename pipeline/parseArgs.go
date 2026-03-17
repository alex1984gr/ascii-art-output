// Package pipeline contains all processing stages for ASCII art generation
package pipeline

import (
	"flag"   // Standard Go package for parsing command-line flags
	"fmt"    // Formatted I/O for building error messages
	"io"     // I/O interfaces, used to discard flag error output
	"os"     // OS-level output for printing usage to stderr
	"strings" // String utilities for replacing escape sequences
)

// usage is the message printed when the user provides wrong arguments
const usage = "Usage: go run . [OPTION] [STRING] [BANNER]\n\nEX: go run . --output=<fileName.txt> something standard"

// Config holds all parsed values needed to run the pipeline
type Config struct {
	Input     string // The text to convert to ASCII art
	Font      string // The banner font to use (standard, shadow, thinkertoy)
	OutFile   string // Output file path; empty means write to stdout
	ColorName string // ANSI color name to apply; empty means no color
	Substring string // Part of input to color; empty means color everything
}

// ParseArgs parses os.Args-style arguments into a Config struct.
// Returns an error if arguments are missing, malformed, or use wrong flag format.
func ParseArgs(args []string) (Config, error) {
	// Reject flags passed with a space instead of = (e.g. --output file.txt)
	for _, arg := range args {
		// Check if the user typed --output or --color without the = sign
		if arg == "--output" || arg == "--color" {
			// Return an error with the correct usage message
			return Config{}, fmt.Errorf("%s", usage)
		}
	}

	// Create a new flag set so we don't pollute the global flag state
	fs := flag.NewFlagSet("ascii-art", flag.ContinueOnError)
	// Discard the default flag error output so we control the message ourselves
	fs.SetOutput(io.Discard)

	// Build Config with safe defaults before parsing so we can populate it directly
	cfg := Config{Font: "standard"}

	// Register --output flag directly into cfg to avoid pointer dereference after parse
	fs.StringVar(&cfg.OutFile, "output", "", "")
	// Register --color flag directly into cfg to avoid pointer dereference after parse
	fs.StringVar(&cfg.ColorName, "color", "", "")

	// Parse the provided arguments; stop and return error on unknown flags
	if err := fs.Parse(args); err != nil {
		return Config{}, fmt.Errorf("%s", usage)
	}

	// fs.Args() returns the non-flag arguments left after parsing
	remaining := fs.Args()

	// At least one positional argument (the input string) is required
	if len(remaining) == 0 {
		return Config{}, fmt.Errorf("%s", usage)
	}

	switch len(remaining) {
	case 1:
		// Only the input string was provided: go run . "hello"
		cfg.Input = remaining[0]
	case 2:
		if cfg.ColorName != "" {
			// --color=<color> <substring> <input>: color only the substring
			cfg.Substring = remaining[0] // The part of the text to colorize
			cfg.Input = remaining[1]     // The full input text
		} else {
			// go run . "hello" standard: input + banner font
			cfg.Input = remaining[0] // The text to render
			cfg.Font = remaining[1]  // The banner font name
		}
	default:
		// --color=<color> <substring> <input> [banner]
		cfg.Substring = remaining[0] // The substring to colorize
		cfg.Input = remaining[1]     // The full input text
		cfg.Font = remaining[2]      // The banner font name
	}

	// Replace literal \n sequences in the input with real newline characters
	cfg.Input = strings.ReplaceAll(cfg.Input, "\\n", "\n")

	// Return the fully populated Config
	return cfg, nil
}

// printUsage writes the usage message to stderr
func printUsage() {
	fmt.Fprintln(os.Stderr, usage)
}
