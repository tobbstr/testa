package assert

import (
	"testing"
)

// NewFatal returns an assert function, which is used to make assertions.
// If any assertion fails using this function, code execution is immediately stopped
// and the test is marked as having failed.
func NewFatal(t *testing.T) func(got interface{}) comparator {
	return func(got interface{}) comparator {
		return fatalComparator{
			t:   t,
			got: got,
		}
	}
}

// New returns an assert function, which is used to make assertions.
// If any assertion fails using this function, code execution is allowed to continue,
// but the test is marked as having failed.
func New(t *testing.T) func(got interface{}) comparator {
	return func(got interface{}) comparator {
		return defaultComparator{
			t:   t,
			got: got,
		}
	}
}
