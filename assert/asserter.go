package assert

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

// nilabe types
var nilableKinds = []reflect.Kind{
	reflect.Chan, reflect.Func, reflect.Map,
	reflect.Interface, reflect.Ptr, reflect.Slice,
}

// New returns an assert function, which is used to make assertions.
// If any assertion fails using this function, code execution is allowed to continue,
// but the test is marked as having failed.
func New(t *testing.T) func(got interface{}) asserter {
	return func(got interface{}) asserter {
		return asserter{
			got:   got,
			t:     t,
			fatal: false,
		}
	}
}

// NewFatal returns an assert function, which is used to make assertions.
// If any assertion fails using this function, code execution is immediately stopped
// and the test is marked as having failed.
func NewFatal(t *testing.T) func(got interface{}) asserter {
	return func(got interface{}) asserter {
		return asserter{
			got:   got,
			t:     t,
			fatal: true,
		}
	}
}

type asserter struct {
	got   interface{}
	t     *testing.T
	fatal bool
}

func (a *asserter) errorf(msg string, want interface{}, hasWant bool) {
	if a.fatal {
		a.t.Fatal(errorMsg(msg, want, a.got, hasWant))
		return
	}

	a.t.Fail()
	a.t.Log(errorMsg(msg, want, a.got, hasWant))
}

// Equals asserts the observed value equals the 'want' argument (expected value).
// They are considered equal if both are nil or if they're deeply equal according to
// reflect.DeepEqual's definition of equal.
func (a asserter) Equals(want interface{}) bool {
	if err := validateArgsForEqualsFn(a.got, want); err != nil {
		a.errorf(fmt.Sprintf("Invalid argument: %v", err), want, true)
		return false
	}
	if !equals(a.got, want) {
		a.errorf("Observed and expected values must be equal", want, true)
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
func (a asserter) IgnoringOrderEqualsElementsIn(want interface{}) bool {
	if !isList(a.got) || !isList(want) {
		a.errorf("Invalid argument", want, true)
		return false
	}

	if isEmpty(a.got) && isEmpty(want) {
		return true
	}

	matchingElemCountFor := make(map[interface{}]int)
	gotElemSequence := reflect.ValueOf(a.got)
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
		a.errorf(fmt.Sprintf("No of matching elements must be equal, want = %d, got = %d", wantElemCount, gotMatchingElemCount), want, true)
		return false
	}

	return true
}

func isList(list interface{}) bool {
	kind := reflect.ValueOf(list).Kind()
	if kind != reflect.Array && kind != reflect.Slice {
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
		a.errorf("Observed value must be empty", nil, false)
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
		a.errorf("Observed value must be a function", nil, false)
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
		a.errorf("Observed value must be nil", nil, false)
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
	isEmpty := isEmpty(a.got)
	if isEmpty {
		a.errorf("Observed value must non-empty", nil, false)
	}
	return !isEmpty
}

// IsNotNil asserts the observed value is not nil. If nil, the function
// under test is marked as having failed.
func (a asserter) IsNotNil() bool {
	isNil := isNil(a.got)
	if isNil {
		a.errorf("Observed value must not be nil", nil, false)
	}
	return !isNil
}

// IsTrue asserts the observed value is true. Otherwise, the function
// under test is marked as having failed. Only a boolean value of true
// returns true, for all other cases it returns false.
func (a asserter) IsTrue() bool {
	isTrue := isTrue((a.got))
	if !isTrue {
		a.errorf("Observed value must be true", nil, false)
	}
	return isTrue
}

func isTrue(got interface{}) bool {
	if got == nil {
		return false
	}

	value := reflect.ValueOf(got)
	kind := value.Kind()

	if kind != reflect.Bool {
		return false
	}

	return value.Bool()
}

// IsFalse asserts the observed value is false. Otherwise, the function
// under test is marked as having failed. Only a boolean value of false
// returns true, for all other cases it returns false.
func (a asserter) IsFalse() bool {
	isFalse := isFalse(a.got)
	if !isFalse {
		a.errorf("Observed value must be false", nil, false)
	}
	return isFalse
}

func isFalse(got interface{}) bool {
	if got == nil {
		return false
	}

	value := reflect.ValueOf(got)
	kind := value.Kind()

	if kind != reflect.Bool {
		return false
	}

	gotTrue := value.Bool()

	return !gotTrue
}

// IsPointerWithSameAddressAs asserts the observed pointer points to the same memory address as
// the 'want' pointer. Both the observed and want values must be pointers. If not, or if
// the pointers don't point to the same memory address, the function under test is marked
// as having failed.
func (a asserter) IsPointerWithSameAddressAs(want interface{}) bool {
	isPointerWithSameAddressAs := isPointerWithSameAddressAs(a.got, want)
	if !isPointerWithSameAddressAs {
		a.errorf("Observed pointer must be the same as the expected", want, true)
	}
	return isPointerWithSameAddressAs
}

func isPointerWithSameAddressAs(got, want interface{}) bool {
	if got == nil || want == nil {
		return false
	}

	gotKind := reflect.ValueOf(got).Kind()
	wantKind := reflect.ValueOf(want).Kind()

	if gotKind != reflect.Ptr || wantKind != reflect.Ptr {
		return false
	}

	if got != want {
		return false
	}
	return true
}

// NotEquals asserts the observed value is not equal to the 'want' argument. It performs the
// same comparison as the Equals method, but inverts the result. If they are equal, the function
// under test is marked as having failed.
func (a asserter) NotEquals(want interface{}) bool {
	if err := validateArgsForEqualsFn(a.got, want); err != nil {
		a.errorf(fmt.Sprintf("Invalid argument: %v", err), want, true)
		return false
	}
	if equals(a.got, want) {
		a.errorf("Observed and expected values must be unequal", want, true)
		return false
	}

	return true
}

// IsJSONEqualTo asserts the observed value is valid JSON and that it equals the 'want' argument.
// If not equal, the function under test is marked as having failed.
func (a asserter) IsJSONEqualTo(want interface{}) bool {
	if isNil(a.got) && isNil(want) {
		return true
	}

	if a.got == nil || want == nil {
		a.errorf("Invalid argument", want, true)
		return false
	}

	gotValue := reflect.ValueOf(a.got)
	gotKind := gotValue.Kind()

	if gotKind != reflect.String && gotKind != reflect.Slice {
		a.errorf("Invalid observed argument type, only string and []byte allowed", want, true)
		return false
	}

	wantValue := reflect.ValueOf(want)
	wantKind := wantValue.Kind()

	if wantKind != reflect.String && wantKind != reflect.Slice {
		a.errorf("Expected value must be a string or slice of bytes", want, true)
		return false
	}

	var gotAsBytes []byte
	if gotKind == reflect.String {
		gotAsBytes = []byte(gotValue.String())
	} else {
		if gotValue.Type() == reflect.TypeOf([]byte(nil)) {
			gotAsBytes = gotValue.Bytes()
		} else {
			a.errorf("Observed value must be a string or slice of bytes", want, true)
			return false
		}
	}
	if isEmpty(gotAsBytes) {
		gotAsBytes = []byte("{}")
	}

	var wantAsBytes []byte
	if wantKind == reflect.String {
		wantAsBytes = []byte(wantValue.String())
	} else {
		if wantValue.Type() == reflect.TypeOf([]byte(nil)) {
			wantAsBytes = wantValue.Bytes()
		} else {
			a.errorf("Expected slice must be a slice of bytes", want, true)
			return false
		}
	}
	if isEmpty(wantAsBytes) {
		wantAsBytes = []byte("{}")
	}

	var got1 interface{}
	var want1 interface{}

	var err error
	err = json.Unmarshal(gotAsBytes, &got1)
	if err != nil {
		a.errorf("Could not JSON-unmarshal observed value", want, true)
		return false
	}
	err = json.Unmarshal(wantAsBytes, &want1)
	if err != nil {
		a.errorf("Could not JSON-unmarshal expected value", want, true)
		return false
	}

	if !reflect.DeepEqual(got1, want1) {
		a.errorf("Values are not equal", want, true)
		return false
	}
	return true
}

// IsWantedError asserts the observed value is an error and that it's wanted. If it's not, the function
// under test is marked as having failed.
// This function's purpose is to simplify testing of error values returned by functions.
// See the below example:
//
//		got, err := FuncToTest()
//		assert(err).IsWantedError(wantErr) // where wantErr is a bool
//
func (a asserter) IsWantedError(wantErr bool) bool {
	if wantErr && isNil(a.got) {
		a.errorf("Observed value must not be nil", wantErr, true)
		return false
	}
	if !wantErr && !isNil(a.got) {
		a.errorf("Observed value must be nil", wantErr, true)
		return false
	}
	if !isNil(a.got) {
		if _, ok := a.got.(error); !ok {
			a.errorf("Observed value must be an error", wantErr, true)
			return false
		}
	}

	return true
}

// IsType asserts the observed value is the wanted type. If it's not, the function under test
// is marked as having failed.
//
// 	Example 1. Asserts got is a slice of bytes
//	assert(got).IsType([]byte{})
//
//	Example 2. Asserts got is a string
//	assert(got).IsType("")
//
//	Example 3. Asserts got is a func with a specific signature
//	assert(got).IsType( func(a, b int) int { return 5 } )
func (a asserter) IsType(want interface{}) bool {
	isType := isType(a.got, want)
	if !isType {
		a.errorf("Observed and expected values must be of the same Type", want, true)
		return false
	}
	return true
}

func isType(got, want interface{}) bool {
	if isNil(got) && isNil(want) {
		return true
	}

	gotValue := reflect.ValueOf(got)
	wantValue := reflect.ValueOf(want)

	if reflect.TypeOf((func())(nil)) == reflect.TypeOf(nil) {
		fmt.Printf("\n\ninterface{}(nil) and nil are equal\n\n")
	}

	if gotValue.Kind() == reflect.Ptr && wantValue.Kind() == reflect.Ptr {
		if gotValue.Elem().IsValid() && wantValue.Elem().IsValid() {
			return reflect.TypeOf(gotValue.Elem().Interface()) == reflect.TypeOf(wantValue.Elem().Interface())
		}
		return false
	}

	return reflect.TypeOf(got) == reflect.TypeOf(want)
}

func isChan(arg interface{}) bool {
	if arg == nil {
		return false
	}
	return reflect.TypeOf(arg).Kind() == reflect.Chan
}

// Implements asserts the observed value implements the wanted interface. If it does not, the function under test
// is marked as having failed.
//
//		got:	The observed value to check for interface implementations.
//				- Only non-nil values are allowed
//				- Both values and pointers are allowed
//
//		want:	The interface the observed value is checked against.
//				- Only pointers to interfaces are allowed
//
//	Example: Asserts *strings.Reader implements io.Reader
//		assert(strings.NewReader("dummy-str")).Implements((*io.Reader)(nil))
func (a asserter) Implements(want interface{}) bool {
	if isNil(a.got) || want == nil {
		a.errorf("Observed/Expected value must non-nil", want, true)
		return false
	}

	wantKind := reflect.ValueOf(want).Kind()

	if wantKind != reflect.Ptr {
		a.errorf("Expected value must be a pointer to an interface", want, true)
		return false
	}

	wantType := reflect.TypeOf(want).Elem()

	if !reflect.TypeOf(a.got).Implements(wantType) {
		a.errorf("Observed value must implement expected interface", want, true)
		return false
	}

	return true
}
