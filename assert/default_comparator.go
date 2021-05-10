package assert

import "testing"

type defaultComparator struct {
	t   *testing.T
	got interface{}
}

func (d defaultComparator) Equals(want interface{}) bool {
	return equals(d.got, want, d.t.Errorf)
}

func (d defaultComparator) EqualsElementsInIgnoringOrder(want interface{}) {
	equalElementsIgnoringOrder(d.got, want, d.t.Errorf)
}

func (d defaultComparator) IsEmpty() {
	if !isEmpty(d.got) {
		d.t.Errorf("expected 'got' value to be empty, found %#v", d.got)
	}
}

func (d defaultComparator) IsNil() {
	isNil(d.got, d.t.Errorf)
}

func (d defaultComparator) IsNotNil() {
	isNotNil(d.got, d.t.Errorf)
}
