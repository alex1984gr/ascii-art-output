// Package pipeline orchestrates all ASCII art generation stages in order
package pipeline

import (
	"fmt"           // Formatted I/O for writing error messages
	"io"            // I/O interfaces (Writer) used for output destination
	"os"            // OS-level file creation and stderr access
	"path/filepath" // File path utilities for safe filename handling
)

// Run is the single entry point called by main.
// It orchestrates every pipeline stage in order and returns 0 on success, 1 on error.
// Stages: parseArgs → ValidateInput → LoadBanner → Tokenize → RenderLines → (color) → WriteOutput
func Run(args []string, stdout io.Writer) int {
	// Stage 1: parse flags (--output, --color) and positional arguments using the flag package
	cfg, err := parseArgs(args)
	if err != nil {
		// parseArgs returns the usage string as the error message
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	// Stage 2: validate the input text (rejects empty, too-long, or non-printable characters)
	if err := ValidateInput(cfg.input); err != nil {
		fmt.Fprintf(os.Stderr, "invalid input: %v\n", err.Error())
		return 1
	}

	// Stage 3: load the chosen banner font file into a map of character → 8 ASCII art lines
	banner, err := LoadBanner(cfg.font)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed loading banner: %v\n", err.Error())
		return 1
	}

	// Stage 4+5: split input into individual character tokens, then render ASCII art lines
	lines := RenderLines(Tokenize(cfg.input), banner)

	// Stage 6 (optional): apply ANSI color after rendering so rendering logic stays pure
	// Color formatting is applied after rendering, ensuring that rendering logic remains
	// pure and independent of presentation concerns.
	if cfg.colorName != "" {
		// Color the full output, or only the ASCII art rows that match the substring
		lines, err = ColorLinesWithBanner(lines, cfg.colorName, cfg.substring, banner)
		if err != nil {
			fmt.Fprintf(os.Stderr, "color error: %v\n", err.Error())
			return 1
		}
	}

	// Stage 7: determine the output destination — file or stdout
	var w io.Writer = stdout // default: write to stdout
	if cfg.outFile != "" {
		// filepath.Base strips any directory components from the user-supplied filename,
		// restricting output to the current directory and preventing path traversal attacks
		// (e.g. --output=../../etc/passwd becomes just "passwd" and is written locally)
		safeName := filepath.Base(cfg.outFile)
		// Create (or overwrite) the output file using only the safe base filename
		f, err := os.Create(safeName) //nolint
		if err != nil {
			// Do not include the filename in the error to avoid echoing user input to stderr
			fmt.Fprintln(os.Stderr, "failed creating output file")
			return 1
		}
		// defer guarantees the file handle is released when Run returns, even on error paths
		defer f.Close() //nolint
		w = f // redirect all subsequent writes to the file instead of stdout
	}

	// Stage 7: write every ASCII art line to the chosen destination, separated by newlines
	if err := WriteOutput(lines, w); err != nil {
		fmt.Fprintf(os.Stderr, "failed writing output: %v\n", err.Error())
		return 1
	}

	// All stages completed successfully; exit code 0 signals success to the shell
	return 0
}
