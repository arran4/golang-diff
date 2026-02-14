package diff

import (
	"strings"
)

func Compare(a, b interface{}, options ...interface{}) string {
	opts := NewOptions(options...)

	aLines := toStringSlice(a)
	bLines := toStringSlice(b)

	diffs := AlignLines(aLines, bLines, opts)
	return FormatDiff(diffs, opts)
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
