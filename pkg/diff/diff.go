package diff

import (
	"strings"
)

func Compare(a, b interface{}, options ...interface{}) string {
	opts := NewOptions(options...)

	aLines := toStringSlice(a)
	bLines := toStringSlice(b)

	diffs := AlignLines(aLines, bLines, opts)
	output := FormatDiff(diffs, opts)
	if opts.TestingT != nil {
		opts.TestingT.Helper()
		for _, diff := range diffs {
			if diff.Type != DiffEqual {
				opts.TestingT.Errorf("%s", output)
				break
			}
		}
	}
	return output
}

func Diff(a, b interface{}, options ...interface{}) string {
	return Compare(a, b, options...)
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
