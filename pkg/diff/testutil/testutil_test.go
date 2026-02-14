package testutil_test

import (
	"github.com/arran4/golang-diff/pkg/diff"
	"github.com/arran4/golang-diff/pkg/diff/testutil"
	"testing"
)

func TestFailIfMismatch(t *testing.T) {
	var _ diff.TestingT = testutil.FailIfMismatch(t)
}
