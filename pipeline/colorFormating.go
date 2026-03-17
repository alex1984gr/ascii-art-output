// Package pipeline contains all the processing stages for ASCII art generation
package pipeline

import (
	"fmt"     // Formatted I/O for building error messages
	"strings" // String utilities for ReplaceAll and ToLower
)

// ansiColors maps human-readable color names to their ANSI terminal escape sequences.
// The terminal reads these sequences and changes the text color accordingly.
var ansiColors = map[string]string{
	// Basic 8 colors — standard ANSI codes 30-37
	"black":   "\033[30m", // ANSI escape code for black text
	"red":     "\033[31m", // ANSI escape code for red text
	"green":   "\033[32m", // ANSI escape code for green text
	"yellow":  "\033[33m", // ANSI escape code for yellow text
	"blue":    "\033[34m", // ANSI escape code for blue text
	"magenta": "\033[35m", // ANSI escape code for magenta text
	"cyan":    "\033[36m", // ANSI escape code for cyan text
	"white":   "\033[37m", // ANSI escape code for white text

	// Bright/bold variants — ANSI codes 90-97
	"bright_black":   "\033[90m", // Bright black (appears as dark gray)
	"bright_red":     "\033[91m", // Bright red
	"bright_green":   "\033[92m", // Bright green
	"bright_yellow":  "\033[93m", // Bright yellow
	"bright_blue":    "\033[94m", // Bright blue
	"bright_magenta": "\033[95m", // Bright magenta
	"bright_cyan":    "\033[96m", // Bright cyan
	"bright_white":   "\033[97m", // Bright white

	// Extended 256-color palette — format: \033[38;5;<n>m
	"orange":       "\033[38;5;208m", // Orange (256-color index 208)
	"dark_orange":  "\033[38;5;166m", // Dark orange (256-color index 166)
	"light_orange": "\033[38;5;214m", // Light orange (256-color index 214)
	"purple":       "\033[38;5;129m", // Purple (256-color index 129)
	"dark_purple":  "\033[38;5;54m",  // Dark purple (256-color index 54)
	"light_purple": "\033[38;5;141m", // Light purple (256-color index 141)
	"pink":         "\033[38;5;213m", // Pink (256-color index 213)
	"hot_pink":     "\033[38;5;198m", // Hot pink (256-color index 198)
	"light_pink":   "\033[38;5;217m", // Light pink (256-color index 217)
	"brown":        "\033[38;5;130m", // Brown (256-color index 130)
	"dark_brown":   "\033[38;5;94m",  // Dark brown (256-color index 94)
	"light_brown":  "\033[38;5;180m", // Light brown (256-color index 180)
	"gold":         "\033[38;5;220m", // Gold (256-color index 220)
	"silver":       "\033[38;5;250m", // Silver (256-color index 250)
	"bronze":       "\033[38;5;136m", // Bronze (256-color index 136)

	// Shades of gray
	"gray":       "\033[90m",        // Gray (alias for bright black)
	"grey":       "\033[90m",        // Grey (alternative British spelling)
	"dark_gray":  "\033[38;5;240m",  // Dark gray (256-color index 240)
	"dark_grey":  "\033[38;5;240m",  // Dark grey (alternative spelling)
	"light_gray": "\033[38;5;252m",  // Light gray (256-color index 252)
	"light_grey": "\033[38;5;252m",  // Light grey (alternative spelling)

	// Shades of red
	"dark_red":  "\033[38;5;88m",  // Dark red (256-color index 88)
	"light_red": "\033[38;5;203m", // Light red (256-color index 203)
	"crimson":   "\033[38;5;160m", // Crimson (256-color index 160)
	"maroon":    "\033[38;5;52m",  // Maroon (256-color index 52)
	"salmon":    "\033[38;5;209m", // Salmon (256-color index 209)
	"coral":     "\033[38;5;203m", // Coral (256-color index 203)

	// Shades of green
	"dark_green":   "\033[38;5;22m",  // Dark green (256-color index 22)
	"light_green":  "\033[38;5;120m", // Light green (256-color index 120)
	"lime":         "\033[38;5;154m", // Lime (256-color index 154)
	"olive":        "\033[38;5;58m",  // Olive (256-color index 58)
	"forest_green": "\033[38;5;28m",  // Forest green (256-color index 28)
	"sea_green":    "\033[38;5;85m",  // Sea green (256-color index 85)

	// Shades of blue
	"dark_blue":  "\033[38;5;18m",  // Dark blue (256-color index 18)
	"light_blue": "\033[38;5;117m", // Light blue (256-color index 117)
	"navy":       "\033[38;5;17m",  // Navy blue (256-color index 17)
	"sky_blue":   "\033[38;5;117m", // Sky blue (256-color index 117)
	"royal_blue": "\033[38;5;63m",  // Royal blue (256-color index 63)
	"steel_blue": "\033[38;5;67m",  // Steel blue (256-color index 67)

	// Shades of cyan
	"dark_cyan":  "\033[38;5;30m",  // Dark cyan (256-color index 30)
	"light_cyan": "\033[38;5;123m", // Light cyan (256-color index 123)
	"aqua":       "\033[38;5;51m",  // Aqua (256-color index 51)
	"turquoise":  "\033[38;5;80m",  // Turquoise (256-color index 80)
	"teal":       "\033[38;5;30m",  // Teal (256-color index 30)

	// Shades of yellow
	"dark_yellow":  "\033[38;5;136m", // Dark yellow (256-color index 136)
	"light_yellow": "\033[38;5;228m", // Light yellow (256-color index 228)
	"khaki":        "\033[38;5;185m", // Khaki (256-color index 185)

	// Shades of magenta
	"dark_magenta":  "\033[38;5;90m",  // Dark magenta (256-color index 90)
	"light_magenta": "\033[38;5;213m", // Light magenta (256-color index 213)
	"violet":        "\033[38;5;177m", // Violet (256-color index 177)
	"indigo":        "\033[38;5;54m",  // Indigo (256-color index 54)

	// Special named colors
	"beige":     "\033[38;5;230m", // Beige (256-color index 230)
	"cream":     "\033[38;5;230m", // Cream (256-color index 230)
	"ivory":     "\033[38;5;255m", // Ivory (256-color index 255)
	"lavender":  "\033[38;5;183m", // Lavender (256-color index 183)
	"mint":      "\033[38;5;121m", // Mint (256-color index 121)
	"peach":     "\033[38;5;217m", // Peach (256-color index 217)
	"rose":      "\033[38;5;211m", // Rose (256-color index 211)
	"ruby":      "\033[38;5;161m", // Ruby (256-color index 161)
	"emerald":   "\033[38;5;35m",  // Emerald (256-color index 35)
	"sapphire":  "\033[38;5;26m",  // Sapphire (256-color index 26)
	"amethyst":  "\033[38;5;134m", // Amethyst (256-color index 134)
	"topaz":     "\033[38;5;214m", // Topaz (256-color index 214)
	"jade":      "\033[38;5;35m",  // Jade (256-color index 35)
	"amber":     "\033[38;5;214m", // Amber (256-color index 214)
	"chocolate": "\033[38;5;94m",  // Chocolate (256-color index 94)
	"coffee":    "\033[38;5;94m",  // Coffee (256-color index 94)
	"sand":      "\033[38;5;215m", // Sand (256-color index 215)
	"tan":       "\033[38;5;180m", // Tan (256-color index 180)

	// Reset — restores the terminal's default text color
	"reset": "\033[0m", // ANSI reset code: clears all active color attributes
}

// ColorLines applies ANSI color to a slice of plain text lines.
// Used by tests and as a simpler alternative when no banner context is needed.
// If substring is empty, every line is fully colored.
// If substring is non-empty, only occurrences of that literal string are colored.
func ColorLines(lines []string, color, substring string) ([]string, error) {
	// Convert the color name to lowercase so lookups are case-insensitive
	code, ok := ansiColors[strings.ToLower(color)]
	// If the color name is not in the map, return a descriptive error
	if !ok {
		return nil, fmt.Errorf("invalid color: %s", color)
	}

	// Allocate a new slice with the same length to hold the colored lines
	colored := make([]string, len(lines))

	// Iterate over every line in the input slice
	for i, line := range lines {
		if substring == "" {
			// No substring: wrap the entire line between the color code and the reset code
			colored[i] = fmt.Sprintf("%s%s%s", code, line, ansiColors["reset"])
		} else {
			// Substring provided: replace every occurrence of it with a colored version
			colored[i] = strings.ReplaceAll(line, substring, fmt.Sprintf("%s%s%s", code, substring, ansiColors["reset"]))
		}
	}

	// Return the slice of colored lines and no error
	return colored, nil
}

// ColorLinesWithBanner applies ANSI color to ASCII art lines using the banner for context.
// If substring is empty, every line is fully colored.
// If substring is non-empty, it is first rendered to ASCII art using the banner,
// and then each matching ASCII art row is colored inside the full output.
func ColorLinesWithBanner(lines []string, color, substring string, banner map[string][]string) ([]string, error) {
	// Convert the color name to lowercase and look it up in the map
	code, ok := ansiColors[strings.ToLower(color)]
	// Return an error if the color name is not recognised
	if !ok {
		return nil, fmt.Errorf("invalid color: %s", color)
	}

	// If no substring was given, color every line in full
	if substring == "" {
		// Allocate result slice with the same length as the input
		colored := make([]string, len(lines))
		for i, line := range lines {
			// Wrap each line with the color code and the reset code
			colored[i] = fmt.Sprintf("%s%s%s", code, line, ansiColors["reset"])
		}
		// Return the fully colored lines
		return colored, nil
	}

	// Tokenize the substring into individual characters so it can be rendered
	subTokens := Tokenize(substring)
	// Render the substring tokens into ASCII art rows using the same banner font
	subLines := RenderLines(subTokens, banner)

	// Allocate result slice with the same length as the full output
	colored := make([]string, len(lines))

	// Process each output line and selectively color the substring's ASCII art pattern
	for i, line := range lines {
		// Determine which row of the 8-row ASCII art block we are on (cycles 0-7)
		subRow := i % 8
		// Only attempt replacement if this row index exists in the substring's rendered art
		if subRow < len(subLines) {
			// Retrieve the ASCII art pattern for this row of the substring
			pattern := subLines[subRow]
			// Replace every occurrence of the pattern in the current line with a colored version
			colored[i] = strings.ReplaceAll(line, pattern, fmt.Sprintf("%s%s%s", code, pattern, ansiColors["reset"]))
		} else {
			// No pattern for this row — keep the line unchanged
			colored[i] = line
		}
	}

	// Return the selectively colored lines
	return colored, nil
}
