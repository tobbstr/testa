package assert

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

// nilabe types
var nilableKinds = []reflect.Kind{
	reflect.Chan, reflect.Func, reflect.Map,
	reflect.Interface, reflect.Ptr, reflect.Slice,
}

type asserter struct {
	got    interface{}
	errorf func(string, ...interface{})
}

// New returns an assert function, which is used to make assertions.
// If any assertion fails using this function, code execution is allowed to continue,
// but the test is marked as having failed.
func New(t *testing.T) func(got interface{}) asserter {
	return func(got interface{}) asserter {
		return asserter{
			errorf: t.Errorf,
			got:    got,
		}
	}
}

// NewFatal returns an assert function, which is used to make assertions.
// If any assertion fails using this function, code execution is immediately stopped
// and the test is marked as having failed.
func NewFatal(t *testing.T) func(got interface{}) asserter {
	return func(got interface{}) asserter {
		return asserter{
			errorf: t.Fatalf,
			got:    got,
		}
	}
}

// Equals asserts the observed value equals the 'want' argument (expected value).
// They are considered equal if both are nil or if they're deeply equal according to
// reflect.DeepEqual's definition of equal.
func (a asserter) Equals(want interface{}) {
	if err := validateArgsForEqualsFn(a.got, want); err != nil {
		a.errorf("validation of args for equals function failed: %v", err)
		return
	}
	if !equals(a.got, want) {
		a.errorf("expected assert(%#v).Equals(%#v) to be true, found false", a.got, want)
		return
	}
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

func equals(got, want interface{}) bool {
	if got == nil && want == nil {
		return true
	}

	if got == nil || want == nil {
		return false
	}

	if !reflect.DeepEqual(got, want) {
		return false
	}
	return true
}

// IgnoringOrderEqualsElementsIn asserts the observed value is equal to the 'want' argument
// (expected value), ignoring order. Valid types for comparison are slices and arrays, but it's
// also valid to compare slices with arrays and vice versa. Comparing other types or if the
// values being compared are not equal, the function under test is marked as having failed.
// Two sequences of elements are equal if their number of elements are the
// same, and if their elements are equal ignoring order.
func (a asserter) IgnoringOrderEqualsElementsIn(want interface{}) {
	ignoringOrderEqualsElementsIn(a.got, want, a.errorf)
}

func ignoringOrderEqualsElementsIn(got interface{}, want interface{}, errorf func(string, ...interface{})) {
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
		wantElemCountFor[wantElem]++
		for g := 0; g < gotElemSequenceLen; g++ {
			gotElem := gotElemSequence.Index(g).Interface()
			if reflect.DeepEqual(wantElem, gotElem) {
				matchingElemCountFor[gotElem]++
			}
		}
	}

	wantElemCount := len(wantElemCountFor)
	gotMatchingElemCount := len(matchingElemCountFor)

	if wantElemCount != gotMatchingElemCount {
		errorf("elements mismatch, expected %d matching elements found %d", wantElemCount, gotMatchingElemCount)
	}
}

func isList(list interface{}, errorf func(string, ...interface{})) bool {
	kind := reflect.ValueOf(list).Kind()
	if kind != reflect.Array && kind != reflect.Slice {
		// fmt.Printf("unsupported argument type for 'list', got %s expected array or slice\n", kind)
		errorf("unsupported argument type for 'list', expected array or slice, found %s", kind)
		return false
	}
	return true
}

// IsEmpty asserts the observed value is empty. If not empty, the function under test is
// marked as having failed.
// Arrays, channels, maps and slices are considered empty if they're nil or has zero length.
// Pointers are considered empty if the referenced values are nil.
// For all other types, the zero value is considered empty.
func (a asserter) IsEmpty() bool {
	isEmpty := isEmpty(a.got)
	if !isEmpty {
		a.errorf("expected observed value to be empty, found %#v", a.got)
	}
	return isEmpty
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

// IsFunction asserts the observed value is a function value. If not, the function under test
// is marked as having failed.
func (a asserter) IsFunction() bool {
	isFunc := isFunc(a.got)
	if !isFunc {
		a.errorf("expected observed value to be a function value, found %#v", a.got)
	}
	return isFunc
}

func isFunc(arg interface{}) bool {
	if arg == nil {
		return false
	}
	return reflect.TypeOf(arg).Kind() == reflect.Func
}

// IsNil asserts the observed value is nil. If not nil, the function
// under test is marked as having failed.
func (a asserter) IsNil() bool {
	isNil := isNil(a.got)
	if !isNil {
		a.errorf("expected the observed value to be nil, found %#v", a.got)
	}
	return isNil
}

func isNil(got interface{}) bool {
	if got == nil {
		return true
	}
	value := reflect.ValueOf(got)
	kind := value.Kind()

	for _, nilableKind := range nilableKinds {
		if kind == nilableKind {
			if value.IsNil() {
				return true
			}
			break
		}
	}
	return false
}

// IsNotEmpty asserts the observed value isn't empty. If empty, the function
// under test is marked as having failed.
func (a asserter) IsNotEmpty() bool {
	return isNotEmpty(a.got, a.errorf)
}

func isNotEmpty(got interface{}, errorf func(string, ...interface{})) bool {
	isEmpty := isEmpty(got)
	if isEmpty {
		errorf("expected %#v not to be empty, found empty", got)
	}
	return isEmpty
}

// IsNotNil asserts the observed value is not nil. If nil, the function
// under test is marked as having failed.
func (a asserter) IsNotNil() bool {
	return isNotNil(a.got, a.errorf)
}

func isNotNil(got interface{}, errorf func(string, ...interface{})) bool {
	if got == nil {
		errorf("expected non-nil reference, found nil")
		return false
	}
	value := reflect.ValueOf(got)
	kind := value.Kind()

	for _, nilableKind := range nilableKinds {
		if kind == nilableKind {
			if value.IsNil() {
				errorf("expected non-nil value for nilable kind, found nil")
				return false
			}
		}
	}
	return true
}

// IsTrue asserts the observed value is true. Otherwise, the function
// under test is marked as having failed.
func (a asserter) IsTrue() bool {
	return isTrue(a.got, a.errorf)
}

func isTrue(got interface{}, errorf func(string, ...interface{})) bool {
	if got == nil {
		errorf("expected observed value to be true, found %v", got)
		return false
	}

	value := reflect.ValueOf(got)
	kind := value.Kind()

	if kind != reflect.Bool {
		errorf("expected observed value to be true, found %v", got)
		return false
	}

	gotTrue := value.Bool()
	if !gotTrue {
		errorf("expected observed value to be true, found %v", got)
		return false
	}
	return true
}

// IsFalse asserts the observed value is false. Otherwise, the function
// under test is marked as having failed.
func (a asserter) IsFalse() bool {
	return isFalse(a.got, a.errorf)
}

func isFalse(got interface{}, errorf func(string, ...interface{})) bool {
	if got == nil {
		errorf("expected observed value to be false, found %v", got)
		return false
	}

	value := reflect.ValueOf(got)
	kind := value.Kind()

	if kind != reflect.Bool {
		errorf("expected observed value to be false, found %v", got)
		return false
	}

	gotTrue := value.Bool()
	if gotTrue {
		errorf("expected observed value to be false, found %v", got)
		return false
	}
	return true
}

// IsPointerWithSameAddressAs asserts the observed pointer points to the same memory address as
// the 'want' pointer. Both the observed and want values must be pointers. If not, or if
// the pointers don't point to the same memory address, the function under test is marked
// as having failed.
func (a asserter) IsPointerWithSameAddressAs(want interface{}) bool {
	return isPointerWithSameAddressAs(a.got, want, a.errorf)
}

func isPointerWithSameAddressAs(got, want interface{}, errorf func(string, ...interface{})) bool {
	if got == nil || want == nil {
		errorf("expected both 'got' and 'want' to be non-nil values, found got %v and want %v", got, want)
		return false
	}

	gotKind := reflect.ValueOf(got).Kind()
	wantKind := reflect.ValueOf(want).Kind()

	if gotKind != reflect.Ptr || wantKind != reflect.Ptr {
		errorf("expected both 'got' and 'want' to be pointers, found got %v and want %v", got, want)
		return false
	}

	if got != want {
		errorf("expected 'want' to point to same memory address as the observed value, found got %v and want %v", got, want)
		return false
	}
	return true
}

// NotEquals asserts the observed value is not equal to the 'want' argument. It performs the
// same comparison as the Equals method, but inverts the result. If they are equal, the function
// under test is marked as having failed.
func (a asserter) NotEquals(want interface{}) {
	if err := validateArgsForEqualsFn(a.got, want); err != nil {
		a.errorf("validation of args for equals function failed: %v", err)
		return
	}
	if equals(a.got, want) {
		a.errorf("expected assert(%#v).NotEquals(%#v) to be true, found false", a.got, want)
		return
	}
}

// IsJSONEqualTo asserts the observed value is valid JSON and that it equals the 'want' argument.
// If not equal, the function under test is marked as having failed.
func (a asserter) IsJSONEqualTo(want interface{}) bool {
	if a.got == nil || want == nil {
		a.errorf("expected both 'got' and 'want' to be non-nil values, found got %v and want %v", a.got, want)
		return false
	}

	gotValue := reflect.ValueOf(a.got)
	gotKind := gotValue.Kind()

	if gotKind != reflect.String && gotKind != reflect.Slice {
		a.errorf("expected got to be a string or a slice of bytes")
		return false
	}

	wantValue := reflect.ValueOf(want)
	wantKind := wantValue.Kind()

	if wantKind != reflect.String && wantKind != reflect.Slice {
		a.errorf("expected want to be a string or a slice of bytes")
		return false
	}

	var gotAsBytes []byte
	if gotKind == reflect.String {
		gotAsBytes = []byte(gotValue.String())
	} else {
		if gotValue.Type() == reflect.TypeOf([]byte(nil)) {
			gotAsBytes = gotValue.Bytes()
		} else {
			a.errorf("expected got to be a string or slice of bytes")
			return false
		}
	}

	var wantAsBytes []byte
	if wantKind == reflect.String {
		wantAsBytes = []byte(wantValue.String())
	} else {
		if wantValue.Type() == reflect.TypeOf([]byte(nil)) {
			wantAsBytes = wantValue.Bytes()
		} else {
			a.errorf("expected want to be a string or slice of bytes")
			return false
		}
	}

	var got1 interface{}
	var want1 interface{}

	var err error
	err = json.Unmarshal(gotAsBytes, &got1)
	if err != nil {
		a.errorf("Error mashalling string 1 :: %s", err.Error())
		return false
	}
	err = json.Unmarshal(wantAsBytes, &want1)
	if err != nil {
		a.errorf("Error mashalling string 2 :: %s", err.Error())
		return false
	}

	if !reflect.DeepEqual(got1, want1) {
		a.errorf("expected observed JSON value to be equal 'want', found got %v and want %v", a.got, want)
		return false
	}
	return true
}

// IsWantedError asserts the observed value is an error and that it's wanted. If it's not, the function
// under test is marked as having failed.
// This function's purpose is to simplify testing of error values returned by functions, and it's best
// used by a Fatal asserter since it then also stops the execution of code. See the below example:
//
//		got, err := FuncToTest()
//		require(err).IsWantedError(wantErr) // where wantErr is a bool
//		if err != nil {
//			return
//		}
//
//		// Observed value assertions
//		assert(got).Equals(want)
func (a asserter) IsWantedError(wantErr bool) bool {
	if wantErr && isNil(a.got) {
		return false
	}
	if !wantErr && isNotNil(a.got, a.errorf) {
		return false
	}
	if _, ok := a.got.(error); !ok {
		a.errorf("expected observed value to be an error, found %#v", a.got)
		return false
	}

	return true
}
