package assert

import (
	"testing"
)

type defaultAssert struct {
	t *testing.T
}

// NewFatal returns an assert interface implementation which does not allow code
// execution to continue when tests fail.
func NewFatal(t *testing.T) func(got interface{}) comparator {
	return func(got interface{}) comparator {
		return fatalComparator{
			t:   t,
			got: got,
		}
	}
}

// New returns an assert interface implementation which allows code execution
// to continue when tests fail.
func New(t *testing.T) func(got interface{}) comparator {
	return func(got interface{}) comparator {
		return defaultComparator{
			t:   t,
			got: got,
		}
	}
}
