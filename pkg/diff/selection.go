package diff

import (
	"fmt"
	"strconv"
	"strings"
)

type Range struct {
	Start int
	End   int
}

type ShortCodeLineSelectionLang struct {
	RangesA []Range
	RangesB []Range
}

// Ensure ShortCodeLineSelectionLang implements fmt.Stringer
var _ fmt.Stringer = (*ShortCodeLineSelectionLang)(nil)

func (s *ShortCodeLineSelectionLang) String() string {
	var parts []string
	if len(s.RangesA) > 0 {
		var sub []string
		for _, r := range s.RangesA {
			sub = append(sub, fmt.Sprintf("%d-%d", r.Start, r.End))
		}
		parts = append(parts, "a:"+strings.Join(sub, ","))
	}
	if len(s.RangesB) > 0 {
		var sub []string
		for _, r := range s.RangesB {
			sub = append(sub, fmt.Sprintf("%d-%d", r.Start, r.End))
		}
		parts = append(parts, "b:"+strings.Join(sub, ","))
	}
	return strings.Join(parts, ",")
}

type ShortCodeWidthSelectionLang struct {
	RangesA []Range
	RangesB []Range
}

// Ensure ShortCodeWidthSelectionLang implements fmt.Stringer
var _ fmt.Stringer = (*ShortCodeWidthSelectionLang)(nil)

func (s *ShortCodeWidthSelectionLang) String() string {
	var parts []string
	if len(s.RangesA) > 0 {
		var sub []string
		for _, r := range s.RangesA {
			sub = append(sub, fmt.Sprintf("%d-%d", r.Start, r.End))
		}
		parts = append(parts, "a:"+strings.Join(sub, ","))
	}
	if len(s.RangesB) > 0 {
		var sub []string
		for _, r := range s.RangesB {
			sub = append(sub, fmt.Sprintf("%d-%d", r.Start, r.End))
		}
		parts = append(parts, "b:"+strings.Join(sub, ","))
	}
	return strings.Join(parts, ",")
}

// Compile parses the line selection string.
// Grammar: [a:]start[-end],[b:]start[-end],...
// Or positional: start[-end],start[-end] -> A, B
func (s WithLineSelectionShortCode) Compile() (*ShortCodeLineSelectionLang, error) {
	str := string(s)
	if str == "" {
		return &ShortCodeLineSelectionLang{}, nil
	}
	parts := strings.Split(str, ",")
	res := &ShortCodeLineSelectionLang{}
	var posArgs []string

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if strings.HasPrefix(p, "a:") {
			r, err := parseRange(strings.TrimPrefix(p, "a:"))
			if err != nil {
				return nil, err
			}
			res.RangesA = append(res.RangesA, r)
		} else if strings.HasPrefix(p, "b:") {
			r, err := parseRange(strings.TrimPrefix(p, "b:"))
			if err != nil {
				return nil, err
			}
			res.RangesB = append(res.RangesB, r)
		} else {
			posArgs = append(posArgs, p)
		}
	}

	for i, arg := range posArgs {
		r, err := parseRange(arg)
		if err != nil {
			return nil, err
		}
		if i == 0 {
			res.RangesA = append(res.RangesA, r)
			// If only one positional arg, apply to B as well?
			if len(posArgs) == 1 {
				res.RangesB = append(res.RangesB, r)
			}
		} else if i == 1 {
			res.RangesB = append(res.RangesB, r)
		}
	}

	return res, nil
}

func parseRange(s string) (Range, error) {
	parts := strings.Split(s, "-")
	if len(parts) == 1 {
		v, err := strconv.Atoi(parts[0])
		if err != nil {
			return Range{}, fmt.Errorf("invalid range: %s", s)
		}
		return Range{Start: v, End: v}, nil
	}
	if len(parts) == 2 {
		start, err1 := strconv.Atoi(parts[0])
		end, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			return Range{}, fmt.Errorf("invalid range: %s", s)
		}
		if start > end {
			return Range{}, fmt.Errorf("invalid range: start > end (%d > %d)", start, end)
		}
		return Range{Start: start, End: end}, nil
	}
	return Range{}, fmt.Errorf("invalid range format: %s", s)
}

func (s WithWidthSelectionShortCode) Compile() (*ShortCodeWidthSelectionLang, error) {
	str := string(s)
	if str == "" {
		return &ShortCodeWidthSelectionLang{}, nil
	}
	parts := strings.Split(str, ",")
	res := &ShortCodeWidthSelectionLang{}
	var posArgs []string

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if strings.HasPrefix(p, "a:") {
			r, err := parseRange(strings.TrimPrefix(p, "a:"))
			if err != nil {
				return nil, err
			}
			res.RangesA = append(res.RangesA, r)
		} else if strings.HasPrefix(p, "b:") {
			r, err := parseRange(strings.TrimPrefix(p, "b:"))
			if err != nil {
				return nil, err
			}
			res.RangesB = append(res.RangesB, r)
		} else {
			posArgs = append(posArgs, p)
		}
	}

	for i, arg := range posArgs {
		r, err := parseRange(arg)
		if err != nil {
			return nil, err
		}
		if i == 0 {
			res.RangesA = append(res.RangesA, r)
			if len(posArgs) == 1 {
				res.RangesB = append(res.RangesB, r)
			}
		} else if i == 1 {
			res.RangesB = append(res.RangesB, r)
		}
	}

	return res, nil
}
