package main

import (
	"flag"
	"fmt"
	"os"

	app "github.com/arran4/golang-diff"
)

type DiffCmd struct {
	*RootCmd
	Flags       *flag.FlagSet
	path1       string
	path2       string
	term        bool
	interactive bool
	maxLines    int
	selectFile  string
}

func (c *RootCmd) NewDiff() Cmd {
	fs := flag.NewFlagSet("diff", flag.ContinueOnError)
	cDiff := &DiffCmd{
		RootCmd: c,
		Flags:   fs,
	}

	fs.BoolVar(&cDiff.term, "term", false, "Terminal mode (colors)")
	fs.BoolVar(&cDiff.term, "t", false, "Terminal mode (colors)")
	fs.BoolVar(&cDiff.interactive, "interactive", false, "Interactive mode")
	fs.BoolVar(&cDiff.interactive, "i", false, "Interactive mode")
	fs.IntVar(&cDiff.maxLines, "max-lines", 1000, "Max lines to search for alignment")
	fs.IntVar(&cDiff.maxLines, "m", 1000, "Max lines to search for alignment")
	fs.StringVar(&cDiff.selectFile, "select-file", "", "Glob pattern to filter files")
	fs.StringVar(&cDiff.selectFile, "s", "", "Glob pattern to filter files")

	return cDiff
}

func (c *DiffCmd) Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s diff <path1> <path2>:\n", os.Args[0])
	c.Flags.PrintDefaults()
}

func (c *DiffCmd) Execute(args []string) error {
	if err := c.Flags.Parse(args); err != nil {
		return err
	}
	remaining := c.Flags.Args()
	if len(remaining) < 2 {
		c.Usage()
		return fmt.Errorf("expected at least 2 arguments, got %d", len(remaining))
	}
	c.path1 = remaining[0]
	c.path2 = remaining[1]

	app.DiffFiles(c.path1, c.path2, c.term, c.interactive, c.maxLines, c.selectFile)
	return nil
}
