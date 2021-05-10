package assert

import (
	"testing"
)

type assert interface {
	// That prepares the observed value, the 'got' argument', for comparison by
	// returning a comparator interface, with which the observed value can be
	// compared.
	That(got interface{}) comparator
}

type defaultAssert struct {
	t *testing.T
}

// NewFatal returns an assert interface implementation which does not allow code
// execution to continue when tests fail.
func NewFatal(t *testing.T) assert {
	return fatalAssert{
		t: t,
	}
}

// New returns an assert interface implementation which allows code execution
// to continue when tests fail.
func New(t *testing.T) assert {
	return defaultAssert{
		t: t,
	}
}

func (c defaultAssert) That(got interface{}) comparator {
	return defaultComparator{
		t:   c.t,
		got: got,
	}
}

type fatalAssert struct {
	t *testing.T
}

func (f fatalAssert) That(got interface{}) comparator {
	return fatalComparator{
		t:   f.t,
		got: got,
	}
}
