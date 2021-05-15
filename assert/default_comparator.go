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
		d.t.Errorf("expected observed value to be empty, found %#v", d.got)
	}
}

func (d defaultComparator) IsFunction() {
	if !isFunc(d.got) {
		d.t.Errorf("expected observed value to be a function value, found %#v", d.got)
	}
}

func (d defaultComparator) IsNil() {
	if !isNil(d.got) {
		d.t.Errorf("expected the observed value to be nil, found %#v", d.got)
	}
}

func (d defaultComparator) IsNotEmpty() {
	isNotEmpty(d.got, d.t.Errorf)
}

func (d defaultComparator) IsNotNil() {
	isNotNil(d.got, d.t.Errorf)
}

func (d defaultComparator) IsTrue() {
	isTrue(d.got, d.t.Errorf)
}

func (d defaultComparator) IsFalse() {
	isFalse(d.got, d.t.Errorf)
}

func (d defaultComparator) PointsToSameAddressAs(want interface{}) {
	pointsToSameAddressAs(d.got, want, d.t.Errorf)
}
