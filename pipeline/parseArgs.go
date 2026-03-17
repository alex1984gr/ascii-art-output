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

// config holds all parsed values needed to run the pipeline
type config struct {
	input     string // The text to convert to ASCII art
	font      string // The banner font to use (standard, shadow, thinkertoy)
	outFile   string // Output file path; empty means write to stdout
	colorName string // ANSI color name to apply; empty means no color
	substring string // Part of input to color; empty means color everything
}

// parseArgs parses os.Args-style arguments into a config struct.
// Returns an error if arguments are missing, malformed, or use wrong flag format.
func parseArgs(args []string) (config, error) {
	// Reject flags passed with a space instead of = (e.g. --output file.txt)
	for _, arg := range args {
		// Check if the user typed --output or --color without the = sign
		if arg == "--output" || arg == "--color" {
			// Return an error with the correct usage message
			return config{}, fmt.Errorf("%s", usage)
		}
	}

	// Create a new flag set so we don't pollute the global flag state
	fs := flag.NewFlagSet("ascii-art", flag.ContinueOnError)
	// Discard the default flag error output so we control the message ourselves
	fs.SetOutput(io.Discard)

	// Build config with safe defaults before parsing so we can populate it directly
	cfg := config{font: "standard"}

	// Register --output flag directly into cfg to avoid pointer dereference after parse
	fs.StringVar(&cfg.outFile, "output", "", "")
	// Register --color flag directly into cfg to avoid pointer dereference after parse
	fs.StringVar(&cfg.colorName, "color", "", "")

	// Parse the provided arguments; stop and return error on unknown flags
	if err := fs.Parse(args); err != nil {
		return config{}, fmt.Errorf("%s", usage)
	}

	// fs.Args() returns the non-flag arguments left after parsing
	remaining := fs.Args()

	// At least one positional argument (the input string) is required
	if len(remaining) == 0 {
		return config{}, fmt.Errorf("%s", usage)
	}

	switch len(remaining) {
	case 1:
		// Only the input string was provided: go run . "hello"
		cfg.input = remaining[0]
	case 2:
		if cfg.colorName != "" {
			// --color=<color> <substring> <input>: color only the substring
			cfg.substring = remaining[0] // The part of the text to colorize
			cfg.input = remaining[1]     // The full input text
		} else {
			// go run . "hello" standard: input + banner font
			cfg.input = remaining[0] // The text to render
			cfg.font = remaining[1]  // The banner font name
		}
	default:
		// --color=<color> <substring> <input> [banner]
		cfg.substring = remaining[0] // The substring to colorize
		cfg.input = remaining[1]     // The full input text
		cfg.font = remaining[2]      // The banner font name
	}

	// Replace literal \n sequences in the input with real newline characters
	cfg.input = strings.ReplaceAll(cfg.input, "\\n", "\n")

	// Return the fully populated config
	return cfg, nil
}

// printUsage writes the usage message to stderr
func printUsage() {
	fmt.Fprintln(os.Stderr, usage)
}
