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

func TestCompareOptions(t *testing.T) {
	a := "line1\nline2\nline3\nline4"
	b := "line1\nline2\nline3 diff\nline4"

	// Limit Lines
	// Compare only first 2 lines. Should be identical.
	// Output should NOT contain "line3 diff".
	res := Compare(a, b, WithMaxLines(2))
	if strings.Contains(res, "line3 diff") {
		t.Errorf("Expected line3 diff to be ignored with MaxLines(2). Output:\n%s", res)
	}
	if !strings.Contains(res, "line2") {
		t.Errorf("Expected line2 to be present. Output:\n%s", res)
	}

	// Line Selection
	// Select line 3 only.
	// Output should contain "line3 diff".
	res = Compare(a, b, WithLineSelectionShortCode("3"))
	if !strings.Contains(res, "line3 diff") {
		t.Errorf("Expected line3 diff to be present with LineSelection(3). Output:\n%s", res)
	}
	if strings.Contains(res, "line1") {
		t.Errorf("Expected line1 to be ignored with LineSelection(3). Output:\n%s", res)
	}

	// Width Selection (Columns)
	// Select columns 1-4.
	// "line1" -> "line"
	// "line3 diff" -> "line"
	// So they become identical.
	res = Compare(a, b, WithWidthSelectionShortCode("1-4"))
	if strings.Contains(res, "diff") { // "diff" word in content should be gone
		t.Errorf("Expected 'diff' content to be stripped by WidthSelection(1-4). Output:\n%s", res)
	}

	// Limit Width
	// Limit width to 5.
	// "line3 diff" -> "line3" (len 5)
	// "line3" -> "line3"
	// So line3 matches line3.
	res = Compare(a, b, WithMaxWidth(5))
	if strings.Contains(res, " diff") {
		t.Errorf("Expected ' diff' suffix to be stripped by MaxWidth(5). Output:\n%s", res)
	}
}
