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
		got         interface{}
		want        interface{}
		wantTFailed bool
	}{
		"should pass for identical nil args": {
			got:         nil,
			want:        nil,
			wantTFailed: false,
		},
		"should pass for equal comparable non-nil args": {
			got:         nonZero["int"],
			want:        nonZero["int"],
			wantTFailed: false,
		},
		"should fail when want non-nil, but get nil": {
			got:         nil,
			want:        nonZero["int"],
			wantTFailed: true,
		},
		"should fail when want nil, but get comparable non-nil": {
			got:         nonZero["struct"],
			want:        nil,
			wantTFailed: true,
		},
		"should fail when want comparable non-nil, but get non-comparable func": {
			got:         nonZero["func"],
			want:        nonZero["int"],
			wantTFailed: true,
		},
		"should fail when want non-comparable func, but get comparable non-nil": {
			got:         nonZero["string"],
			want:        func() {},
			wantTFailed: true,
		},
		"should fail for unequal comparable non-nil args": {
			got:         nonZero["int"],
			want:        nonZero["string"],
			wantTFailed: true,
		},
	}

	t.Parallel()
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(subT *testing.T) {
			// when
			instanceT := &testing.T{}
			assert := New(instanceT)
			assert(tc.got).Equals(tc.want)

			// then
			if tc.wantTFailed == instanceT.Failed() {
				return
			}
			subT.Errorf("expected assert(%#v).Equals(%#v) to be %v, found %v", tc.got, tc.want, !tc.wantTFailed, !instanceT.Failed())
		})
	}

}

func TestIsNil(t *testing.T) {
	testCases := map[string]struct {
		got         interface{}
		givenT      *testing.T
		wantTFailed bool
	}{
		"should fail when get non-nilable": {
			got:         nonZero["byte"],
			givenT:      &testing.T{},
			wantTFailed: true,
		},
		"should fail when get non-nil chan": {
			got:         nonZero["chan"],
			givenT:      &testing.T{},
			wantTFailed: true,
		},
		"should pass when get nil chan": {
			got:         zero["chan"],
			givenT:      &testing.T{},
			wantTFailed: false,
		},
		"should fail when get non-nil func": {
			got:         nonZero["func"],
			givenT:      &testing.T{},
			wantTFailed: true,
		},
		"should pass when get nil func": {
			got:         zero["func"],
			givenT:      &testing.T{},
			wantTFailed: false,
		},
		"should fail when get non-nil map": {
			got:         nonZero["map"],
			givenT:      &testing.T{},
			wantTFailed: true,
		},
		"should pass when get nil map": {
			got:         zero["map"],
			givenT:      &testing.T{},
			wantTFailed: false,
		},
		"should fail when get non-nil interface": {
			got:         nonZero["interface"],
			givenT:      &testing.T{},
			wantTFailed: true,
		},
		"should pass when get nil interface": {
			got:         zero["interface"],
			givenT:      &testing.T{},
			wantTFailed: false,
		},
		"should fail when get non-nil pointer": {
			got:         nonZero["ptr"],
			givenT:      &testing.T{},
			wantTFailed: true,
		},
		"should pass when get nil pointer": {
			got:         zero["ptr"],
			givenT:      &testing.T{},
			wantTFailed: false,
		},
		"should fail when get non-nil slice": {
			got:         nonZero["slice"],
			givenT:      &testing.T{},
			wantTFailed: true,
		},
		"should pass when get nil slice": {
			got:         zero["slice"],
			givenT:      &testing.T{},
			wantTFailed: false,
		},
	}

	t.Parallel()
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(subT *testing.T) {
			// when
			assert := New(tc.givenT)
			assert(tc.got).IsNil()

			gotTFailed := tc.givenT.Failed()

			// then
			if gotTFailed == tc.wantTFailed {
				return
			}
			subT.Errorf("expected assert(%#v).IsNil() to be %v, found %v", tc.got, !tc.wantTFailed, !gotTFailed)
		})
	}
}

func TestIsNotNil(t *testing.T) {
	testCases := map[string]struct {
		got         interface{}
		givenT      *testing.T
		wantTFailed bool
	}{
		"should pass when get non-nilable": {
			got:         nonZero["int"],
			givenT:      &testing.T{},
			wantTFailed: false,
		},
		"should pass when get non-nil chan": {
			got:         nonZero["chan"],
			givenT:      &testing.T{},
			wantTFailed: false,
		},
		"should fail when get nil chan": {
			got:         zero["chan"],
			givenT:      &testing.T{},
			wantTFailed: true,
		},
		"should pass when get non-nil func": {
			got:         nonZero["func"],
			givenT:      &testing.T{},
			wantTFailed: false,
		},
		"should fail when get nil func": {
			got:         zero["func"],
			givenT:      &testing.T{},
			wantTFailed: true,
		},
		"should pass when get non-nil map": {
			got:         nonZero["map"],
			givenT:      &testing.T{},
			wantTFailed: false,
		},
		"should fail when get nil map": {
			got:         zero["map"],
			givenT:      &testing.T{},
			wantTFailed: true,
		},
		"should pass when get non-nil interface": {
			got:         nonZero["interface"],
			givenT:      &testing.T{},
			wantTFailed: false,
		},
		"should fail when get nil interface": {
			got:         zero["interface"],
			givenT:      &testing.T{},
			wantTFailed: true,
		},
		"should pass when get non-nil pointer": {
			got:         nonZero["ptr"],
			givenT:      &testing.T{},
			wantTFailed: false,
		},
		"should fail when get nil pointer": {
			got:         zero["ptr"],
			givenT:      &testing.T{},
			wantTFailed: true,
		},
		"should pass when get non-nil slice": {
			got:         nonZero["slice"],
			givenT:      &testing.T{},
			wantTFailed: false,
		},
		"should fail when get nil slice": {
			got:         zero["slice"],
			givenT:      &testing.T{},
			wantTFailed: true,
		},
	}

	t.Parallel()
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(subT *testing.T) {
			// when
			assert := New(tc.givenT)
			assert(tc.got).IsNotNil()

			gotTFailed := tc.givenT.Failed()

			// then
			if gotTFailed == tc.wantTFailed {
				return
			}
			subT.Errorf("expected assert(%#v).IsNotNil() to be %v, found %v", tc.got, !tc.wantTFailed, !gotTFailed)
		})
	}
}

func TestEqualsElementsInIgnoringOrder(t *testing.T) {
	testCases := map[string]struct {
		wantTFailed bool
		want        interface{}
		got         interface{}
		givenT      *testing.T
	}{
		"should fail when get bool": {
			wantTFailed: true,
			want:        nonZero["slice"],
			got:         nonZero["bool"],
			givenT:      &testing.T{},
		},
		"should fail when get byte": {
			wantTFailed: true,
			want:        nonZero["slice"],
			got:         nonZero["byte"],
			givenT:      &testing.T{},
		},
		"should fail when get chan": {
			wantTFailed: true,
			want:        nonZero["slice"],
			got:         nonZero["chan"],
			givenT:      &testing.T{},
		},
		"should fail when get complex128": {
			wantTFailed: true,
			want:        nonZero["slice"],
			got:         nonZero["complex128"],
			givenT:      &testing.T{},
		},
		"should fail when get complex64": {
			wantTFailed: true,
			want:        nonZero["slice"],
			got:         nonZero["complex64"],
			givenT:      &testing.T{},
		},
		"should fail when get float32": {
			wantTFailed: true,
			want:        nonZero["slice"],
			got:         nonZero["float32"],
			givenT:      &testing.T{},
		},
		"should fail when get float64": {
			wantTFailed: true,
			want:        nonZero["slice"],
			got:         nonZero["float64"],
			givenT:      &testing.T{},
		},
		"should fail when get func": {
			wantTFailed: true,
			want:        nonZero["slice"],
			got:         nonZero["func"],
			givenT:      &testing.T{},
		},
		"should fail when get int": {
			wantTFailed: true,
			want:        nonZero["slice"],
			got:         nonZero["int"],
			givenT:      &testing.T{},
		},
		"should fail when get int8": {
			wantTFailed: true,
			want:        nonZero["slice"],
			got:         nonZero["int8"],
			givenT:      &testing.T{},
		},
		"should fail when get int16": {
			wantTFailed: true,
			want:        nonZero["slice"],
			got:         nonZero["int16"],
			givenT:      &testing.T{},
		},
		"should fail when get int32": {
			wantTFailed: true,
			want:        nonZero["slice"],
			got:         nonZero["int32"],
			givenT:      &testing.T{},
		},
		"should fail when get int64": {
			wantTFailed: true,
			want:        nonZero["slice"],
			got:         nonZero["int64"],
			givenT:      &testing.T{},
		},
		"should fail when get map": {
			wantTFailed: true,
			want:        nonZero["slice"],
			got:         nonZero["map"],
			givenT:      &testing.T{},
		},
		"should fail when get ptr": {
			wantTFailed: true,
			want:        nonZero["slice"],
			got:         nonZero["ptr"],
			givenT:      &testing.T{},
		},
		"should fail when get string": {
			wantTFailed: true,
			want:        nonZero["slice"],
			got:         nonZero["string"],
			givenT:      &testing.T{},
		},
		"should fail when get struct": {
			wantTFailed: true,
			want:        nonZero["slice"],
			got:         nonZero["struct"],
			givenT:      &testing.T{},
		},
		"should pass when want slice and get the same set of elements in order": {
			wantTFailed: false,
			want:        []string{"a", "b", "c"},
			got:         []string{"a", "b", "c"},
			givenT:      &testing.T{},
		},
		"should pass when want slice and get the same set of elements not in order": {
			wantTFailed: false,
			want:        []string{"a", "b", "c"},
			got:         []string{"b", "c", "a"},
			givenT:      &testing.T{},
		},
		"should fail when want slice but get subset not in order": {
			wantTFailed: true,
			want:        []string{"a", "b", "c"},
			got:         []string{"c", "a"},
			givenT:      &testing.T{},
		},
		"should fail when want slice but get subset in order": {
			wantTFailed: true,
			want:        []string{"a", "b", "c"},
			got:         []string{"b", "c"},
			givenT:      &testing.T{},
		},
		"should pass when want array and get the same set of elements in order": {
			wantTFailed: false,
			want:        [3]string{"a", "b", "c"},
			got:         [3]string{"a", "b", "c"},
			givenT:      &testing.T{},
		},
		"should pass when want array and get the same set of elements not in order": {
			wantTFailed: false,
			want:        [3]string{"a", "b", "c"},
			got:         [3]string{"b", "c", "a"},
			givenT:      &testing.T{},
		},
		"should fail when want array but get subset not in order": {
			wantTFailed: true,
			want:        [3]string{"a", "b", "c"},
			got:         [2]string{"c", "a"},
			givenT:      &testing.T{},
		},
		"should fail when want array but get subset in order": {
			wantTFailed: true,
			want:        [3]string{"a", "b", "c"},
			got:         [2]string{"b", "c"},
			givenT:      &testing.T{},
		},
		"should pass when want array but get slice with same set of elements": {
			wantTFailed: false,
			want:        [3]string{"a", "b", "c"},
			got:         []string{"a", "b", "c"},
			givenT:      &testing.T{},
		},
		"should pass when want slice but get array with same set of elements": {
			wantTFailed: false,
			want:        []string{"a", "b", "c"},
			got:         [3]string{"a", "b", "c"},
			givenT:      &testing.T{},
		},
		"should pass when want and get empty slice": {
			wantTFailed: false,
			want:        []string{},
			got:         []string{},
			givenT:      &testing.T{},
		},
	}

	t.Parallel()
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(subT *testing.T) {
			// when
			assert := New(tc.givenT)
			assert(tc.got).EqualsElementsInIgnoringOrder(tc.want)

			gotTFailed := tc.givenT.Failed()

			// then
			if gotTFailed == tc.wantTFailed {
				return
			}
			subT.Errorf("expected assert(%#v).EqualsElementsInIgnoringOrder(%#v) to be %v, found %v", tc.got, tc.want, !tc.wantTFailed, !gotTFailed)
		})
	}
}
