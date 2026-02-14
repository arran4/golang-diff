# golang-diff

`golang-diff` is a Go library and CLI tool for comparing text files line by line and character by character. It provides colored terminal output and an interactive viewer.

## Features

- **Line Alignment**: Uses a hash-based lookahead algorithm to align lines efficiently.
- **Character-Level Diff**: Computes precise character-level differences within modified lines.
- **Terminal Colors**: Highlighting for additions (green) and deletions (red).
- **Interactive Mode**: View differences in `less` (or your preferred pager) for easy navigation.
- **Configurable Lookahead**: Adjust the search depth for line alignment.

## Installation

To install the CLI tool:

```bash
go install github.com/arran4/golang-diff/cmd/diff@latest
```

This will install the `diff` binary to your `$GOPATH/bin` directory.

## Usage (CLI)

The basic usage is:

```bash
diff compare <file1> <file2> [flags]
```

### Flags

- `--term` / `-t`: Enable terminal colors (default: false).
- `--interactive` / `-i`: Open the diff in an interactive pager (`less`) (default: false).
- `--max-lines` / `-m`: Set the maximum number of lines to search ahead for alignment (default: 1000).

### Examples

Compare two files with color output:

```bash
diff compare file1.txt file2.txt -t
```

Compare two files interactively:

```bash
diff compare file1.txt file2.txt -i
```

## Output Guide

The output displays the two files side-by-side. The center column contains a symbol indicating the type of difference between the lines.

### Symbols

| Symbol | Meaning | Color (Term Mode) |
| :--- | :--- | :--- |
| `==` | Lines are identical | Green |
| `1d` | One continuous difference block | Red |
| `2d` | Two difference blocks | Red |
| `d` | Multiple difference blocks | Red |
| `w` | Whitespace difference only | Yellow |
| `q` | Mixed character and whitespace difference | Red |
| `$` | End of Line (EOL) difference (e.g., CRLF vs LF) | Yellow |

### Colors

When terminal mode is enabled (`-t`):
- **Green**: Added or identical content.
- **Red**: Deleted or modified content.
- **Yellow**: Whitespace or EOL differences.

### Example Output

```text
Hello world    ==  Hello world
This is a test 1d  This is a test case
```

In this example:
- The first line is identical (`==`).
- The second line has one difference (`1d`): " case" is added on the right side.

## Usage (Library)

You can also use `golang-diff` as a library in your Go projects.

### Import

```go
import "github.com/arran4/golang-diff/pkg/diff"
```

### Example

```go
package main

import (
	"fmt"
	"github.com/arran4/golang-diff/pkg/diff"
)

func main() {
	text1 := "Hello world\nThis is a test"
	text2 := "Hello World\nThis is a test case"

	// Compare strings with options
	output := diff.Compare(text1, text2, diff.TermMode(true))

	fmt.Println(output)
}
```

## Development

To run the tests:

```bash
go test ./...
```
