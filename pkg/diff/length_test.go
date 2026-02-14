package diff

import (
	"testing"
)

func TestAlignLines_DifferentLengths(t *testing.T) {
	opts := NewOptions()

	tests := []struct {
		name string
		a    []string
		b    []string
		want int // Expected number of aligned lines
	}{
		{
			name: "b longer at end",
			a:    []string{"a", "b"},
			b:    []string{"a", "b", "c", "d"},
			want: 4,
		},
		{
			name: "a longer at end",
			a:    []string{"a", "b", "c", "d"},
			b:    []string{"a", "b"},
			want: 4,
		},
		{
			name: "b longer at start",
			a:    []string{"c", "d"},
			b:    []string{"a", "b", "c", "d"},
			want: 4,
		},
		{
			name: "a longer at start",
			a:    []string{"a", "b", "c", "d"},
			b:    []string{"c", "d"},
			want: 4,
		},
		{
			name: "mixed insertion/deletion",
			a:    []string{"a", "c", "e"},
			b:    []string{"a", "b", "c", "d", "e"},
			want: 5,
		},
		{
			name: "completely different",
			a:    []string{"a", "b"},
			b:    []string{"c", "d", "e"},
			want: 3,
		},
		{
			name: "empty a",
			a:    []string{},
			b:    []string{"a", "b"},
			want: 2,
		},
		{
			name: "empty b",
			a:    []string{"a", "b"},
			b:    []string{},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AlignLines(tt.a, tt.b, opts)
			if len(got) != tt.want {
				t.Errorf("AlignLines(%v, %v) returned %d lines, want %d", tt.a, tt.b, len(got), tt.want)
				for i, l := range got {
					t.Logf("%d: %s | %s (%v)", i, l.Left, l.Right, l.Type)
				}
			}
		})
	}
}

func TestAlignLines_Lookahead(t *testing.T) {
	// Restrict lookahead to 1
	opts := NewOptions()
	opts.MaxLines = 1

	a := []string{"a", "x", "y", "z", "b"}
	b := []string{"a", "b"}
	// "b" is at index 4 in a, index 1 in b.
	// Difference is 3 lines.
	// Lookahead 1 means it will only look 1 line ahead.
	// So when at "x" (ai=1) and "b" (bi=1):
	// It will look for "x" in b at [1, 2] -> not found.
	// It will look for "b" in a at [1, 2] -> "x", "y" -> not found. "b" is at 4.
	// So it should fail to align "b" with "b" immediately.
	// It will treat "x" vs "b" as modification or separate lines.
	// Then "y", "z" vs nothing... eventually "b" vs nothing?

	got := AlignLines(a, b, opts)

	// If it aligned properly, we'd expect:
	// a | a
	// x |
	// y |
	// z |
	// b | b
	// Total 5 lines.

	// If it failed to align because of lookahead:
	// a | a (match)
	// x | b (mod)
	// y |   (del)
	// z |   (del)
	// b |   (del)
	// Total 5 lines.

	// Wait, if it treats "x" vs "b" as modification, it consumes both.
	// ai=2 (y), bi=2 (end)
	// y | (del)
	// z | (del)
	// b | (del)
	// Total 1 (match) + 1 (mod) + 3 (del) = 5 lines.

	// Let's verify what happens.
	foundAlignment := false
	for _, l := range got {
		if l.Left == "b" && l.Right == "b" && l.Type == DiffEqual {
			foundAlignment = true
			break
		}
	}

	if foundAlignment {
		t.Log("Found alignment despite short lookahead (maybe expected if logic falls through)")
	} else {
		t.Log("Did not align 'b' due to short lookahead")
	}

	// Now try with sufficient lookahead
	opts2 := NewOptions()
	opts2.MaxLines = 5
	got2 := AlignLines(a, b, opts2)
	foundAlignment2 := false
	for _, l := range got2 {
		if l.Left == "b" && l.Right == "b" && l.Type == DiffEqual {
			foundAlignment2 = true
			break
		}
	}
	if !foundAlignment2 {
		t.Error("Should have aligned 'b' with sufficient lookahead")
		for i, l := range got2 {
			t.Logf("%d: %s | %s (%v)", i, l.Left, l.Right, l.Type)
		}
	}
}
