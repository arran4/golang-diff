package diff

import (
	"strings"
	"testing"
)

func TestAlignLines(t *testing.T) {
	opts := NewOptions()
	a := []string{"a", "b", "c"}
	b := []string{"a", "c"}

	diffs := AlignLines(a, b, opts)
	// Expected alignment:
	// a  | a (Equal)
	// b  |   (Delete / Diff1)
	// c  | c (Equal)

	if len(diffs) != 3 {
		t.Errorf("Expected 3 diff lines, got %d", len(diffs))
	}

	if diffs[0].Type != DiffEqual {
		t.Errorf("Line 0 should be equal, got %v", diffs[0].Type)
	}
	// Deletion (b vs "") is classified as Diff1 (one continuous difference = whole line)
	if diffs[1].Type != Diff1 {
		t.Errorf("Line 1 should be Diff1 (deletion), got %v", diffs[1].Type)
	}
	if diffs[2].Type != DiffEqual {
		t.Errorf("Line 2 should be equal, got %v", diffs[2].Type)
	}
}

func TestComputeDiffType(t *testing.T) {
	tests := []struct {
		a, b     string
		expected DiffType
	}{
		{"abc", "abc", DiffEqual},
		{"abc", "axc", Diff1},      // 1 continuous block (b->x)
		{"abcde", "azcxe", Diff2},  // 2 blocks (b->z, d->x) separated by c
		{"abc", "abd", Diff1},      // 1 block (c->d)
		{"   ", " \t ", DiffSpace}, // Whitespace only
		{"a b", "a  b", DiffSpace}, // Whitespace only (insertion of space)
		{"a", "b", Diff1},          // 1 block
		{"", "a", Diff1},           // 1 block (whole line)
		{"1.2.3.4", "1x2x3x4", "3d"},
		{"1.2.3.4.5.6.7.8.9.10.11", "1x2x3x4x5x6x7x8x9x10x11", "+d"},
	}

	for _, tt := range tests {
		got, _ := ComputeDiffType(tt.a, tt.b)
		if got != tt.expected {
			t.Errorf("ComputeDiffType(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.expected)
		}
	}
}

func TestCompareOutput(t *testing.T) {
	a := "line1\nline2\nline3"
	b := "line1\nline2 modified\nline3"

	output := Compare(a, b)
	if !strings.Contains(output, "line1") {
		t.Error("Output missing line1")
	}
	if !strings.Contains(output, "line2 modified") {
		t.Error("Output missing modified line")
	}
	// Check buffer symbol
	// line2 vs line2 modified -> " modified" inserted.
	// Contains space and characters -> DiffMixed (q)
	if !strings.Contains(output, " q  ") {
		t.Error("Output missing buffer symbol q for modified line. Output:\n" + output)
	}
}
