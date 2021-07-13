package assert

import (
	"fmt"
	"testing"
)

var (
	zero = map[string]interface{}{
		"array":      [3]int{0, 0, 0},
		"bool":       false,
		"byte":       byte(0),
		"chan":       (chan (struct{}))(nil),
		"complex64":  complex64(0),
		"complex128": complex128(0),
		"float32":    float32(0),
		"float64":    float64(0),
		"func":       (func())(nil),
		"int":        int(0),
		"int8":       int8(0),
		"int16":      int16(0),
		"int32":      int32(0),
		"int64":      int64(0),
		"interface":  (interface{})(nil),
		"map":        (map[string]string)(nil),
		"ptr":        (*string)(nil),
		"rune":       rune(0),
		"slice":      ([]int)(nil),
		"string":     "",
		"struct":     struct{ i int }{i: 0},
		"uint":       uint(0),
		"uint8":      uint8(0),
		"uint16":     uint16(0),
		"uint32":     uint32(0),
		"uint64":     uint64(0),
		"error":      nil,
	}
	nonZero = map[string]interface{}{
		"array":      [3]int{1, 2, 3},
		"bool":       true,
		"byte":       byte(3),
		"chan":       make(chan (int)),
		"complex64":  2 + 2i,
		"complex128": complex(float64(5), float64(2)),
		"float32":    float32(5.1),
		"float64":    float64(2.3),
		"func":       func() {},
		"int":        int(3),
		"int8":       int8(3),
		"int16":      int16(3),
		"int32":      int32(3),
		"int64":      int64(3),
		"interface":  interface{}(3),
		"map":        map[string]string{"key": "value"},
		"ptr":        &[]int{1, 2, 3},
		"rune":       rune(3),
		"slice":      []int{1, 2, 3},
		"string":     "dummy-string",
		"struct":     struct{ i int }{i: 5},
		"uint":       uint(3),
		"uint8":      uint8(3),
		"uint16":     uint16(3),
		"uint32":     uint32(3),
		"uint64":     uint64(3),
		"error":      fmt.Errorf("dummy-error"),
	}
)

func TestEquals(t *testing.T) {
	type args struct {
		got  interface{}
		want interface{}
	}
	tests := []struct {
		name           string
		args           args
		want           bool
		wantTestFailed bool
	}{
		{
			name: "should pass for identical nil args",
			args: args{
				got:  nil,
				want: nil,
			},
			want:           true,
			wantTestFailed: false,
		},
		{
			name: "should pass for equal comparable non-nil args",
			args: args{
				got:  nonZero["int"],
				want: nonZero["int"],
			},
			want:           true,
			wantTestFailed: false,
		},
		{
			name: "should fail when want non-nil, but get nil",
			args: args{
				got:  nil,
				want: nonZero["int"],
			},
			want:           false,
			wantTestFailed: true,
		},
		{
			name: "should fail when want nil, but get comparable non-nil",
			args: args{
				got:  nonZero["struct"],
				want: nil,
			},
			want:           false,
			wantTestFailed: true,
		},
		{
			name: "should fail when want comparable non-nil, but get non-comparable func",
			args: args{
				got:  nonZero["func"],
				want: nonZero["int"],
			},
			want:           false,
			wantTestFailed: true,
		},
		{
			name: "should fail when want non-comparable func, but get comparable non-nil",
			args: args{
				got:  nonZero["string"],
				want: func() {},
			},
			want:           false,
			wantTestFailed: true,
		},
		{
			name: "should fail for unequal comparable non-nil args",
			args: args{
				got:  nonZero["int"],
				want: nonZero["string"],
			},
			want:           false,
			wantTestFailed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			dummyT := &testing.T{}
			dummyAssert := New(dummyT)

			// When
			got := dummyAssert(tt.args.got).Equals(tt.args.want)

			// Then
			if got != tt.want {
				t.Error(errorMsg(fmt.Sprintf("expected got = %#v, found %#v", tt.want, got), tt.want, got, true))
				return
			}
			if dummyT.Failed() != tt.wantTestFailed {
				t.Error(errorMsg(fmt.Sprintf("expected dummyT.Failed() = %#v, found %#v", tt.wantTestFailed, dummyT.Failed()), tt.wantTestFailed, dummyT.Failed(), true))
			}
		})
	}
}

func TestIgnoringOrderEqualsElementsIn(t *testing.T) {
	type args struct {
		got  interface{}
		want interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should fail when get bool",
			args: args{
				got:  nonZero["bool"],
				want: nonZero["slice"],
			},
			want: false,
		},
		{
			name: "should fail when get chan",
			args: args{
				got:  nonZero["chan"],
				want: nonZero["slice"],
			},
			want: false,
		},
		{
			name: "should fail when get complex64",
			args: args{
				got:  nonZero["complex64"],
				want: nonZero["slice"],
			},
			want: false,
		},
		{
			name: "should fail when get complex128",
			args: args{
				got:  nonZero["complex128"],
				want: nonZero["slice"],
			},
			want: false,
		},
		{
			name: "should fail when get float32",
			args: args{
				got:  nonZero["float32"],
				want: nonZero["slice"],
			},
			want: false,
		},
		{
			name: "should fail when get float64",
			args: args{
				got:  nonZero["float64"],
				want: nonZero["slice"],
			},
			want: false,
		},
		{
			name: "should fail when get func",
			args: args{
				got:  nonZero["func"],
				want: nonZero["slice"],
			},
			want: false,
		},
		{
			name: "should fail when get int",
			args: args{
				got:  nonZero["int"],
				want: nonZero["slice"],
			},
			want: false,
		},
		{
			name: "should fail when get int8",
			args: args{
				got:  nonZero["int8"],
				want: nonZero["slice"],
			},
			want: false,
		},
		{
			name: "should fail when get int16",
			args: args{
				got:  nonZero["int16"],
				want: nonZero["slice"],
			},
			want: false,
		},
		{
			name: "should fail when get int32",
			args: args{
				got:  nonZero["int32"],
				want: nonZero["slice"],
			},
			want: false,
		},
		{
			name: "should fail when get int64",
			args: args{
				got:  nonZero["int64"],
				want: nonZero["slice"],
			},
			want: false,
		},
		{
			name: "should fail when get map",
			args: args{
				got:  nonZero["map"],
				want: nonZero["slice"],
			},
			want: false,
		},
		{
			name: "should fail when get ptr",
			args: args{
				got:  nonZero["ptr"],
				want: nonZero["slice"],
			},
			want: false,
		},
		{
			name: "should fail when get string",
			args: args{
				got:  nonZero["string"],
				want: nonZero["slice"],
			},
			want: false,
		},
		{
			name: "should fail when get struct",
			args: args{
				got:  nonZero["struct"],
				want: nonZero["slice"],
			},
			want: false,
		},
		{
			name: "should pass when want slice and get the same set of elements in order",
			args: args{
				got:  []string{"a", "b", "c"},
				want: []string{"a", "b", "c"},
			},
			want: true,
		},
		{
			name: "should pass when want slice and get the same set of elements not in order",
			args: args{
				got:  []string{"b", "c", "a"},
				want: []string{"a", "b", "c"},
			},
			want: true,
		},
		{
			name: "should fail when want slice but get subset not in order",
			args: args{
				got:  []string{"c", "a"},
				want: []string{"a", "b", "c"},
			},
			want: false,
		},
		{
			name: "should fail when want slice but get subset in order",
			args: args{
				got:  []string{"b", "c"},
				want: []string{"a", "b", "c"},
			},
			want: false,
		},
		{
			name: "should pass when want array and get the same set of elements in order",
			args: args{
				got:  [3]string{"a", "b", "c"},
				want: [3]string{"a", "b", "c"},
			},
			want: true,
		},
		{
			name: "should pass when want array and get the same set of elements not in order",
			args: args{
				got:  [3]string{"b", "c", "a"},
				want: [3]string{"a", "b", "c"},
			},
			want: true,
		},
		{
			name: "should fail when want array but get subset not in order",
			args: args{
				got:  [2]string{"c", "a"},
				want: [3]string{"a", "b", "c"},
			},
			want: false,
		},
		{
			name: "should fail when want array but get subset in order",
			args: args{
				got:  [2]string{"b", "c"},
				want: [3]string{"a", "b", "c"},
			},
			want: false,
		},
		{
			name: "should pass when want array but get slice with same set of elements",
			args: args{
				got:  []string{"a", "b", "c"},
				want: [3]string{"a", "b", "c"},
			},
			want: true,
		},
		{
			name: "should pass when want slice but get array with same set of elements",
			args: args{
				got:  [3]string{"a", "b", "c"},
				want: []string{"a", "b", "c"},
			},
			want: true,
		},
		{
			name: "should pass when want and get empty slice",
			args: args{
				got:  []string{},
				want: []string{},
			},
			want: true,
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			dummyT := &testing.T{}
			dummyAssert := New(dummyT)
			assert := New(t)

			// When
			got := dummyAssert(tt.args.got).IgnoringOrderEqualsElementsIn(tt.args.want)

			// Then
			assert(got).Equals(tt.want)
			assert(dummyT.Failed()).Equals(!tt.want)
		})
	}
}

func TestIsEmpty(t *testing.T) {
	nonEmptyChan := make(chan int, 10)
	nonEmptyChan <- 4

	type args struct {
		got interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should pass when get empty array",
			args: args{
				got: [0]int{},
			},
			want: true,
		},
		{
			name: "should fail when get non-empty array",
			args: args{
				got: zero["array"],
			},
			want: false,
		},
		{
			name: "should pass when get zero-valued bool",
			args: args{
				got: zero["bool"],
			},
			want: true,
		},
		{
			name: "should fail when get true bool",
			args: args{
				got: true,
			},
			want: false,
		},
		{
			name: "should pass when get zero-valued byte",
			args: args{
				got: zero["byte"],
			},
			want: true,
		},
		{
			name: "should fail when get non-zero-valued byte",
			args: args{
				got: nonZero["byte"],
			},
			want: false,
		},
		{
			name: "should pass when get empty chan",
			args: args{
				got: zero["chan"],
			},
			want: true,
		},
		{
			name: "should fail when get non-empty chan",
			args: args{
				got: nonEmptyChan,
			},
			want: false,
		},
		{
			name: "should pass when get zero-valued complex64",
			args: args{
				got: zero["complex64"],
			},
			want: true,
		},
		{
			name: "should fail when get non-zero-valued complex64",
			args: args{
				got: nonZero["complex64"],
			},
			want: false,
		},
		{
			name: "should pass when get zero-valued complex128",
			args: args{
				got: zero["complex128"],
			},
			want: true,
		},
		{
			name: "should fail when get non-zero-valued complex128",
			args: args{
				got: nonZero["complex128"],
			},
			want: false,
		},
		{
			name: "should pass when get zero-valued float32",
			args: args{
				got: zero["float32"],
			},
			want: true,
		},
		{
			name: "should fail when get non-zero-valued float32",
			args: args{
				got: nonZero["float32"],
			},
			want: false,
		},
		{
			name: "should pass when get zero-valued func",
			args: args{
				got: zero["func"],
			},
			want: true,
		},
		{
			name: "should fail when get non-zero-valued func",
			args: args{
				got: nonZero["func"],
			},
			want: false,
		},
		{
			name: "should pass when get zero-valued int",
			args: args{
				got: zero["int"],
			},
			want: true,
		},
		{
			name: "should fail when get non-zero-valued int",
			args: args{
				got: nonZero["int"],
			},
			want: false,
		},
		{
			name: "should pass when get zero-valued interface",
			args: args{
				got: zero["interface"],
			},
			want: true,
		},
		{
			name: "should fail when get non-zero-valued interface",
			args: args{
				got: nonZero["interface"],
			},
			want: false,
		},
		{
			name: "should pass when get zero-valued map",
			args: args{
				got: zero["map"],
			},
			want: true,
		},
		{
			name: "should fail when get non-zero-valued map",
			args: args{
				got: nonZero["map"],
			},
			want: false,
		},
		{
			name: "should pass when get zero-valued ptr",
			args: args{
				got: zero["ptr"],
			},
			want: true,
		},
		{
			name: "should fail when get non-zero-valued ptr",
			args: args{
				got: nonZero["ptr"],
			},
			want: false,
		},
		{
			name: "should pass when get nil slice",
			args: args{
				got: zero["slice"],
			},
			want: true,
		},
		{
			name: "should pass when get empty slice",
			args: args{
				got: []int{},
			},
			want: true,
		},
		{
			name: "should fail when get non-zero-valued slice",
			args: args{
				got: nonZero["slice"],
			},
			want: false,
		},
		{
			name: "should pass when get zero-valued string",
			args: args{
				got: zero["string"],
			},
			want: true,
		},
		{
			name: "should fail when get non-zero-valued string",
			args: args{
				got: nonZero["string"],
			},
			want: false,
		},
		{
			name: "should pass when get zero-valued struct",
			args: args{
				got: zero["struct"],
			},
			want: true,
		},
		{
			name: "should fail when get non-zero-valued struct",
			args: args{
				got: nonZero["struct"],
			},
			want: false,
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			dummyT := &testing.T{}
			dummyAssert := New(dummyT)
			assert := New(t)

			// When
			got := dummyAssert(tt.args.got).IsEmpty()

			// Then
			assert(got).Equals(tt.want)
			assert(dummyT.Failed()).Equals(!tt.want)
		})
	}
}

func TestIsFunction(t *testing.T) {
	type args struct {
		got interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should fail when get array",
			args: args{
				got: nonZero["array"],
			},
			want: false,
		},
		{
			name: "should fail when get bool",
			args: args{
				got: nonZero["bool"],
			},
			want: false,
		},
		{
			name: "should fail when get byte",
			args: args{
				got: nonZero["byte"],
			},
			want: false,
		},
		{
			name: "should fail when get chan",
			args: args{
				got: nonZero["chan"],
			},
			want: false,
		},
		{
			name: "should fail when get complex64",
			args: args{
				got: nonZero["complex64"],
			},
			want: false,
		},
		{
			name: "should fail when get complex128",
			args: args{
				got: nonZero["complex128"],
			},
			want: false,
		},
		{
			name: "should fail when get float32",
			args: args{
				got: nonZero["float32"],
			},
			want: false,
		},
		{
			name: "should pass when get zero-valued func",
			args: args{
				got: zero["func"],
			},
			want: true,
		},
		{
			name: "should pass when get non-zero-valued func",
			args: args{
				got: nonZero["func"],
			},
			want: true,
		},
		{
			name: "should fail when get int",
			args: args{
				got: nonZero["int"],
			},
			want: false,
		},
		{
			name: "should fail when get interface",
			args: args{
				got: nonZero["interface"],
			},
			want: false,
		},
		{
			name: "should fail when get map",
			args: args{
				got: nonZero["map"],
			},
			want: false,
		},
		{
			name: "should fail when get ptr",
			args: args{
				got: nonZero["ptr"],
			},
			want: false,
		},
		{
			name: "should fail when get slice",
			args: args{
				got: nonZero["slice"],
			},
			want: false,
		},
		{
			name: "should fail when get string",
			args: args{
				got: nonZero["string"],
			},
			want: false,
		},
		{
			name: "should fail when get struct",
			args: args{
				got: nonZero["struct"],
			},
			want: false,
		},
		{
			name: "should fail when get nil",
			args: args{
				got: nil,
			},
			want: false,
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			dummyT := &testing.T{}
			dummyAssert := New(dummyT)
			assert := New(t)

			// When
			got := dummyAssert(tt.args.got).IsFunction()

			// Then
			assert(got).Equals(tt.want)
			assert(dummyT.Failed()).Equals(!tt.want)
		})
	}
}

func TestIsNil(t *testing.T) {
	type args struct {
		got interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should fail when get non-nilable",
			args: args{
				got: nonZero["byte"],
			},
			want: false,
		},
		{
			name: "should fail when get non-nil chan",
			args: args{
				got: nonZero["chan"],
			},
			want: false,
		},
		{
			name: "should pass when get nil chan",
			args: args{
				got: zero["chan"],
			},
			want: true,
		},
		{
			name: "should fail when get non-nil func",
			args: args{
				got: nonZero["func"],
			},
			want: false,
		},
		{
			name: "should pass when get nil func",
			args: args{
				got: zero["func"],
			},
			want: true,
		},
		{
			name: "should fail when get non-nil map",
			args: args{
				got: nonZero["map"],
			},
			want: false,
		},
		{
			name: "should pass when get nil map",
			args: args{
				got: zero["map"],
			},
			want: true,
		},
		{
			name: "should fail when get non-nil interface",
			args: args{
				got: nonZero["interface"],
			},
			want: false,
		},
		{
			name: "should pass when get nil interface",
			args: args{
				got: zero["interface"],
			},
			want: true,
		},
		{
			name: "should fail when get non-nil pointer",
			args: args{
				got: nonZero["ptr"],
			},
			want: false,
		},
		{
			name: "should pass when get nil pointer",
			args: args{
				got: zero["ptr"],
			},
			want: true,
		},
		{
			name: "should fail when get non-nil slice",
			args: args{
				got: nonZero["slice"],
			},
			want: false,
		},
		{
			name: "should pass when get nil slice",
			args: args{
				got: zero["slice"],
			},
			want: true,
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			dummyT := &testing.T{}
			dummyAssert := New(dummyT)
			assert := New(t)

			// When
			got := dummyAssert(tt.args.got).IsNil()

			// Then
			assert(got).Equals(tt.want)
			assert(dummyT.Failed()).Equals(!tt.want)
		})
	}
}

func TestIsNotEmpty(t *testing.T) {
	nonEmptyChan := make(chan int, 10)
	nonEmptyChan <- 4

	type args struct {
		got interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should fail when get empty array",
			args: args{
				got: [0]int{},
			},
			want: false,
		},
		{
			name: "should pass when get non-empty array",
			args: args{
				got: zero["array"],
			},
			want: true,
		},
		{
			name: "should fail when get zero-valued bool",
			args: args{
				got: zero["bool"],
			},
			want: false,
		},
		{
			name: "should pass when get true bool",
			args: args{
				got: true,
			},
			want: true,
		},
		{
			name: "should fail when get zero-valued byte",
			args: args{
				got: zero["byte"],
			},
			want: false,
		},
		{
			name: "should pass when get non-zero-valued byte",
			args: args{
				got: nonZero["byte"],
			},
			want: true,
		},
		{
			name: "should fail when get empty chan",
			args: args{
				got: zero["chan"],
			},
			want: false,
		},
		{
			name: "should pass when get non-empty chan",
			args: args{
				got: nonEmptyChan,
			},
			want: true,
		},
		{
			name: "should fail when get zero-valued complex64",
			args: args{
				got: zero["complex64"],
			},
			want: false,
		},
		{
			name: "should pass when get non-zero-valued complex64",
			args: args{
				got: nonZero["complex64"],
			},
			want: true,
		},
		{
			name: "should fail when get zero-valued complex128",
			args: args{
				got: zero["complex128"],
			},
			want: false,
		},
		{
			name: "should pass when get non-zero-valued complex128",
			args: args{
				got: nonZero["complex128"],
			},
			want: true,
		},
		{
			name: "should fail when get zero-valued float32",
			args: args{
				got: zero["float32"],
			},
			want: false,
		},
		{
			name: "should pass when get non-zero-valued float32",
			args: args{
				got: nonZero["float32"],
			},
			want: true,
		},
		{
			name: "should fail when get zero-valued func",
			args: args{
				got: zero["func"],
			},
			want: false,
		},
		{
			name: "should pass when get non-zero-valued func",
			args: args{
				got: nonZero["func"],
			},
			want: true,
		},
		{
			name: "should fail when get zero-valued int",
			args: args{
				got: zero["int"],
			},
			want: false,
		},
		{
			name: "should pass when get non-zero-valued int",
			args: args{
				got: nonZero["int"],
			},
			want: true,
		},
		{
			name: "should fail when get zero-valued interface",
			args: args{
				got: zero["interface"],
			},
			want: false,
		},
		{
			name: "should pass when get non-zero-valued interface",
			args: args{
				got: nonZero["interface"],
			},
			want: true,
		},
		{
			name: "should fail when get zero-valued map",
			args: args{
				got: zero["map"],
			},
			want: false,
		},
		{
			name: "should pass when get non-zero-valued map",
			args: args{
				got: nonZero["map"],
			},
			want: true,
		},
		{
			name: "should fail when get zero-valued ptr",
			args: args{
				got: zero["ptr"],
			},
			want: false,
		},
		{
			name: "should pass when get non-zero-valued ptr",
			args: args{
				got: nonZero["ptr"],
			},
			want: true,
		},
		{
			name: "should fail when get nil slice",
			args: args{
				got: zero["slice"],
			},
			want: false,
		},
		{
			name: "should fail when get empty slice",
			args: args{
				got: []int{},
			},
			want: false,
		},
		{
			name: "should pass when get non-zero-valued slice",
			args: args{
				got: nonZero["slice"],
			},
			want: true,
		},
		{
			name: "should fail when get zero-valued string",
			args: args{
				got: zero["string"],
			},
			want: false,
		},
		{
			name: "should pass when get non-zero-valued string",
			args: args{
				got: nonZero["string"],
			},
			want: true,
		},
		{
			name: "should fail when get zero-valued struct",
			args: args{
				got: zero["struct"],
			},
			want: false,
		},
		{
			name: "should pass when get non-zero-valued struct",
			args: args{
				got: nonZero["struct"],
			},
			want: true,
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			dummyT := &testing.T{}
			dummyAssert := New(dummyT)
			assert := New(t)

			// When
			got := dummyAssert(tt.args.got).IsNotEmpty()

			// Then
			assert(got).Equals(tt.want)
			assert(dummyT.Failed()).Equals(!tt.want)
		})
	}
}

func TestIsNotNil(t *testing.T) {
	type args struct {
		got interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should pass when get non-nilable",
			args: args{
				got: nonZero["int"],
			},
			want: true,
		},
		{
			name: "should pass when get non-nil chan",
			args: args{
				got: nonZero["chan"],
			},
			want: true,
		},
		{
			name: "should fail when get nil chan",
			args: args{
				got: zero["chan"],
			},
			want: false,
		},
		{
			name: "should pass when get non-nil func",
			args: args{
				got: nonZero["func"],
			},
			want: true,
		},
		{
			name: "should fail when get nil func",
			args: args{
				got: zero["func"],
			},
			want: false,
		},
		{
			name: "should pass when get non-nil map",
			args: args{
				got: nonZero["map"],
			},
			want: true,
		},
		{
			name: "should fail when get nil map",
			args: args{
				got: zero["map"],
			},
			want: false,
		},
		{
			name: "should pass when get non-nil interface",
			args: args{
				got: nonZero["interface"],
			},
			want: true,
		},
		{
			name: "should fail when get nil interface",
			args: args{
				got: zero["interface"],
			},
			want: false,
		},
		{
			name: "should pass when get non-nil pointer",
			args: args{
				got: nonZero["ptr"],
			},
			want: true,
		},
		{
			name: "should fail when get nil pointer",
			args: args{
				got: zero["ptr"],
			},
			want: false,
		},
		{
			name: "should pass when get non-nil slice",
			args: args{
				got: nonZero["slice"],
			},
			want: true,
		},
		{
			name: "should fail when get nil slice",
			args: args{
				got: zero["slice"],
			},
			want: false,
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			dummyT := &testing.T{}
			dummyAssert := New(dummyT)
			assert := New(t)

			// When
			got := dummyAssert(tt.args.got).IsNotNil()

			// Then
			assert(got).Equals(tt.want)
			assert(dummyT.Failed()).Equals(!tt.want)
		})
	}
}

func TestIsTrue(t *testing.T) {
	type args struct {
		got interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should fail when get nil",
			args: args{
				got: nil,
			},
			want: false,
		},
		{
			name: "should fail when get non-bool",
			args: args{
				got: nonZero["int"],
			},
			want: false,
		},
		{
			name: "should pass when get true bool",
			args: args{
				got: nonZero["bool"],
			},
			want: true,
		},
		{
			name: "should fail when get false bool",
			args: args{
				got: zero["bool"],
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			dummyT := &testing.T{}
			dummyAssert := New(dummyT)
			assert := New(t)

			// When
			got := dummyAssert(tt.args.got).IsTrue()

			// Then
			assert(got).Equals(tt.want)
			assert(dummyT.Failed()).Equals(!tt.want)
		})
	}
}

func TestIsFalse(t *testing.T) {
	type args struct {
		got interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should fail when get nil",
			args: args{
				got: nil,
			},
			want: false,
		},
		{
			name: "should fail when get non-bool",
			args: args{
				got: nonZero["int"],
			},
			want: false,
		},
		{
			name: "should fail when get true bool",
			args: args{
				got: nonZero["bool"],
			},
			want: false,
		},
		{
			name: "should pass when get false bool",
			args: args{
				got: zero["bool"],
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			dummyT := &testing.T{}
			dummyAssert := New(dummyT)
			assert := New(t)

			// When
			got := dummyAssert(tt.args.got).IsFalse()

			// Then
			assert(got).Equals(tt.want)
			assert(dummyT.Failed()).Equals(!tt.want)
		})
	}
}

func TestIsPointerWithSameAddressAs(t *testing.T) {
	type args struct {
		got  interface{}
		want interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should fail when get nil",
			args: args{
				got:  nil,
				want: nonZero["ptr"],
			},
			want: false,
		},
		{
			name: "should fail when want nil",
			args: args{
				got:  nonZero["ptr"],
				want: nil,
			},
			want: false,
		},
		{
			name: "should fail when get non-pointer",
			args: args{
				got:  nonZero["string"],
				want: &struct{}{},
			},
			want: false,
		},
		{
			name: "should fail when get pointer that points to another address",
			args: args{
				got:  nonZero["ptr"],
				want: &struct{}{},
			},
			want: false,
		},
		{
			name: "should pass when want and get pointer that points to same address",
			args: args{
				got:  nonZero["ptr"],
				want: nonZero["ptr"],
			},
			want: true,
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			dummyT := &testing.T{}
			dummyAssert := New(dummyT)
			assert := New(t)

			// When
			got := dummyAssert(tt.args.got).IsPointerWithSameAddressAs(tt.args.want)

			// Then
			assert(got).Equals(tt.want)
			assert(dummyT.Failed()).Equals(!tt.want)
		})
	}
}

func TestNotEquals(t *testing.T) {
	type args struct {
		got  interface{}
		want interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should fail for identical nil args",
			args: args{
				got:  nil,
				want: nil,
			},
			want: false,
		},
		{
			name: "should fail for equal comparable non-nil args",
			args: args{
				got:  nonZero["int"],
				want: nonZero["int"],
			},
			want: false,
		},
		{
			name: "should pass when want non-nil, but get nil",
			args: args{
				got:  nil,
				want: nonZero["int"],
			},
			want: true,
		},
		{
			name: "should pass when want nil, but get comparable non-nil",
			args: args{
				got:  nonZero["struct"],
				want: nil,
			},
			want: true,
		},
		{
			name: "should fail when want comparable non-nil, but get non-comparable func",
			args: args{
				got:  nonZero["func"],
				want: nonZero["int"],
			},
			want: false,
		},
		{
			name: "should fail when want non-comparable func, but get comparable non-nil",
			args: args{
				got:  nonZero["string"],
				want: func() {},
			},
			want: false,
		},
		{
			name: "should pass for unequal comparable non-nil args",
			args: args{
				got:  nonZero["int"],
				want: nonZero["string"],
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			dummyT := &testing.T{}
			dummyAssert := New(dummyT)
			assert := New(t)

			// When
			got := dummyAssert(tt.args.got).NotEquals(tt.args.want)

			// Then
			assert(got).Equals(tt.want)
			assert(dummyT.Failed()).Equals(!tt.want)
		})
	}
}

func TestIsJSONEqualTo(t *testing.T) {
	validJSONObjectAsString := `{"dummy-field": "dummy-value"}`
	validJSONObjectAsByteSlice := []byte(validJSONObjectAsString)

	invalidJSONObjectAsString := `{"dummy-field": "dummy-value"`

	differentValidJSONObjectAsString := `{"dummy-field": "dummy-value-2"}`

	type args struct {
		got  interface{}
		want interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should pass when want and get same valid JSON object as string",
			args: args{
				got:  validJSONObjectAsString,
				want: validJSONObjectAsString,
			},
			want: true,
		},
		{
			name: "should pass when want and get same valid JSON object as byte slices",
			args: args{
				got:  validJSONObjectAsByteSlice,
				want: validJSONObjectAsByteSlice,
			},
			want: true,
		},
		{
			name: "should pass when want same valid JSON object as byte slice, but get it as string",
			args: args{
				got:  validJSONObjectAsString,
				want: validJSONObjectAsByteSlice,
			},
			want: true,
		},
		{
			name: "should pass when want same valid JSON object as string, but get it as byte slice",
			args: args{
				got:  validJSONObjectAsByteSlice,
				want: validJSONObjectAsString,
			},
			want: true,
		},
		{
			name: "should fail when want valid JSON object, but get array",
			args: args{
				got:  nonZero["array"],
				want: validJSONObjectAsString,
			},
			want: false,
		},
		{
			name: "should fail when want valid JSON object, but get bool",
			args: args{
				got:  nonZero["bool"],
				want: validJSONObjectAsString,
			},
			want: false,
		},
		{
			name: "should fail when want valid JSON object, but get byte",
			args: args{
				got:  nonZero["byte"],
				want: validJSONObjectAsString,
			},
			want: false,
		},
		{
			name: "should fail when want valid JSON object, but get chan",
			args: args{
				got:  nonZero["chan"],
				want: validJSONObjectAsString,
			},
			want: false,
		},
		{
			name: "should fail when want valid JSON object, but get complex64",
			args: args{
				got:  nonZero["complex64"],
				want: validJSONObjectAsString,
			},
			want: false,
		},
		{
			name: "should fail when want valid JSON object, but get complex128",
			args: args{
				got:  nonZero["complex128"],
				want: validJSONObjectAsString,
			},
			want: false,
		},
		{
			name: "should fail when want valid JSON object, but get float32",
			args: args{
				got:  nonZero["float32"],
				want: validJSONObjectAsString,
			},
			want: false,
		},
		{
			name: "should fail when want valid JSON object, but get float64",
			args: args{
				got:  nonZero["float64"],
				want: validJSONObjectAsString,
			},
			want: false,
		},
		{
			name: "should fail when want valid JSON object, but get func",
			args: args{
				got:  nonZero["func"],
				want: validJSONObjectAsString,
			},
			want: false,
		},
		{
			name: "should fail when want valid JSON object, but get int",
			args: args{
				got:  nonZero["int"],
				want: validJSONObjectAsString,
			},
			want: false,
		},
		{
			name: "should fail when want valid JSON object, but get int8",
			args: args{
				got:  nonZero["int8"],
				want: validJSONObjectAsString,
			},
			want: false,
		},
		{
			name: "should fail when want valid JSON object, but get int16",
			args: args{
				got:  nonZero["int16"],
				want: validJSONObjectAsString,
			},
			want: false,
		},
		{
			name: "should fail when want valid JSON object, but get int32",
			args: args{
				got:  nonZero["int32"],
				want: validJSONObjectAsString,
			},
			want: false,
		},
		{
			name: "should fail when want valid JSON object, but get int64",
			args: args{
				got:  nonZero["int64"],
				want: validJSONObjectAsString,
			},
			want: false,
		},
		{
			name: "should fail when want valid JSON object, but get interface",
			args: args{
				got:  nonZero["interface"],
				want: validJSONObjectAsString,
			},
			want: false,
		},
		{
			name: "should fail when want valid JSON object, but get map",
			args: args{
				got:  nonZero["map"],
				want: validJSONObjectAsString,
			},
			want: false,
		},
		{
			name: "should fail when want valid JSON object, but get nil",
			args: args{
				got:  nil,
				want: validJSONObjectAsString,
			},
			want: false,
		},
		{
			name: "should fail when want valid JSON object, but get ptr",
			args: args{
				got:  nonZero["ptr"],
				want: validJSONObjectAsString,
			},
			want: false,
		},
		{
			name: "should fail when want valid JSON object, but get rune",
			args: args{
				got:  nonZero["rune"],
				want: validJSONObjectAsString,
			},
			want: false,
		},
		{
			name: "should fail when want valid JSON object, but get slice",
			args: args{
				got:  nonZero["slice"],
				want: validJSONObjectAsString,
			},
			want: false,
		},
		{
			name: "should fail when want valid JSON object, but get string",
			args: args{
				got:  nonZero["string"],
				want: validJSONObjectAsString,
			},
			want: false,
		},
		{
			name: "should fail when want valid JSON object, but get struct",
			args: args{
				got:  nonZero["struct"],
				want: validJSONObjectAsString,
			},
			want: false,
		},
		{
			name: "should fail when want valid JSON object, but get uint",
			args: args{
				got:  nonZero["uint"],
				want: validJSONObjectAsString,
			},
			want: false,
		},
		{
			name: "should fail when want valid JSON object, but get uint8",
			args: args{
				got:  nonZero["uint8"],
				want: validJSONObjectAsString,
			},
			want: false,
		},
		{
			name: "should fail when want valid JSON object, but get uint16",
			args: args{
				got:  nonZero["uint16"],
				want: validJSONObjectAsString,
			},
			want: false,
		},
		{
			name: "should fail when want valid JSON object, but get uint32",
			args: args{
				got:  nonZero["uint32"],
				want: validJSONObjectAsString,
			},
			want: false,
		},
		{
			name: "should fail when want valid JSON object, but get uint64",
			args: args{
				got:  nonZero["uint64"],
				want: validJSONObjectAsString,
			},
			want: false,
		},
		{
			name: "should fail when want arg is non-(string or []byte), but get valid JSON object",
			args: args{
				got:  validJSONObjectAsString,
				want: nonZero["bool"],
			},
			want: false,
		},
		{
			name: "should fail when want arg is slice of non-bytes, but get valid JSON object",
			args: args{
				got:  validJSONObjectAsString,
				want: []bool{},
			},
			want: false,
		},
		{
			name: "should fail when want valid JSON object, but get invalid JSON object",
			args: args{
				got:  invalidJSONObjectAsString,
				want: validJSONObjectAsString,
			},
			want: false,
		},
		{
			name: "should fail when want invalid JSON object, but get valid JSON object",
			args: args{
				got:  validJSONObjectAsString,
				want: invalidJSONObjectAsString,
			},
			want: false,
		},
		{
			name: "should fail when want same valid JSON object, but get another",
			args: args{
				got:  differentValidJSONObjectAsString,
				want: validJSONObjectAsString,
			},
			want: false,
		},
		{
			name: "should fail when want valid JSON object, but get empty string",
			args: args{
				got:  "",
				want: validJSONObjectAsString,
			},
			want: false,
		},
		{
			name: "should fail when want empty string, but get valid JSON object",
			args: args{
				got:  validJSONObjectAsString,
				want: "",
			},
			want: false,
		},
		{
			name: "should pass when want and get empty string",
			args: args{
				got:  "",
				want: "",
			},
			want: true,
		},
		{
			name: "should pass when want and get nil",
			args: args{
				got:  nil,
				want: nil,
			},
			want: true,
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			dummyT := &testing.T{}
			dummyAssert := New(dummyT)
			assert := New(t)

			// When
			got := dummyAssert(tt.args.got).IsJSONEqualTo(tt.args.want)

			// Then
			assert(got).Equals(tt.want)
			assert(dummyT.Failed()).Equals(!tt.want)
		})
	}
}

func TestIsWantedError(t *testing.T) {
	type args struct {
		got     interface{}
		wantErr bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should return true when wantErr and get non-nil error",
			args: args{
				got:     nonZero["error"],
				wantErr: true,
			},
			want: true,
		},
		{
			name: "should return true when not wantErr and get nil error",
			args: args{
				got:     zero["error"],
				wantErr: false,
			},
			want: true,
		},
		{
			name: "should return false when wantErr and get nil error",
			args: args{
				got:     zero["error"],
				wantErr: true,
			},
			want: false,
		},
		{
			name: "should return false when not wantErr and get non-nil error",
			args: args{
				got:     nonZero["error"],
				wantErr: false,
			},
			want: false,
		},
		{
			name: "should return false when wantErr and get non-error",
			args: args{
				got:     nonZero["struct"],
				wantErr: true,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			dummyT := &testing.T{}
			dummyAssert := New(dummyT)
			assert := New(t)
			require := NewFatal(t)

			// When
			got := dummyAssert(tt.args.got).IsWantedError(tt.args.wantErr)

			// Then
			assert(got).Equals(tt.want)
			require(dummyT.Failed()).Equals(!tt.want)
		})
	}
}
