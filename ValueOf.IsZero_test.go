package tinyreflect_test

import (
	"testing"

	"github.com/cdvelop/tinyreflect"
)

// TestIsZero tests the IsZero method for all supported types.
func TestIsZero(t *testing.T) {
	tr := tinyreflect.New()

	type SimpleStruct struct{ A int; B string }
	type NestedStruct struct{ S SimpleStruct; C bool }

	tests := []struct {
		name     string
		value    any
		expected bool
	}{
		// Basic types
		{"string zero", "", true},
		{"string non-zero", "hello", false},
		{"bool zero", false, true},
		{"bool non-zero", true, false},
		{"int zero", 0, true},
		{"int non-zero", 42, false},
		{"uint zero", uint(0), true},
		{"float64 zero", float64(0.0), true},

		// Pointer types
		{"nil pointer", (*int)(nil), true},
		{"non-nil pointer", func() *int { i := 42; return &i }(), false},

		// Slice types
		{"nil slice", (func() []int { var s []int; return s })(), true},
		{"empty slice", []int{}, false},
		{"non-empty slice", []int{1, 2, 3}, false},

		// Map types
		{"nil map", (func() map[string]int { var m map[string]int; return m })(), true},
		{"empty map", map[string]int{}, false},
		{"non-empty map", map[string]int{"key": 42}, false},

		// Struct types
		{"struct all zero", SimpleStruct{A: 0, B: ""}, true},
		{"struct A non-zero", SimpleStruct{A: 42, B: ""}, false},
		{"nested struct all zero", NestedStruct{S: SimpleStruct{A: 0, B: ""}, C: false}, true},
		{"nested struct S.A non-zero", NestedStruct{S: SimpleStruct{A: 42, B: ""}, C: false}, false},

		// Interface types
		{"nil interface", nil, true},
		{"interface with zero int", (interface{})(0), true},
		{"interface with non-zero int", (interface{})(42), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := tr.ValueOf(tt.value)
			result := v.IsZero()
			if result != tt.expected {
				t.Errorf("IsZero() = %v, expected %v for %T(%#v)", result, tt.expected, tt.value, tt.value)
			}
		})
	}
}

// TestIsZero_ZeroValue tests the IsZero method on a zero Value struct.
func TestIsZero_ZeroValue(t *testing.T) {
	var v tinyreflect.Value // Zero value of the struct itself
	if !v.IsZero() {
		t.Error("IsZero() on a zero Value struct should be true")
	}
}

// BenchmarkIsZero benchmarks the IsZero method.
func BenchmarkIsZero(b *testing.B) {
	tr := tinyreflect.New()
	testCases := []struct {
		name  string
		value any
	}{
		{"Int", 0},
		{"String", ""},
		{"Bool", false},
		{"Float", 3.14},
		{"Slice", []int{1, 2}},
		{"Map", map[string]int{"key": 1}},
		{"Struct", struct{ A int }{A: 42}},
	}

	for _, tc := range testCases {
		v := tr.ValueOf(tc.value)
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = v.IsZero()
			}
		})
	}
}