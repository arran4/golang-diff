package diff

import (
	"fmt"
	"strings"
)

func FormatDiff(lines []DiffLine, opts *Options) string {
	maxLeft := 0
	for _, line := range lines {
		if len(line.Left) > maxLeft {
			maxLeft = len(line.Left)
		}
	}

	var sb strings.Builder
	for _, line := range lines {
		leftStr, rightStr := renderDiffLine(line, opts)

		// Padding logic: must pad based on VISIBLE length (line.Left), not ANSI length.
		visibleLeftLen := len(line.Left)
		padding := ""
		if maxLeft > visibleLeftLen {
			padding = strings.Repeat(" ", maxLeft-visibleLeftLen)
		}

		symbol := string(line.Type)
		// Buffer width: space + 2 chars + space = 4 chars
		buffer := fmt.Sprintf(" %-2s ", symbol)

		if opts.TermMode {
			buffer = colorizeSymbol(buffer, line.Type)
		}

		sb.WriteString(leftStr)
		sb.WriteString(padding)
		sb.WriteString(buffer)
		sb.WriteString(rightStr)
		sb.WriteString("\n")
	}
	return sb.String()
}

func renderDiffLine(line DiffLine, opts *Options) (string, string) {
	if !opts.TermMode {
		return line.Left, line.Right
	}

	if line.Type == DiffEqual {
		return line.Left, line.Right
	}

	// If no ops (e.g. empty lines or special case), fallback
	if len(line.Ops) == 0 {
		// Should generally have ops if ComputeDiffType was called.
		// If line.Left and line.Right are non-empty but ops empty, it means match?
		// But type != DiffEqual.
		return line.Left, line.Right
	}

	var leftSB, rightSB strings.Builder

	for _, op := range line.Ops {
		switch op.Type {
		case OpMatch:
			leftSB.WriteString(op.Content)
			rightSB.WriteString(op.Content)
		case OpDelete:
			leftSB.WriteString(colorize(op.Content, "31")) // Red
		case OpInsert:
			rightSB.WriteString(colorize(op.Content, "32")) // Green
		}
	}

	return leftSB.String(), rightSB.String()
}

func colorize(s, code string) string {
	return "\033[" + code + "m" + s + "\033[0m"
}

func colorizeSymbol(s string, t DiffType) string {
	var code string
	switch t {
	case DiffEqual:
		code = "32" // Green
	case DiffSpace, DiffEOL:
		code = "33" // Yellow
	default:
		code = "31" // Red
	}
	return "\033[" + code + "m" + s + "\033[0m"
}
