package testutil

import (
	"testing"
	"github.com/arran4/golang-diff/pkg/diff"
)

func FailIfMismatch(t *testing.T) diff.TestingT {
	return t
}
