package diff

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/tools/txtar"
)

type DirTest struct {
	Name     string `json:"name"`
	Type     string `json:"type"` // "diff", "patch" (implies diff)
	Dir1     string `json:"dir1"`
	Dir2     string `json:"dir2"`
	Filter   string `json:"filter"`
	Expected string `json:"expected"`
}

func TestDirTxtar(t *testing.T) {
	// Walk the testdata directory from the embed.FS
	err := fs.WalkDir(testData, "testdata", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".txtar") || !strings.Contains(path, "dir_") {
			return nil
		}

		t.Run(filepath.Base(path), func(t *testing.T) {
			content, err := testData.ReadFile(path)
			if err != nil {
				t.Fatalf("Failed to read file: %v", err)
			}
			ar := txtar.Parse(content)

			// Extract test-config.json
			var tests []DirTest
			var testsData []byte
			for _, f := range ar.Files {
				if f.Name == "test-config.json" {
					testsData = f.Data
					break
				}
			}
			if len(testsData) > 0 {
				if err := json.Unmarshal(testsData, &tests); err != nil {
					t.Fatalf("Failed to parse test-config.json: %v", err)
				}
			} else {
				t.Fatalf("No test-config.json found in %s", path)
			}

			// Setup Temp Dir with file contents
			tempDir := t.TempDir()
			for _, f := range ar.Files {
				if f.Name == "test-config.json" || f.Name == "expected.diff" {
					continue
				}
				p := filepath.Join(tempDir, f.Name)
				if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
					t.Fatalf("MkdirAll failed: %v", err)
				}
				if err := os.WriteFile(p, f.Data, 0644); err != nil {
					t.Fatalf("WriteFile failed: %v", err)
				}
			}

			// Get expected diff content if available in archive (default fallback)
			var globalExpected string
			for _, f := range ar.Files {
				if f.Name == "expected.diff" {
					globalExpected = string(f.Data)
					break
				}
			}

			// Change CWD to tempDir to test relative paths
			cwd, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := os.Chdir(cwd); err != nil {
					t.Errorf("Failed to change directory back to %s: %v", cwd, err)
				}
			}()
			if err := os.Chdir(tempDir); err != nil {
				t.Fatal(err)
			}

			for _, test := range tests {
				t.Run(test.Name, func(t *testing.T) {
					dir1 := test.Dir1
					dir2 := test.Dir2

					var opts []interface{}
					if test.Filter != "" {
						opts = append(opts, FileFilter(func(path string) bool {
							matched, _ := filepath.Match(test.Filter, filepath.Base(path))
							return matched
						}))
					}

					output, err := Diff(dir1, dir2, opts...)
					if err != nil {
						t.Fatalf("Diff failed: %v", err)
					}

					// Determine expected output
					expected := globalExpected
					if test.Expected != "" {
						// Look for specific expected file in archive
						found := false
						for _, f := range ar.Files {
							if f.Name == test.Expected {
								expected = string(f.Data)
								found = true
								break
							}
						}
						if !found && test.Expected != "expected.diff" {
							t.Fatalf("Expected file %s not found in archive", test.Expected)
						}
					}

					// Normalize line endings and trim spaces for comparison
					gotTrimmed := trimTrailingSpaces(output)
					expectedTrimmed := trimTrailingSpaces(expected)

					if gotTrimmed != expectedTrimmed {
						t.Errorf("Mismatch:\nExpected:\n%q\nGot:\n%q\nDiff:\n%s", expectedTrimmed, gotTrimmed, Compare(expectedTrimmed, gotTrimmed))
					}

					if strings.Contains(test.Type, "patch") {
						// Apply patch to current directory (tempDir)
						if err := Apply(output, "."); err != nil {
							t.Fatalf("Apply failed: %v", err)
						}
						// Verify (very basic verify: run diff again, should be empty or identical?)
						// Wait, Applying patch makes dir1 look like dir2?
						// "Diff dir1 dir2" shows how to change dir1 to match dir2.
						// Applying it to dir1 should make dir1 == dir2.
						// So Diff(dir1, dir2) should return equality everywhere?
						// "Equal lines" are part of the output.
						// If identical, Diff output shows equality.
						// Let's verify specific file content if possible, or just re-diff.
						// But re-diffing involves walking.
						// Let's rely on standard verification logic if we had one.
						// For now, let's trust Apply didn't error.
						// We can check if file2.txt in dir1 is modified.

						// In "basic" test:
						// dir1/file2.txt: "content2" -> "content2 modified"
						// dir1/file3.txt: (missing) -> "new file"

						if test.Name == "basic" {
							// Check relative paths because we chdired
							c2, _ := os.ReadFile("dir1/file2.txt")
							if strings.TrimSpace(string(c2)) != "content2 modified" {
								t.Errorf("Patch failed for file2.txt. Got %q", string(c2))
							}
							c3, err := os.ReadFile("dir1/file3.txt")
							if err != nil {
								t.Errorf("Patch failed: file3.txt missing")
							} else if strings.TrimSpace(string(c3)) != "new file" {
								t.Errorf("Patch failed for file3.txt. Got %q", string(c3))
							}
						}
					}
				})
			}
		})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
