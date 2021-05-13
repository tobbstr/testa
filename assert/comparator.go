package assert

import (
	"errors"
	"reflect"
)

// nilabe types
var nilableKinds []reflect.Kind = []reflect.Kind{
	reflect.Chan, reflect.Func, reflect.Map,
	reflect.Interface, reflect.Ptr, reflect.Slice,
}

type comparator interface {
	// Equals compares for equality, the 'want' argument (expected value) with the observed.
	// They are considered equal if both are nil or if they're deeply equal according to
	// reflect.DeepEqual's definition of equal.
	Equals(want interface{}) bool

	// EqualsElementsInIgnoringOrder compares for equality, ignoring order, the 'want' argument
	// (expected value) with the observed. Valid types for comparison are slices and arrays,
	// but it's also valid to compare slices with arrays and vice versa. Comparing other types
	// or if the values being compared are not equal, the function under test is marked as
	// having failed.
	// Two sequences of elements are equal if their number of elements are the
	// same, and if their elements are equal ignoring order.
	EqualsElementsInIgnoringOrder(want interface{})

	// IsEmpty checks whether the observed value is empty. If not empty, the function
	// under test is marked as having failed.
	// Arrays, channels, maps and slices are considered empty if they're nil or has zero
	// length.
	// Pointers are considered empty if the dereferenced values are nil.
	// For all other types, the zero value is considered empty.
	IsEmpty()

	// IsNil checks whether the observed value is nil. If not nil, the function
	// under test is marked as having failed.
	IsNil()

	// IsNotEmpty checks whether the observed value isn't empty. If empty, the function
	// under test is marked as having failed.
	IsNotEmpty()

	// IsNotNil checks whether the observed value is not nil. If nil, the function
	// under test is marked as having failed.
	IsNotNil()
}

func equals(got, want interface{}, errorf func(string, ...interface{})) bool {
	if err := validateArgsForEqualsFn(got, want); err != nil {
		errorf("validation of args for Equal method failed: %v", err)
		return false
	}

	if got == nil && want == nil {
		return true
	}

	if got == nil || want == nil {
		errorf("Error: got %v, expected: %v", got, want)
		return false
	}

	if !reflect.DeepEqual(got, want) {
		errorf("Error: got %v, expected: %v", got, want)
		return false
	}
	return true
}

func validateArgsForEqualsFn(a, b interface{}) error {
	if a == nil && b == nil {
		return nil
	}
	if isFunc(a) || isFunc(b) {
		return errors.New("cannot compare equality for function type")
	}
	return nil
}

func isFunc(arg interface{}) bool {
	if arg == nil {
		return false
	}
	return reflect.TypeOf(arg).Kind() == reflect.Func
}

func isNil(got interface{}, errorf func(string, ...interface{})) {
	if got == nil {
		return
	}
	value := reflect.ValueOf(got)
	kind := value.Kind()

	for _, nilableKind := range nilableKinds {
		if kind == nilableKind {
			if value.IsNil() {
				return
			}
			break
		}
	}
	errorf("expected nil value, found %+v", got)
}

func isNotEmpty(got interface{}, errorf func(string, ...interface{})) {
	if isEmpty(got) {
		errorf("expected %#v not to be empty, found empty", got)
	}
}

func isNotNil(got interface{}, errorf func(string, ...interface{})) {
	if got == nil {
		errorf("expected non-nil reference, found nil")
	}
	value := reflect.ValueOf(got)
	kind := value.Kind()

	for _, nilableKind := range nilableKinds {
		if kind == nilableKind {
			if value.IsNil() {
				errorf("expected non-nil value for nilable kind, found nil")
				return
			}
			return
		}
	}
}

func equalElementsIgnoringOrder(got interface{}, want interface{}, errorf func(string, ...interface{})) {
	if !isList(got, errorf) || !isList(want, errorf) {
		return
	}

	if isEmpty(got) && isEmpty(want) {
		return
	}

	matchingElemCountFor := make(map[interface{}]int)
	gotElemSequence := reflect.ValueOf(got)
	gotElemSequenceLen := gotElemSequence.Len()
	wantElemCountFor := make(map[interface{}]int)
	wantElemSequence := reflect.ValueOf(want)
	wantElemSequenceLen := wantElemSequence.Len()

	for w := 0; w < wantElemSequenceLen; w++ {
		wantElem := wantElemSequence.Index(w).Interface()
		wantElemCountFor[wantElem] = wantElemCountFor[wantElem] + 1
		for g := 0; g < gotElemSequenceLen; g++ {
			gotElem := gotElemSequence.Index(g).Interface()
			if reflect.DeepEqual(wantElem, gotElem) {
				matchingElemCountFor[gotElem] = matchingElemCountFor[gotElem] + 1
			}
		}
	}

	wantElemCount := len(wantElemCountFor)
	gotMatchingElemCount := len(matchingElemCountFor)

	if wantElemCount != gotMatchingElemCount {
		errorf("elements mismatch, expected %d matching elements found %d", wantElemCount, gotMatchingElemCount)
	}
}

func isEmpty(obj interface{}) bool {
	if obj == nil {
		return true
	}

	objValue := reflect.ValueOf(obj)
	objKind := objValue.Kind()

	switch objKind {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		return objValue.Len() == 0
	case reflect.Ptr:
		if objValue.IsNil() {
			return true
		}
		return isEmpty(objValue.Elem().Interface())
	default:
		objTypeZeroValue := reflect.Zero(objValue.Type())
		return reflect.DeepEqual(obj, objTypeZeroValue.Interface())
	}
}

func isList(list interface{}, errorf func(string, ...interface{})) bool {
	kind := reflect.ValueOf(list).Kind()
	if kind != reflect.Array && kind != reflect.Slice {
		//fmt.Printf("unsupported argument type for 'list', got %s expected array or slice\n", kind)
		errorf("unsupported argument type for 'list', expected array or slice, found %s", kind)
		return false
	}
	return true
}