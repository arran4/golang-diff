package diff

type TermMode bool
type Interactive bool
type MaxLines int
type LineUpFunc func(a, b []string, opts *Options) []DiffLine
type FileFilter func(path string) bool

type Options struct {
	TermMode    bool
	Interactive bool
	MaxLines    int
	LineUpFunc  LineUpFunc
	FileFilter  FileFilter
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
		MaxLines: 1000,
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
		case MaxLines:
			opts.MaxLines = int(v)
		case int:
			opts.MaxLines = v
		case LineUpFunc:
			opts.LineUpFunc = v
		case FileFilter:
			opts.FileFilter = v
		}
	}
	return opts
}
