package diff

import (
	_ "embed"
	"fmt"
	"io"
)

//go:embed selection_help.txt
var selectionHelp string

func PrintSelectionHelp(w io.Writer) {
	_, _ = fmt.Fprint(w, selectionHelp)
}
