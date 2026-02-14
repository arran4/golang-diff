package diff

import (
	"crypto/sha1"
	"fmt"
	"unicode"
)

func CalculateHash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))[:8]
}

func AlignLines(a, b []string, opts *Options) []DiffLine {
	if opts.LineUpFunc != nil {
		return opts.LineUpFunc(a, b, opts)
	}

	var result []DiffLine
	ai, bi := 0, 0
	maxLookahead := opts.MaxLines
	if maxLookahead <= 0 {
		maxLookahead = 1000
	}

	// Precompute hashes for quick lookup
	// We need lookups for 'a' (given b[bi]) and 'b' (given a[ai])
	// However, since we iterate forward, we can just compute on demand or precompute all.
	// Precomputing is cleaner.
	
	aHashes := make(map[string][]int)
	for i, line := range a {
		h := CalculateHash(line)
		aHashes[h] = append(aHashes[h], i)
	}
	
	bHashes := make(map[string][]int)
	for i, line := range b {
		h := CalculateHash(line)
		bHashes[h] = append(bHashes[h], i)
	}

	for ai < len(a) || bi < len(b) {
		if ai >= len(a) {
			t, ops := ComputeDiffType("", b[bi])
			result = append(result, DiffLine{Left: "", Right: b[bi], Type: t, Ops: ops})
			bi++
			continue
		}
		if bi >= len(b) {
			t, ops := ComputeDiffType(a[ai], "")
			result = append(result, DiffLine{Left: a[ai], Right: "", Type: t, Ops: ops})
			ai++
			continue
		}

		ha := CalculateHash(a[ai])
		hb := CalculateHash(b[bi])

		if ha == hb {
			// Hash match -> Align
			// Even if content differs (collision), we treat as aligned modified line.
			t, ops := ComputeDiffType(a[ai], b[bi])
			result = append(result, DiffLine{Left: a[ai], Right: b[bi], Type: t, Ops: ops})
			ai++
			bi++
			continue
		}

		// Lookahead using hashes
		// Find nearest match in b for a[ai]
		bestBj := -1
		if candidates, ok := bHashes[ha]; ok {
			for _, idx := range candidates {
				if idx > bi {
					if idx > bi+maxLookahead {
						break // Optimization: candidates are sorted? No, append order is sorted.
						// Yes, append happens in loop i=0..len. So sorted.
					}
					bestBj = idx
					break // Found nearest
				}
			}
		}

		// Find nearest match in a for b[bi]
		bestAj := -1
		if candidates, ok := aHashes[hb]; ok {
			for _, idx := range candidates {
				if idx > ai {
					if idx > ai+maxLookahead {
						break
					}
					bestAj = idx
					break
				}
			}
		}

		if bestBj != -1 && (bestAj == -1 || bestBj-bi < bestAj-ai) {
			// Insertion in b (skip b until bestBj)
			for bi < bestBj {
				t, ops := ComputeDiffType("", b[bi])
				result = append(result, DiffLine{Left: "", Right: b[bi], Type: t, Ops: ops})
				bi++
			}
		} else if bestAj != -1 {
			// Deletion in a (skip a until bestAj)
			for ai < bestAj {
				t, ops := ComputeDiffType(a[ai], "")
				result = append(result, DiffLine{Left: a[ai], Right: "", Type: t, Ops: ops})
				ai++
			}
		} else {
			// Modification (no match found nearby)
			t, ops := ComputeDiffType(a[ai], b[bi])
			result = append(result, DiffLine{Left: a[ai], Right: b[bi], Type: t, Ops: ops})
			ai++
			bi++
		}
	}
	return result
}

func ComputeDiffType(a, b string) (DiffType, []Operation) {
	if a == b {
		return DiffEqual, nil
	}
	
	if a == "" && b == "" {
	    return DiffEqual, nil
	}

	ops := getEditScript(a, b)

	blocks := 0
	inDiff := false
	hasCharDiff := false
	hasSpaceDiff := false
	isOnlyEOL := true
	hasEOLDiff := false

	for _, op := range ops {
		switch op.Type {
		case OpMatch:
			inDiff = false
		case OpInsert, OpDelete:
			if !inDiff {
				blocks++
				inDiff = true
			}
			for _, r := range op.Content {
				if unicode.IsSpace(r) {
					hasSpaceDiff = true
					if r == '\r' || r == '\n' {
						hasEOLDiff = true
					} else {
						isOnlyEOL = false
					}
				} else {
					hasCharDiff = true
					isOnlyEOL = false
				}
			}
		}
	}

	if blocks == 0 {
		return DiffEqual, ops
	}

	if hasEOLDiff && isOnlyEOL && !hasCharDiff {
		return DiffEOL, ops
	}

	if !hasCharDiff && hasSpaceDiff {
		return DiffSpace, ops
	}
	if hasCharDiff && hasSpaceDiff {
		return DiffMixed, ops
	}

	if blocks == 1 {
		return Diff1, ops
	}
	if blocks == 2 {
		return Diff2, ops
	}

	return DiffChar, ops
}

func getEditScript(s1, s2 string) []Operation {
	r1, r2 := []rune(s1), []rune(s2)
	n, m := len(r1), len(r2)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}

	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if r1[i-1] == r2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				if dp[i-1][j] >= dp[i][j-1] {
					dp[i][j] = dp[i][j-1]
				} else {
					dp[i][j] = dp[i][j-1]
				}
			}
		}
	}

	var ops []Operation
	i, j := n, m
	for i > 0 || j > 0 {
		if i > 0 && j > 0 && r1[i-1] == r2[j-1] {
			ops = append(ops, Operation{Type: OpMatch, Content: string(r1[i-1])})
			i--
			j--
		} else if j > 0 && (i == 0 || dp[i][j-1] > dp[i-1][j]) {
			ops = append(ops, Operation{Type: OpInsert, Content: string(r2[j-1])})
			j--
		} else {
			ops = append(ops, Operation{Type: OpDelete, Content: string(r1[i-1])})
			i--
		}
	}
	
	// Reverse ops
	for k := 0; k < len(ops)/2; k++ {
		ops[k], ops[len(ops)-1-k] = ops[len(ops)-1-k], ops[k]
	}
	return ops
}
