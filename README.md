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
