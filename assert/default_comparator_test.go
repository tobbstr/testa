package assert

import (
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
	}
)

func TestEquals(t *testing.T) {
	// given
	testCases := map[string]struct {
		got              interface{}
		want             interface{}
		wantAssertPassed bool
	}{
		"should pass for identical nil args": {
			got:              nil,
			want:             nil,
			wantAssertPassed: true,
		},
		"should pass for equal comparable non-nil args": {
			got:              nonZero["int"],
			want:             nonZero["int"],
			wantAssertPassed: true,
		},
		"should fail when want non-nil, but get nil": {
			got:              nil,
			want:             nonZero["int"],
			wantAssertPassed: false,
		},
		"should fail when want nil, but get comparable non-nil": {
			got:              nonZero["struct"],
			want:             nil,
			wantAssertPassed: false,
		},
		"should fail when want comparable non-nil, but get non-comparable func": {
			got:              nonZero["func"],
			want:             nonZero["int"],
			wantAssertPassed: false,
		},
		"should fail when want non-comparable func, but get comparable non-nil": {
			got:              nonZero["string"],
			want:             func() {},
			wantAssertPassed: false,
		},
		"should fail for unequal comparable non-nil args": {
			got:              nonZero["int"],
			want:             nonZero["string"],
			wantAssertPassed: false,
		},
	}

	t.Parallel()
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(subT *testing.T) {
			// given
			givenT := &testing.T{}
			assert := New(givenT)

			// when
			assert(tc.got).Equals(tc.want)

			// then
			gotAssertPassed := !givenT.Failed()
			if gotAssertPassed == tc.wantAssertPassed {
				return
			}
			subT.Errorf("expected assert(%#v).Equals(%#v) to be %v, found %v", tc.got, tc.want, tc.wantAssertPassed, gotAssertPassed)
		})
	}
}

func TestEqualsElementsInIgnoringOrder(t *testing.T) {
	testCases := map[string]struct {
		wantAssertPassed bool
		want             interface{}
		got              interface{}
		givenT           *testing.T
	}{
		"should fail when get bool": {
			wantAssertPassed: false,
			want:             nonZero["slice"],
			got:              nonZero["bool"],
			givenT:           &testing.T{},
		},
		"should fail when get byte": {
			wantAssertPassed: false,
			want:             nonZero["slice"],
			got:              nonZero["byte"],
			givenT:           &testing.T{},
		},
		"should fail when get chan": {
			wantAssertPassed: false,
			want:             nonZero["slice"],
			got:              nonZero["chan"],
			givenT:           &testing.T{},
		},
		"should fail when get complex128": {
			wantAssertPassed: false,
			want:             nonZero["slice"],
			got:              nonZero["complex128"],
			givenT:           &testing.T{},
		},
		"should fail when get complex64": {
			wantAssertPassed: false,
			want:             nonZero["slice"],
			got:              nonZero["complex64"],
			givenT:           &testing.T{},
		},
		"should fail when get float32": {
			wantAssertPassed: false,
			want:             nonZero["slice"],
			got:              nonZero["float32"],
			givenT:           &testing.T{},
		},
		"should fail when get float64": {
			wantAssertPassed: false,
			want:             nonZero["slice"],
			got:              nonZero["float64"],
			givenT:           &testing.T{},
		},
		"should fail when get func": {
			wantAssertPassed: false,
			want:             nonZero["slice"],
			got:              nonZero["func"],
			givenT:           &testing.T{},
		},
		"should fail when get int": {
			wantAssertPassed: false,
			want:             nonZero["slice"],
			got:              nonZero["int"],
			givenT:           &testing.T{},
		},
		"should fail when get int8": {
			wantAssertPassed: false,
			want:             nonZero["slice"],
			got:              nonZero["int8"],
			givenT:           &testing.T{},
		},
		"should fail when get int16": {
			wantAssertPassed: false,
			want:             nonZero["slice"],
			got:              nonZero["int16"],
			givenT:           &testing.T{},
		},
		"should fail when get int32": {
			wantAssertPassed: false,
			want:             nonZero["slice"],
			got:              nonZero["int32"],
			givenT:           &testing.T{},
		},
		"should fail when get int64": {
			wantAssertPassed: false,
			want:             nonZero["slice"],
			got:              nonZero["int64"],
			givenT:           &testing.T{},
		},
		"should fail when get map": {
			wantAssertPassed: false,
			want:             nonZero["slice"],
			got:              nonZero["map"],
			givenT:           &testing.T{},
		},
		"should fail when get ptr": {
			wantAssertPassed: false,
			want:             nonZero["slice"],
			got:              nonZero["ptr"],
			givenT:           &testing.T{},
		},
		"should fail when get string": {
			wantAssertPassed: false,
			want:             nonZero["slice"],
			got:              nonZero["string"],
			givenT:           &testing.T{},
		},
		"should fail when get struct": {
			wantAssertPassed: false,
			want:             nonZero["slice"],
			got:              nonZero["struct"],
			givenT:           &testing.T{},
		},
		"should pass when want slice and get the same set of elements in order": {
			wantAssertPassed: true,
			want:             []string{"a", "b", "c"},
			got:              []string{"a", "b", "c"},
			givenT:           &testing.T{},
		},
		"should pass when want slice and get the same set of elements not in order": {
			wantAssertPassed: true,
			want:             []string{"a", "b", "c"},
			got:              []string{"b", "c", "a"},
			givenT:           &testing.T{},
		},
		"should fail when want slice but get subset not in order": {
			wantAssertPassed: false,
			want:             []string{"a", "b", "c"},
			got:              []string{"c", "a"},
			givenT:           &testing.T{},
		},
		"should fail when want slice but get subset in order": {
			wantAssertPassed: false,
			want:             []string{"a", "b", "c"},
			got:              []string{"b", "c"},
			givenT:           &testing.T{},
		},
		"should pass when want array and get the same set of elements in order": {
			wantAssertPassed: true,
			want:             [3]string{"a", "b", "c"},
			got:              [3]string{"a", "b", "c"},
			givenT:           &testing.T{},
		},
		"should pass when want array and get the same set of elements not in order": {
			wantAssertPassed: true,
			want:             [3]string{"a", "b", "c"},
			got:              [3]string{"b", "c", "a"},
			givenT:           &testing.T{},
		},
		"should fail when want array but get subset not in order": {
			wantAssertPassed: false,
			want:             [3]string{"a", "b", "c"},
			got:              [2]string{"c", "a"},
			givenT:           &testing.T{},
		},
		"should fail when want array but get subset in order": {
			wantAssertPassed: false,
			want:             [3]string{"a", "b", "c"},
			got:              [2]string{"b", "c"},
			givenT:           &testing.T{},
		},
		"should pass when want array but get slice with same set of elements": {
			wantAssertPassed: true,
			want:             [3]string{"a", "b", "c"},
			got:              []string{"a", "b", "c"},
			givenT:           &testing.T{},
		},
		"should pass when want slice but get array with same set of elements": {
			wantAssertPassed: true,
			want:             []string{"a", "b", "c"},
			got:              [3]string{"a", "b", "c"},
			givenT:           &testing.T{},
		},
		"should pass when want and get empty slice": {
			wantAssertPassed: true,
			want:             []string{},
			got:              []string{},
			givenT:           &testing.T{},
		},
	}

	t.Parallel()
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(subT *testing.T) {
			// given
			assert := New(tc.givenT)

			// when
			assert(tc.got).EqualsElementsInIgnoringOrder(tc.want)

			// then
			gotAssertPassed := !tc.givenT.Failed()
			if gotAssertPassed == tc.wantAssertPassed {
				return
			}
			subT.Errorf("expected assert(%#v).EqualsElementsInIgnoringOrder(%#v) to be %v, found %v", tc.got, tc.want, tc.wantAssertPassed, gotAssertPassed)
		})
	}
}

func TestIsEmpty(t *testing.T) {
	nonEmptyChan := make(chan int, 10)
	nonEmptyChan <- 4
	testcases := map[string]struct {
		got              interface{}
		givenT           *testing.T
		wantAssertPassed bool
	}{
		"should pass when get empty array": {
			got:              [0]int{},
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get non-empty array": {
			got:              zero["array"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get zero-valued bool": {
			got:              zero["bool"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get true bool": {
			got:              true,
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get zero-valued byte": {
			got:              zero["byte"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get non-zero-valued byte": {
			got:              nonZero["byte"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get empty chan": {
			got:              zero["chan"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get non-empty chan": {
			got:              nonEmptyChan,
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get zero-valued complex64": {
			got:              zero["complex64"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get non-zero-valued complex64": {
			got:              nonZero["complex64"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get zero-valued complex128": {
			got:              zero["complex128"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get non-zero-valued complex128": {
			got:              nonZero["complex128"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get zero-valued float32": {
			got:              zero["float32"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get non-zero-valued float32": {
			got:              nonZero["float32"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get zero-valued func": {
			got:              zero["func"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get non-zero-valued func": {
			got:              nonZero["func"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get zero-valued int": {
			got:              zero["int"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get non-zero-valued int": {
			got:              nonZero["int"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get zero-valued interface": {
			got:              zero["interface"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get non-zero-valued interface": {
			got:              nonZero["interface"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get zero-valued map": {
			got:              zero["map"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get non-zero-valued map": {
			got:              nonZero["map"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get zero-valued ptr": {
			got:              zero["ptr"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get non-zero-valued ptr": {
			got:              nonZero["ptr"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get nil slice": {
			got:              zero["slice"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should pass when get empty slice": {
			got:              []int{},
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get non-zero-valued slice": {
			got:              nonZero["slice"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get zero-valued string": {
			got:              zero["string"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get non-zero-valued string": {
			got:              nonZero["string"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get zero-valued struct": {
			got:              zero["struct"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get non-zero-valued struct": {
			got:              nonZero["struct"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
	}

	t.Parallel()
	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(subT *testing.T) {
			// given
			assert := New(tc.givenT)

			// when
			assert(tc.got).IsEmpty()

			// then
			gotAssertPassed := !tc.givenT.Failed()
			if gotAssertPassed == tc.wantAssertPassed {
				return
			}
			subT.Errorf("expected assert(%#v).IsEmpty() to be %v, found %v", tc.got, tc.wantAssertPassed, gotAssertPassed)
		})
	}
}

func TestIsNil(t *testing.T) {
	testCases := map[string]struct {
		got              interface{}
		givenT           *testing.T
		wantAssertPassed bool
	}{
		"should fail when get non-nilable": {
			got:              nonZero["byte"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should fail when get non-nil chan": {
			got:              nonZero["chan"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get nil chan": {
			got:              zero["chan"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get non-nil func": {
			got:              nonZero["func"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get nil func": {
			got:              zero["func"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get non-nil map": {
			got:              nonZero["map"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get nil map": {
			got:              zero["map"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get non-nil interface": {
			got:              nonZero["interface"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get nil interface": {
			got:              zero["interface"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get non-nil pointer": {
			got:              nonZero["ptr"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get nil pointer": {
			got:              zero["ptr"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get non-nil slice": {
			got:              nonZero["slice"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get nil slice": {
			got:              zero["slice"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
	}

	t.Parallel()
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(subT *testing.T) {
			// given
			assert := New(tc.givenT)

			// when
			assert(tc.got).IsNil()

			// then
			gotAssertPassed := !tc.givenT.Failed()
			if gotAssertPassed == tc.wantAssertPassed {
				return
			}
			subT.Errorf("expected assert(%#v).IsNil() to be %v, found %v", tc.got, tc.wantAssertPassed, gotAssertPassed)
		})
	}
}

func TestIsNotEmpty(t *testing.T) {
	nonEmptyChan := make(chan int, 10)
	nonEmptyChan <- 4
	testcases := map[string]struct {
		got              interface{}
		givenT           *testing.T
		wantAssertPassed bool
	}{
		"should fail when get empty array": {
			got:              [0]int{},
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get non-empty array": {
			got:              zero["array"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get zero-valued bool": {
			got:              zero["bool"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get true bool": {
			got:              true,
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get zero-valued byte": {
			got:              zero["byte"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get non-zero-valued byte": {
			got:              nonZero["byte"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get empty chan": {
			got:              zero["chan"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get non-empty chan": {
			got:              nonEmptyChan,
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get zero-valued complex64": {
			got:              zero["complex64"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get non-zero-valued complex64": {
			got:              nonZero["complex64"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get zero-valued complex128": {
			got:              zero["complex128"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get non-zero-valued complex128": {
			got:              nonZero["complex128"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get zero-valued float32": {
			got:              zero["float32"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get non-zero-valued float32": {
			got:              nonZero["float32"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get zero-valued func": {
			got:              zero["func"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get non-zero-valued func": {
			got:              nonZero["func"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get zero-valued int": {
			got:              zero["int"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get non-zero-valued int": {
			got:              nonZero["int"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get zero-valued interface": {
			got:              zero["interface"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get non-zero-valued interface": {
			got:              nonZero["interface"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get zero-valued map": {
			got:              zero["map"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get non-zero-valued map": {
			got:              nonZero["map"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get zero-valued ptr": {
			got:              zero["ptr"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get non-zero-valued ptr": {
			got:              nonZero["ptr"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get nil slice": {
			got:              zero["slice"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should fail when get empty slice": {
			got:              []int{},
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get non-zero-valued slice": {
			got:              nonZero["slice"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get zero-valued string": {
			got:              zero["string"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get non-zero-valued string": {
			got:              nonZero["string"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get zero-valued struct": {
			got:              zero["struct"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get non-zero-valued struct": {
			got:              nonZero["struct"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
	}

	t.Parallel()
	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(subT *testing.T) {
			// given
			assert := New(tc.givenT)

			// when
			assert(tc.got).IsNotEmpty()

			// then
			gotAssertPassed := !tc.givenT.Failed()
			if gotAssertPassed == tc.wantAssertPassed {
				return
			}
			subT.Errorf("expected assert(%#v).IsNotEmpty() to be %v, found %v", tc.got, tc.wantAssertPassed, gotAssertPassed)
		})
	}
}

func TestIsNotNil(t *testing.T) {
	testCases := map[string]struct {
		got              interface{}
		givenT           *testing.T
		wantAssertPassed bool
	}{
		"should pass when get non-nilable": {
			got:              nonZero["int"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should pass when get non-nil chan": {
			got:              nonZero["chan"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get nil chan": {
			got:              zero["chan"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get non-nil func": {
			got:              nonZero["func"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get nil func": {
			got:              zero["func"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get non-nil map": {
			got:              nonZero["map"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get nil map": {
			got:              zero["map"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get non-nil interface": {
			got:              nonZero["interface"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get nil interface": {
			got:              zero["interface"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get non-nil pointer": {
			got:              nonZero["ptr"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get nil pointer": {
			got:              zero["ptr"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get non-nil slice": {
			got:              nonZero["slice"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get nil slice": {
			got:              zero["slice"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
	}

	t.Parallel()
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(subT *testing.T) {
			// given
			assert := New(tc.givenT)

			// when
			assert(tc.got).IsNotNil()

			// then
			gotAssertPassed := !tc.givenT.Failed()
			if gotAssertPassed == tc.wantAssertPassed {
				return
			}
			subT.Errorf("expected assert(%#v).IsNotNil() to be %v, found %v", tc.got, tc.wantAssertPassed, gotAssertPassed)
		})
	}
}

func TestIsTrue(t *testing.T) {
	testCases := map[string]struct {
		got              interface{}
		givenT           *testing.T
		wantAssertPassed bool
	}{
		"should fail when get nil": {
			got:              nil,
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should fail when get non-bool": {
			got:              nonZero["int"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get true bool": {
			got:              nonZero["bool"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
		"should fail when get false bool": {
			got:              zero["bool"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
	}

	t.Parallel()
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(subT *testing.T) {
			// given
			assert := New(tc.givenT)

			// when
			assert(tc.got).IsTrue()

			// then
			gotAssertPassed := !tc.givenT.Failed()
			if gotAssertPassed == tc.wantAssertPassed {
				return
			}
			subT.Errorf("expected assert(%#v).IsTrue() to be %v, found %v", tc.got, tc.wantAssertPassed, gotAssertPassed)
		})
	}
}

func TestIsFalse(t *testing.T) {
	testCases := map[string]struct {
		got              interface{}
		givenT           *testing.T
		wantAssertPassed bool
	}{
		"should fail when get nil": {
			got:              nil,
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should fail when get non-bool": {
			got:              nonZero["int"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should fail when get true bool": {
			got:              nonZero["bool"],
			givenT:           &testing.T{},
			wantAssertPassed: false,
		},
		"should pass when get false bool": {
			got:              zero["bool"],
			givenT:           &testing.T{},
			wantAssertPassed: true,
		},
	}

	t.Parallel()
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(subT *testing.T) {
			// given
			assert := New(tc.givenT)

			// when
			assert(tc.got).IsFalse()

			// then
			gotAssertPassed := !tc.givenT.Failed()
			if gotAssertPassed == tc.wantAssertPassed {
				return
			}
			subT.Errorf("expected assert(%#v).IsFalse() to be %v, found %v", tc.got, tc.wantAssertPassed, gotAssertPassed)
		})
	}
}
