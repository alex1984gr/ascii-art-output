# ASCII-Art-Output

ASCII-Art-Output is a Go CLI project that renders text into ASCII art using predefined banner styles (`standard`, `shadow`, `thinkertoy`), adds optional ANSI color support, and can write the result to a file using the `--output` flag. The project uses a clear pipeline architecture and Test-Driven Development (TDD) principles for maintainability and correctness. It runs entirely in the command-line environment.

## Usage

Run without options:

```
go run . "Hello" standard
```

Write output to a file:

```
go run . --output=banner.txt "Hello" standard
```

Run with color (entire output):

```
go run . --color=red "Hello" shadow
```

Run with color on a substring only:

```
go run . --color=red "He" "Hello" shadow
```

Combine output file and color:

```
go run . --output=banner.txt --color=blue "Hello" thinkertoy
```

Build and run binary:

```
go build -o ascii-art-output
./ascii-art-output --output=banner.txt "Hello" standard
```

> **Note:** Flags must use the `=` format. `--output file.txt` (with a space) is not accepted and will print a usage error.

## Flags

| Flag | Format | Description |
|------|--------|-------------|
| `--output` | `--output=<fileName.txt>` | Write ASCII art to a file instead of stdout |
| `--color` | `--color=<colorName>` | Apply ANSI color to the output |

## Pipeline Architecture

```
Input
  ↓
parseArgs        ← parses flags and positional arguments (flag package)
  ↓
ValidateInput    ← ensures input length and allowed characters
  ↓
Tokenize         ← splits input into individual runes/tokens
  ↓
LoadBanner       ← loads the banner file into a structured map
  ↓
RenderLines      ← constructs ASCII art lines from tokens and banner
  ↓
ColorLinesWithBanner (optional) ← applies ANSI color codes
  ↓
WriteOutput      ← prints or writes the final ASCII art to file/stdout
```

## File Structure

```
pipeline/
  parseArgs.go      ← flag parsing and argument validation (uses standard flag package)
  pipeline.go       ← orchestrates all pipeline stages
  validateInput.go  ← input validation logic
  tokenize.go       ← tokenization logic
  loadBanner.go     ← banner file loading
  renderLines.go    ← ASCII art rendering
  colorFormating.go ← ANSI color formatting (ColorLines, ColorLinesWithBanner)
  writeOutput.go    ← output writing to io.Writer
```

## Argument Parsing

Argument parsing uses the standard Go `flag` package via a dedicated `parseArgs.go` file. This keeps `pipeline.go` focused only on orchestration.

Supported argument formats:

```
go run . [STRING]
go run . [STRING] [BANNER]
go run . --output=<file> [STRING] [BANNER]
go run . --color=<color> [STRING] [BANNER]
go run . --color=<color> [SUBSTRING] [STRING] [BANNER]
go run . --output=<file> --color=<color> [STRING] [BANNER]
go run . --output=<file> --color=<color> [SUBSTRING] [STRING] [BANNER]
```

## Implementation Details

Banners map each printable character (ASCII 32–126) to 8 ASCII art lines:

```go
map[string][]string
```

Rendering algorithm:

1. Iterate each input character
2. Retrieve banner representation
3. Concatenate rows for all characters (8 rows total)
4. Repeat for multiline input

Time complexity: O(n × h), where n = number of characters, h = banner height (8).

Color formatting is applied **after rendering**, ensuring that rendering logic remains pure and independent of presentation concerns.

The `--output` flag uses `filepath.Base` to restrict output to the current directory, preventing path traversal attacks (e.g. `--output=../../etc/passwd` is rejected).

## Testing (TDD)

Development uses TDD: write failing tests, implement minimal logic, refactor while keeping tests passing.
Test coverage includes banner loading, tokenization, rendering, multiline input, color formatting, output file writing, and error handling. Independent testing ensures predictable behavior and safe refactoring.

Run all tests:

```
go test ./tests/...
```

## Error Handling

The program returns meaningful errors for:
- Missing arguments or wrong flag format (prints usage message)
- Invalid banner name or missing banner file
- Unsupported or invalid input characters
- File creation failures

`main.go` acts only as the entry point and delegates all logic to `pipeline.Run`.

## Design Principles

Single Responsibility Principle, Separation of Concerns, Modular structure, Predictable rendering, Maintainable and extensible codebase.
