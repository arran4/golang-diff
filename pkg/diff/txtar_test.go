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

type decoder func(string) ([]byte, error)

var decoders = map[string]decoder{
	".gostr": func(s string) ([]byte, error) {
		s, err := strconv.Unquote(strings.TrimSpace(s))
		if err != nil {
			return nil, err
		}
		return []byte(s), nil
	},
}

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
				name := f.Name
				ext := filepath.Ext(name)
				base := strings.TrimSuffix(name, ext)

				// Handle decoding if an extension matches a registered decoder
				var data []byte
				if decode, ok := decoders[ext]; ok {
					var err error
					data, err = decode(string(f.Data))
					if err != nil {
						t.Fatalf("Failed to decode %s: %v", name, err)
					}
					// Update base name to remove the decoded extension (e.g. input1.txt.gostr -> input1.txt)
					// But wait, base is input1.txt if ext is .gostr
					// Then we might have another extension .txt
				} else {
					data = f.Data
					// If no decoder, base is the name without the last extension.
					// e.g. input1.txt -> base input1
					// e.g. input1.txt.gostr -> ext .gostr -> base input1.txt
					// So if we decoded, we want to match against "input1.txt" etc.
					// If we didn't decode, we want to match against "input1.txt" etc.
				}

				// Re-normalize name for matching
				// If decoded, name effectively becomes the base name (e.g. input1.txt)
				// If not decoded, use original name.
				matchName := name
				if _, ok := decoders[ext]; ok {
					matchName = base
				}

				switch matchName {
				case "input1.txt":
					input1 = data
				case "input2.txt":
					input2 = data
				case "expected.txt":
					expected = data
				case "options.json":
					optionsJSON = data
				case "documentation.md":
					doc = data
				}
			}

			if doc == nil {
				t.Error("Missing documentation.md")
			}

			if input1 == nil || input2 == nil || expected == nil {
				t.Fatalf("Missing required files (input1.txt, input2.txt, expected.txt) or their decoded variants")
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
				t.Logf("Diff of Expected vs Got:\n%s", Compare(expectedStr, got))
			}
		})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
