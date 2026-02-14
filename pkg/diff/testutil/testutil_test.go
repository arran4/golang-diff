package testutil_test

import (
	"testing"
	"github.com/arran4/golang-diff/pkg/diff"
	"github.com/arran4/golang-diff/pkg/diff/testutil"
)

func TestFailIfMismatch(t *testing.T) {
	var _ diff.TestingT = testutil.FailIfMismatch(t)
}
