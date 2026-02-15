package diff

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestApplyExtended(t *testing.T) {
	tempDir := t.TempDir()

	patch := `Diff "file1.txt" "file2.txt"
1.2.3.4 3d 1x2x3x4
        ==
`
	// The "3d" separator should be recognized and parsed correctly.
	// Apply modifies "file1.txt" (the first path in the Diff header) relative to tempDir.
	err := Apply(patch, tempDir)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tempDir, "file1.txt"))
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if strings.TrimSpace(string(content)) != "1x2x3x4" {
		t.Errorf("Expected 1x2x3x4, got %q", string(content))
	}
}

func TestApplyPlusD(t *testing.T) {
	tempDir := t.TempDir()

	patch := `Diff "file_plus.txt" "file_plus.txt"
1.2.3.4.5.6.7.8.9.10 +d 1x2x3x4x5x6x7x8x9x10
                     ==
`
	if err := Apply(patch, tempDir); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tempDir, "file_plus.txt"))
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	expected := "1x2x3x4x5x6x7x8x9x10"
	if strings.TrimSpace(string(content)) != expected {
		t.Errorf("Expected %q, got %q", expected, string(content))
	}
}
