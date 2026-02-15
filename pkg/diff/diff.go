package diff

import (
	"fmt"
	"strings"
)

func Compare(a, b interface{}, options ...interface{}) string {
	opts := NewOptions(options...)

	aLines := toStringSlice(a)
	bLines := toStringSlice(b)

	if opts.LinesSelection != nil {
		compiled, err := opts.LinesSelection.Compile()
		if err != nil {
			return fmt.Sprintf("Error compiling line selection: %v", err)
		}
		aLines = applyLineSelection(aLines, compiled.RangesA)
		bLines = applyLineSelection(bLines, compiled.RangesB)
	}

	if opts.LimitLines > 0 {
		if len(aLines) > opts.LimitLines {
			aLines = aLines[:opts.LimitLines]
		}
		if len(bLines) > opts.LimitLines {
			bLines = bLines[:opts.LimitLines]
		}
	}

	if opts.WidthSelection != nil {
		compiled, err := opts.WidthSelection.Compile()
		if err != nil {
			return fmt.Sprintf("Error compiling width selection: %v", err)
		}
		aLines = applyWidthSelection(aLines, compiled.RangesA)
		bLines = applyWidthSelection(bLines, compiled.RangesB)
	}

	if opts.LimitWidth > 0 {
		aLines = applyLimitWidth(aLines, opts.LimitWidth)
		bLines = applyLimitWidth(bLines, opts.LimitWidth)
	}

	diffs := AlignLines(aLines, bLines, opts)
	return FormatDiff(diffs, opts)
}

func applyLineSelection(lines []string, ranges []Range) []string {
	if len(ranges) == 0 {
		return lines
	}
	var result []string
	for _, r := range ranges {
		start := r.Start - 1
		end := r.End
		if start < 0 {
			start = 0
		}
		if start >= len(lines) {
			continue
		}
		if end > len(lines) {
			end = len(lines)
		}
		if start < end {
			result = append(result, lines[start:end]...)
		}
	}
	return result
}

func applyWidthSelection(lines []string, ranges []Range) []string {
	if len(ranges) == 0 {
		return lines
	}
	var result []string
	for _, line := range lines {
		var sb strings.Builder
		runes := []rune(line)
		for _, r := range ranges {
			start := r.Start - 1
			end := r.End
			if start < 0 {
				start = 0
			}
			if start >= len(runes) {
				continue
			}
			if end > len(runes) {
				end = len(runes)
			}
			if start < end {
				sb.WriteString(string(runes[start:end]))
			}
		}
		result = append(result, sb.String())
	}
	return result
}

func applyLimitWidth(lines []string, limit int) []string {
	var result []string
	for _, line := range lines {
		runes := []rune(line)
		if len(runes) > limit {
			result = append(result, string(runes[:limit]))
		} else {
			result = append(result, line)
		}
	}
	return result
}

func toStringSlice(v interface{}) []string {
	switch t := v.(type) {
	case string:
		return strings.Split(t, "\n")
	case []string:
		return t
	case []byte:
		return strings.Split(string(t), "\n")
	default:
		return []string{}
	}
}
