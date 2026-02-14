package app

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
//	maxLines: --max-lines -m (default: 1000) Max lines to search for alignment
func CompareFiles(file1 string, file2 string, term bool, interactive bool, maxLines int) {
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
		diff.MaxLines(maxLines),
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

// DiffFiles is a subcommand 'diff diff'
// Compares two paths (files or directories) recursively.
//
// Flags:
//	path1: @1 Path 1
//	path2: @2 Path 2
//	term: --term -t Terminal mode (colors)
//	interactive: --interactive -i Interactive mode
//	maxLines: --max-lines -m (default: 1000) Max lines to search for alignment
//	selectFile: --select-file -s Glob pattern to filter files
func DiffFiles(path1, path2 string, term bool, interactive bool, maxLines int, selectFile string) {
	opts := []interface{}{
		diff.TermMode(term),
		diff.Interactive(interactive),
		diff.MaxLines(maxLines),
	}

	if selectFile != "" {
		filter := diff.FileFilter(func(path string) bool {
			// Check basename match first (standard glob behavior usually)
			matched, _ := filepath.Match(selectFile, filepath.Base(path))
			if matched {
				return true
			}
			// Check full relative path match
			matched, _ = filepath.Match(selectFile, path)
			return matched
		})
		opts = append(opts, filter)
	}

	output, err := diff.Diff(path1, path2, opts...)
	if err != nil {
		fmt.Printf("Error running diff: %v\n", err)
		return
	}

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

// PatchFiles is a subcommand 'diff patch'
// Applies a patch file to a target directory.
//
// Flags:
//	patchFile: @1 Patch file path
//	targetDir: @2 Target directory path
func PatchFiles(patchFile string, targetDir string) {
	content, err := os.ReadFile(patchFile)
	if err != nil {
		fmt.Printf("Error reading patch file %s: %v\n", patchFile, err)
		return
	}

	if err := diff.Apply(string(content), targetDir); err != nil {
		fmt.Printf("Error applying patch: %v\n", err)
		return
	}
	fmt.Println("Patch applied successfully.")
}
