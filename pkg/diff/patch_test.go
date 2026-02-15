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
	// 3d is recognized as a separator
	// Left side: 1.2.3.4
	// Right side: 1x2x3x4

	// Apply creates file2.txt?
	// Apply applies patch to targetDir.
	// The patch header says Diff "file1.txt" "file2.txt".
	// Apply logic:
	// p1 = file1.txt. target = join(targetDir, p1).
	// So it writes to file1.txt?
	// Wait, Apply writes to p1? Or p2?
	// Let's check patch.go logic.
	// "Determine target file path... currentFile = filepath.Join(targetDir, p1)"
	// So it modifies p1 (file1.txt).

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
