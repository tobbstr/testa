package assert

import "testing"

type fatalComparator struct {
	t   *testing.T
	got interface{}
}

func (f fatalComparator) Equals(want interface{}) {
	if err := validateArgsForEqualsFn(f.got, want); err != nil {
		f.t.Fatalf("validation of args for equals function failed: %v", err)
		return
	}
	if !equals(f.got, want) {
		f.t.Fatalf("expected assert(%#v).Equals(%#v) to be true, found false", f.got, want)
		return
	}
}

func (f fatalComparator) IgnoringOrderEqualsElementsIn(want interface{}) {
	ignoringOrderEqualsElementsIn(f.got, want, f.t.Fatalf)
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

func (f fatalComparator) IsPointerWithSameAddressAs(want interface{}) {
	isPointerWithSameAddressAs(f.got, want, f.t.Fatalf)
}

func (f fatalComparator) NotEquals(want interface{}) {
	if err := validateArgsForEqualsFn(f.got, want); err != nil {
		f.t.Fatalf("validation of args for equals function failed: %v", err)
		return
	}
	if equals(f.got, want) {
		f.t.Fatalf("expected assert(%#v).NotEquals(%#v) to be true, found false", f.got, want)
		return
	}
}
