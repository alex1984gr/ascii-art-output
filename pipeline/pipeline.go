// Package pipeline orchestrates the ASCII art generation process
package pipeline

import (
	"fmt"     // Formatted I/O functions for error messages
	"io"      // Basic I/O interfaces like Writer
	"os"      // Operating system functionality for file operations
	"strings" // String manipulation functions
)

// Run executes the complete ASCII art pipeline from input to output.
// It parses command-line arguments, validates input, loads the banner,
// renders ASCII art, applies optional color, and writes the output.
func Run(args []string, stdout io.Writer) int {
	// Initialize default configuration values
	font := "standard" // Default banner font
	outFile := ""      // Empty means write to stdout
	fontSet := false
	var input string     // Text to convert to ASCII art
	var colorName string // Color to apply
	var substring string // Substring to color (empty means color entire output)

	// Parse command-line arguments manually
	// Expected formats:
	// 1. go run . --color=<color> <substring> "text"
	// 2. go run . --color=<color> "text" (no substring, color entire output)
	// 3. go run . "text" (no color)
	for i := 0; i < len(args); i++ {
		arg := args[i]
		// Check for --font= flag
		if strings.HasPrefix(arg, "--font=") {
			// Extract font name by removing prefix
			font = strings.TrimPrefix(arg, "--font=")
			fontSet = true
			// Check for --out= or --output= flag
		} else if strings.HasPrefix(arg, "--out=") || strings.HasPrefix(arg, "--output=") {
			if strings.HasPrefix(arg, "--output=") {
				outFile = strings.TrimPrefix(arg, "--output=")
			} else {
				outFile = strings.TrimPrefix(arg, "--out=")
			}
			// Check for --color= flag
		} else if strings.HasPrefix(arg, "--color=") {
			// Extract color name
			colorName = strings.TrimPrefix(arg, "--color=")
			// Next argument should be either substring or input text
			if i+1 < len(args) {
				i++
				// If there's another argument after this, this is the substring
				if i+1 < len(args) {
					substring = args[i]
					i++
					input = args[i]
				} else {
					// Only one argument left, it's the input text
					input = args[i]
				}
			}
			// If not a flag and input not set yet, this is the input text
		} else if input == "" {
			input = arg
		} else if !fontSet && isBannerName(arg) {
			// Support positional [STRING] [BANNER] when --font is not provided.
			font = arg
		}
	}

	// Validate that user provided input
	if input == "" {
		// Print error to stderr
		fmt.Fprintln(os.Stderr, "no input provided")
		return 1 // Return failure exit code
	}

	// Convert escaped newlines from CLI input into actual newline characters.
	input = strings.ReplaceAll(input, "\\n", "\n")

	// Validate input text (length, allowed characters, etc.)
	if err := ValidateInput(input); err != nil {
		fmt.Fprintln(os.Stderr, "invalid input:", err)
		return 1
	}

	// Tokenize input string into individual characters
	tokens := Tokenize(input)
	// Load banner file and parse into character-to-glyph map
	banner, err := LoadBanner(font)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed loading banner:", err)
		return 1
	}

	// Render ASCII art lines from tokens using banner
	lines := RenderLines(tokens, banner)

	// Apply color if --color flag was provided
	if colorName != "" {
		// Apply ANSI color codes to rendered lines
		// Pass the banner so ColorLines can render the substring if needed
		lines, err = ColorLinesWithBanner(lines, colorName, substring, banner)
		if err != nil {
			fmt.Fprintln(os.Stderr, "color error:", err)
			return 1
		}
	}

	// Determine output destination (file or stdout)
	var w io.Writer = stdout
	if outFile != "" {
		// Create output file (overwrites if exists)
		f, err := os.Create(outFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed creating output file:", err)
			return 1
		}
		// Ensure file is closed when function returns
		defer f.Close()
		// Set writer to file instead of stdout
		w = f
	}

	// Write final ASCII art lines to output destination
	if err := WriteOutput(lines, w); err != nil {
		fmt.Fprintln(os.Stderr, "failed writing output:", err)
		return 1
	}

	// Return success exit code
	return 0
}

func isBannerName(name string) bool {
	return name == "standard" || name == "shadow" || name == "thinkertoy"
}
