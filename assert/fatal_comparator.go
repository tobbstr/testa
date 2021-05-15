package assert

import "testing"

type fatalComparator struct {
	t   *testing.T
	got interface{}
}

func (f fatalComparator) Equals(want interface{}) bool {
	return equals(f.got, want, f.t.Fatalf)
}

func (f fatalComparator) EqualsElementsInIgnoringOrder(want interface{}) {
	equalElementsIgnoringOrder(f.got, want, f.t.Fatalf)
}

func (f fatalComparator) IsEmpty() {
	if !isEmpty(f.got) {
		f.t.Fatalf("expected observed value to be empty, found %#v", f.got)
	}
}

func (f fatalComparator) IsFunction() {
	if !isFunc(f.got) {
		f.t.Fatalf("expected observed value to be a function value, found %#v", f.got)
	}
}

func (f fatalComparator) IsNil() {
	if !isNil(f.got) {
		f.t.Errorf("expected the observed value to be nil, found %#v", f.got)
	}
}

func (f fatalComparator) IsNotEmpty() {
	isNotEmpty(f.got, f.t.Fatalf)
}

func (f fatalComparator) IsNotNil() {
	isNotNil(f.got, f.t.Fatalf)
}

func (f fatalComparator) IsTrue() {
	isTrue(f.got, f.t.Fatalf)
}

func (f fatalComparator) IsFalse() {
	isFalse(f.got, f.t.Fatalf)
}
