package testutil

import (
	"github.com/arran4/golang-diff/pkg/diff"
	"testing"
)

func FailIfMismatch(t *testing.T) diff.TestingT {
	return t
}
