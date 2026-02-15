# ASCII-Art-Color
ASCII-Art-Color is a Go CLI project that renders text into ASCII art using predefined banner styles (`standard`, `shadow`, `thinkertoy`) and adds optional ANSI color support. The project uses a clear pipeline architecture and Test-Driven Development (TDD) principles for maintainability and correctness. It runs entirely in the command-line environment.

## Usage
Run without color:

go run . "Hello" standard

Run with color:

go run . --color=red "Hello" shadow 

Build and run binary:

go build -o ascii-art-color
./ascii-art-color "Hello" thinkertoy

## Pipeline Architecture
Input

  ↓
ValidateInput
  ↓
Tokenize
  ↓
LoadBanner
  ↓
RenderLines
  ↓
Color Formatting (optional)
  ↓
WriteOutput

- ValidateInput → ensures input length and allowed characters

- Tokenize → splits input into individual runes/tokens

- LoadBanner → loads the banner file into a structured map

- RenderLines → constructs ASCII art lines from tokens and banner

- Color Formatting (optional) → applies ANSI color codes

- WriteOutput → prints or writes the final ASCII art to file/stdout


## Implementation Details

Banners map each printable character (ASCII 32–126) to 8 ASCII art lines:

map[rune][]string


Rendering algorithm:

1. Iterate each input character

2. Retrieve banner representation

3. Concatenate rows for all characters (8 rows total)

4. Repeat for multiline input
   Time complexity: O(n × h), where n = number of characters, h = banner height (8).
   Color formatting is applied after rendering using ANSI codes (e.g., \033[31mASCII\033[0m) so rendering logic stays independent.

## Testing (TDD)

Development uses TDD: write failing tests, implement minimal logic, refactor while keeping tests passing.
Test coverage includes banner loading, tokenization, rendering, multiline input, color formatting, and error handling. Independent testing ensures predictable behavior and safe refactoring.

## Error Handling

The program returns meaningful errors for invalid banner selection, missing arguments, unsupported characters, and file loading failures. main acts only as entry point for fatal errors.

## Design Principles

Single Responsibility Principle, Separation of Concerns, Modular structure, Predictable rendering, Maintainable and extensible codebase.