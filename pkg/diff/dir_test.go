package diff

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/tools/txtar"
)

func TestDirectoryDiff(t *testing.T) {
	// Define test data using txtar format
	archiveData := `
-- dir1/file1.txt --
content1
-- dir1/file2.txt --
content2
-- dir2/file1.txt --
content1
-- dir2/file2.txt --
content2 modified
-- dir2/file3.txt --
new file
`
	ar := txtar.Parse([]byte(archiveData))
	tempDir := t.TempDir()

	// Unpack archive to tempDir
	for _, f := range ar.Files {
		path := filepath.Join(tempDir, f.Name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, f.Data, 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Change CWD to tempDir to test relative paths
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(cwd); err != nil {
			t.Errorf("Failed to change directory back to %s: %v", cwd, err)
		}
	}()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}

	// 1. Test Diff
	output, err := Diff("dir1", "dir2")
	if err != nil {
		t.Fatalf("Diff failed: %v", err)
	}

	// Check file1.txt (identical)
	if !strings.Contains(output, "content1") {
		t.Error("Missing content1 in diff output")
	}

	// Check file2.txt (modified)
	if !strings.Contains(output, "content2 modified") {
		t.Error("Missing modified content in diff output")
	}

	// Check file3.txt (new file)
	if !strings.Contains(output, "new file") {
		t.Error("Missing new file content in diff output")
	}

	// 2. Test Patch
	// Modify dir1/file2.txt to junk
	if err := os.WriteFile(filepath.Join("dir1", "file2.txt"), []byte("junk"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Apply patch to "." (current dir)
	if err := Apply(output, "."); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	// Check dir1/file2.txt
	c2, _ := os.ReadFile(filepath.Join("dir1", "file2.txt"))
	// Patch restores file with trailing newline if original had one (or if Apply adds it).
	// walker.go does not strip newlines from content.
	// patch.go reconstructs file by Joining lines.
	// If lines were "content", Join adds nothing?
	// Wait, processBlock: data := strings.Join(content, "\n").
	// If content is ["line1"], data is "line1". No trailing newline.
	// If original file had "line1\n", Split gives ["line1", ""].
	// processBlock sees ["line1", ""]. Join gives "line1\n".
	// So newlines are preserved.
	// In txtar, "content2 modified" implies "content2 modified\n".
	if string(c2) != "content2 modified\n" {
		t.Errorf("Patch failed for file2.txt. Got %q", string(c2))
	}

	// Check dir1/file3.txt (new file)
	c3, err := os.ReadFile(filepath.Join("dir1", "file3.txt"))
	if err != nil {
		t.Errorf("Patch failed: file3.txt missing")
	} else if string(c3) != "new file\n" {
		t.Errorf("Patch failed for file3.txt. Got %q", string(c3))
	}
}

func TestFileFilter(t *testing.T) {
	archiveData := `
-- dir1/file.go --
package main
-- dir1/readme.md --
# Readme
-- dir2/file.go --
package main
-- dir2/readme.md --
# Readme modified
`
	ar := txtar.Parse([]byte(archiveData))
	tempDir := t.TempDir()
	for _, f := range ar.Files {
		path := filepath.Join(tempDir, f.Name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(path, f.Data, 0644); err != nil {
			t.Fatalf("Failed to write file: %v", err)
		}
	}

	dir1 := filepath.Join(tempDir, "dir1")
	dir2 := filepath.Join(tempDir, "dir2")

	opts := []interface{}{
		FileFilter(func(path string) bool {
			matched, _ := filepath.Match("*.go", filepath.Base(path))
			return matched
		}),
	}

	output, err := Diff(dir1, dir2, opts...)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(output, "file.go") {
		t.Error("file.go missing")
	}
	if strings.Contains(output, "readme.md") {
		t.Error("readme.md should be filtered out")
	}
}
