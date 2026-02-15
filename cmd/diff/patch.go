package main

import (
	"flag"
	"fmt"
	"os"

	app "github.com/arran4/golang-diff"
)

type PatchCmd struct {
	*RootCmd
	Flags     *flag.FlagSet
	patchFile string
	targetDir string
}

func (c *RootCmd) NewPatch() Cmd {
	fs := flag.NewFlagSet("patch", flag.ContinueOnError)
	cPatch := &PatchCmd{
		RootCmd: c,
		Flags:   fs,
	}
	return cPatch
}

func (c *PatchCmd) Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s patch <patchFile> <targetDir>:\n", os.Args[0])
	c.Flags.PrintDefaults()
}

func (c *PatchCmd) Execute(args []string) error {
	if err := c.Flags.Parse(args); err != nil {
		return err
	}
	remaining := c.Flags.Args()
	if len(remaining) < 2 {
		c.Usage()
		return fmt.Errorf("expected at least 2 arguments, got %d", len(remaining))
	}
	c.patchFile = remaining[0]
	c.targetDir = remaining[1]

	app.PatchFiles(c.patchFile, c.targetDir)
	return nil
}
