package diff

import (
	"embed"
	"encoding/json"
	"io/fs"
	"path/filepath"
	"strconv"
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
				case "input1.txt.gostr":
					s, err := strconv.Unquote(strings.TrimSpace(string(f.Data)))
					if err != nil {
						t.Fatalf("Failed to unquote input1.txt.gostr: %v", err)
					}
					input1 = []byte(s)
				case "input2.txt":
					input2 = f.Data
				case "input2.txt.gostr":
					s, err := strconv.Unquote(strings.TrimSpace(string(f.Data)))
					if err != nil {
						t.Fatalf("Failed to unquote input2.txt.gostr: %v", err)
					}
					input2 = []byte(s)
				case "expected.txt":
					expected = f.Data
				case "expected.txt.gostr":
					s, err := strconv.Unquote(strings.TrimSpace(string(f.Data)))
					if err != nil {
						t.Fatalf("Failed to unquote expected.txt.gostr: %v", err)
					}
					expected = []byte(s)
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
				t.Fatalf("Missing required files (input1.txt, input2.txt, expected.txt) or their .gostr variants")
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

			if got != expectedStr {
				t.Errorf("Mismatch for %s:\nExpected:\n%q\nGot:\n%q", path, expectedStr, got)
				// Also print a diff of the output to help debugging
				t.Logf("Diff of Expected vs Got:\n%s", Compare(expectedStr, got))
			}
		})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
