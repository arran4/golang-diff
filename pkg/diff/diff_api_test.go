package diff_test

import (
	"fmt"
	"testing"

	"github.com/arran4/golang-diff/pkg/diff"
)

type mockTestingT struct {
	failed   bool
	errorMsg string
}

func (m *mockTestingT) Helper() {}

func (m *mockTestingT) Errorf(format string, args ...interface{}) {
	m.failed = true
	m.errorMsg = fmt.Sprintf(format, args...)
}

func TestDiffFailOnMismatch(t *testing.T) {
	mock := &mockTestingT{}

	a := "a"
	b := "b"

	// Should fail
	diff.Diff(a, b, mock)

	if !mock.failed {
		t.Error("Expected Diff to fail with mismatching inputs")
	}
	if mock.errorMsg == "" {
		t.Error("Expected error message to be set")
	}
}

func TestDiffPassOnMatch(t *testing.T) {
	mock := &mockTestingT{}

	a := "a"
	b := "a"

	// Should pass
	diff.Diff(a, b, mock)

	if mock.failed {
		t.Error("Expected Diff to pass with matching inputs")
	}
}

func TestDiffPassOnEqual(t *testing.T) {
	mock := &mockTestingT{}
	a := []string{"a", "b"}
	b := []string{"a", "b"}
	diff.Diff(a, b, mock)
	if mock.failed {
		t.Error("Expected Diff to pass with equal inputs")
	}
}

func TestDiffFormatting(t *testing.T) {
	mock := &mockTestingT{}
	// Ensure we don't crash or behave weirdly if formatting string is passed
	a := "a"
	b := "b%s"
	diff.Diff(a, b, mock)
	if !mock.failed {
		t.Error("Should fail")
	}
	// Verify errorMsg contains "b%s" literally, not formatted
	// But Diff output will contain "b%s"
	// Errorf("%s", output) -> output is printed as string.
	// mock.Errorf calls Sprintf(format, args...) -> Sprintf("%s", output) -> output.
	// So output is preserved.
}
