package app

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/arran4/golang-diff/pkg/diff"
)

// CompareFiles is a subcommand `diff compare`
// Compares two files side by side.
//
// Flags:
//
//	file1: @1 File 1 path
//	file2: @2 File 2 path
//	term: --term -t Terminal mode (colors)
//	interactive: --interactive -i Interactive mode
//	searchDepth: --search-depth -s (default: 1000) Max lines to search for alignment
//	limitLines: --max-lines (default: 0) Max lines to compare
//	limitWidth: --max-width (default: 0) Max width
//	linesSelection: --lines -l Line selection
//	widthSelection: --width -w Width selection
func CompareFiles(file1 string, file2 string, term bool, interactive bool, searchDepth int, limitLines, limitWidth int, linesSelection, widthSelection string) {
	c1, err := os.ReadFile(file1)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", file1, err)
		return
	}
	c2, err := os.ReadFile(file2)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", file2, err)
		return
	}

	opts := []interface{}{
		diff.TermMode(term),
		diff.Interactive(interactive),
		diff.SearchDepth(searchDepth),
		diff.WithMaxLines(limitLines),
		diff.WithMaxWidth(limitWidth),
		diff.WithLineSelectionShortCode(linesSelection),
		diff.WithWidthSelectionShortCode(widthSelection),
	}

	output := diff.Compare(string(c1), string(c2), opts...)

	if interactive {
		// Use less for paging
		cmd := exec.Command("less", "-R")
		cmd.Stdin = strings.NewReader(output)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			// Fallback if less fails or not found
			fmt.Print(output)
		}
	} else {
		fmt.Print(output)
	}
}
