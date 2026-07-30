package diff

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func Diff(path1, path2 string, options ...interface{}) (string, error) {
	opts := NewOptions(options...)

	// Initial check to handle file vs dir at root level
	_, err1 := os.Stat(path1)
	_, err2 := os.Stat(path2)

	// If errors are not IsNotExist, return them
	if err1 != nil && !os.IsNotExist(err1) {
		return "", err1
	}
	if err2 != nil && !os.IsNotExist(err2) {
		return "", err2
	}

	return walk("", path1, path2, opts)
}

func walk(relPath, root1, root2 string, opts *Options) (string, error) {
	path1 := filepath.Join(root1, relPath)
	path2 := filepath.Join(root2, relPath)

	fi1, err1 := os.Stat(path1)
	fi2, err2 := os.Stat(path2)

	exists1 := err1 == nil
	exists2 := err2 == nil

	if !exists1 && !exists2 {
		return "", nil
	}

	// Determine if directories
	isDir1 := exists1 && fi1.IsDir()
	isDir2 := exists2 && fi2.IsDir()

	// Handle type mismatch (File vs Dir)
	if exists1 && exists2 && isDir1 != isDir2 {
		return fmt.Sprintf("File %s is a directory while file %s is a regular file\n", path1, path2), nil
	}

	isDir := isDir1 || isDir2

	// If it is a file (or one is a file and other missing)
	if !isDir {
		// Apply filter
		// We use relPath for filtering. If relPath is empty (root file), use the base name?
		// But usually filter works on the path being walked.
		// If Diff("file1", "file2") is called, relPath is empty.
		// If opts.FileFilter is provided, what should it check?
		// Probably the relative path is appropriate.
		// If relPath is empty, maybe we shouldn't filter? Or check "."?
		// Let's assume FileFilter handles empty string or we check the name.
		checkPath := relPath

		if checkPath != "" && opts.FileFilter != nil && !opts.FileFilter(checkPath) {
			return "", nil
		}

		c1 := ""
		c2 := ""
		if exists1 {
			b, err := os.ReadFile(path1)
			if err != nil {
				return "", err
			}
			c1 = string(b)
		}
		if exists2 {
			b, err := os.ReadFile(path2)
			if err != nil {
				return "", err
			}
			c2 = string(b)
		}

		// reuse Compare logic
		lines1 := strings.Split(c1, "\n")
		lines2 := strings.Split(c2, "\n")
		diffs := AlignLines(lines1, lines2, opts)
		output := FormatDiff(diffs, opts)

		header := fmt.Sprintf("Diff %q %q\n", path1, path2)
		return header + output, nil
	}

	// Directory recursion
	entries := make(map[string]struct{})
	if exists1 {
		des, err := os.ReadDir(path1)
		if err != nil {
			return "", err
		}
		for _, de := range des {
			entries[de.Name()] = struct{}{}
		}
	}
	if exists2 {
		des, err := os.ReadDir(path2)
		if err != nil {
			return "", err
		}
		for _, de := range des {
			entries[de.Name()] = struct{}{}
		}
	}

	var names []string
	for name := range entries {
		names = append(names, name)
	}
	sort.Strings(names)

	var sb strings.Builder
	for _, name := range names {
		childRel := filepath.Join(relPath, name)
		res, err := walk(childRel, root1, root2, opts)
		if err != nil {
			return "", err
		}
		sb.WriteString(res)
	}

	return sb.String(), nil
}
