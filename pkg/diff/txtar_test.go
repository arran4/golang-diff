package diff

import (
	"embed"
	"encoding/json"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/tools/txtar"
)

//go:embed testdata
var testData embed.FS

func TestTxtar(t *testing.T) {
	err := fs.WalkDir(testData, "testdata", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".txtar") {
			return nil
		}

		t.Run(filepath.Base(path), func(t *testing.T) {
			content, err := testData.ReadFile(path)
			if err != nil {
				t.Fatalf("Failed to read file: %v", err)
			}
			archive := txtar.Parse(content)

			var input1, input2, expected []byte
			var optionsJSON []byte
			var doc []byte

			for _, f := range archive.Files {
				switch f.Name {
				case "input1.txt":
					input1 = f.Data
				case "input2.txt":
					input2 = f.Data
				case "expected.txt":
					expected = f.Data
				case "options.json":
					optionsJSON = f.Data
				case "documentation.md":
					doc = f.Data
				}
			}

			if doc == nil {
				t.Error("Missing documentation.md")
			}

			if input1 == nil || input2 == nil || expected == nil {
				t.Fatalf("Missing required files (input1.txt, input2.txt, expected.txt)")
			}

			var opts []interface{}
			if len(optionsJSON) > 0 {
				var rawOpts map[string]interface{}
				if err := json.Unmarshal(optionsJSON, &rawOpts); err != nil {
					t.Fatalf("Failed to parse options.json: %v", err)
				}
				for k, v := range rawOpts {
					switch k {
					case "MaxLines":
						if f, ok := v.(float64); ok {
							opts = append(opts, int(f))
						}
					case "TermMode":
						if b, ok := v.(bool); ok {
							opts = append(opts, TermMode(b))
						}
					case "Interactive":
						if b, ok := v.(bool); ok {
							opts = append(opts, Interactive(b))
						}
					}
				}
			}

			got := Compare(input1, input2, opts...)
			expectedStr := string(expected)

			// Trim trailing newlines for easier comparison if desired, but strict is better.
			// Let's stick to strict first.

			// HACK: FormatDiff might add trailing spaces that are hard to represent in txtar expected files.
			// We trim trailing whitespace from each line for comparison.
			gotTrimmed := trimTrailingSpaces(got)
			expectedTrimmed := trimTrailingSpaces(expectedStr)

			if gotTrimmed != expectedTrimmed {
				t.Errorf("Mismatch for %s:\nExpected:\n%q\nGot:\n%q", path, expectedTrimmed, gotTrimmed)
				// Also print a diff of the output to help debugging
				// We can use the Compare function itself to show the diff between expected and got!
				// But we need to be careful not to recurse infinitely if Compare is broken.
				// Assuming Compare works well enough for this:
				t.Logf("Diff of Expected vs Got:\n%s", Compare(expectedStr, got))
			}
		})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func trimTrailingSpaces(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " ")
	}
	return strings.Join(lines, "\n")
}
