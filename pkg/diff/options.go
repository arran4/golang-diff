package diff

type TermMode bool
type Interactive bool
type SearchDepth int
type WithSearchDepth int
type WithMaxLines int
type WithMaxWidth int

type LineSelection interface {
	Compile() (*ShortCodeLineSelectionLang, error)
}
type WidthSelection interface {
	Compile() (*ShortCodeWidthSelectionLang, error)
}
type WithLineSelectionShortCode string
type WithWidthSelectionShortCode string

type LineUpFunc func(a, b []string, opts *Options) []DiffLine

type Options struct {
	TermMode       bool
	Interactive    bool
	SearchDepth    int
	LimitLines     int
	LimitWidth     int
	LinesSelection LineSelection
	WidthSelection WidthSelection
	LineUpFunc     LineUpFunc
}

type DiffType string

const (
	DiffEqual DiffType = "=="
	Diff1     DiffType = "1d" // One continuous difference
	Diff2     DiffType = "2d" // Two differences
	DiffChar  DiffType = "d"  // Generic difference (3+)
	DiffSpace DiffType = "w"  // Whitespace only
	DiffMixed DiffType = "q"  // Character and whitespace
	DiffEOL   DiffType = "$"  // EOL difference
)

type OpType int

const (
	OpMatch OpType = iota
	OpInsert
	OpDelete
)

type Operation struct {
	Type    OpType
	Content string
}

type DiffLine struct {
	Left  string
	Right string
	Type  DiffType
	Ops   []Operation
}

func NewOptions(args ...interface{}) *Options {
	opts := &Options{
		SearchDepth: 1000,
	}
	for _, arg := range args {
		switch v := arg.(type) {
		case TermMode:
			opts.TermMode = bool(v)
		case Interactive:
			opts.Interactive = bool(v)
		case bool:
			// Heuristic: if boolean true is passed, maybe default to TermMode?
			// But type switch is safer.
		case SearchDepth:
			opts.SearchDepth = int(v)
		case WithSearchDepth:
			opts.SearchDepth = int(v)
		case int:
			opts.SearchDepth = v
		case WithMaxLines:
			opts.LimitLines = int(v)
		case WithMaxWidth:
			opts.LimitWidth = int(v)
		case LineSelection:
			opts.LinesSelection = v
		case WidthSelection:
			opts.WidthSelection = v
		case LineUpFunc:
			opts.LineUpFunc = v
		}
	}
	return opts
}
