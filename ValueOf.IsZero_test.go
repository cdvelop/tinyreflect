package tinyreflect

import (
	"testing"
)

// TestIsZero tests the IsZero method for all supported types
func TestIsZero(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		expected bool
	}{
		// String tests
		{"string zero", "", true},
		{"string non-zero", "hello", false},

		// Bool tests
		{"bool zero", false, true},
		{"bool non-zero", true, false},

		// Integer tests
		{"int zero", 0, true},
		{"int non-zero", 42, false},
		{"int8 zero", int8(0), true},
		{"int8 non-zero", int8(42), false},
		{"int16 zero", int16(0), true},
		{"int16 non-zero", int16(42), false},
		{"int32 zero", int32(0), true},
		{"int32 non-zero", int32(42), false},
		{"int64 zero", int64(0), true},
		{"int64 non-zero", int64(42), false},

		// Unsigned integer tests
		{"uint zero", uint(0), true},
		{"uint non-zero", uint(42), false},
		{"uint8 zero", uint8(0), true},
		{"uint8 non-zero", uint8(42), false},
		{"uint16 zero", uint16(0), true},
		{"uint16 non-zero", uint16(42), false},
		{"uint32 zero", uint32(0), true},
		{"uint32 non-zero", uint32(42), false},
		{"uint64 zero", uint64(0), true},
		{"uint64 non-zero", uint64(42), false},
		{"uintptr zero", uintptr(0), true},
		{"uintptr non-zero", uintptr(42), false},

		// Float tests
		{"float32 zero", float32(0.0), true},
		{"float32 non-zero", float32(3.14), false},
		{"float64 zero", float64(0.0), true},
		{"float64 non-zero", float64(3.14), false},

		// Pointer tests
		{"nil pointer", (*int)(nil), true},
		{"non-nil pointer", func() *int { i := 42; return &i }(), false},

		// Slice tests
		{"nil slice", func() []int { var s []int; return s }(), true},
		{"empty slice", []int{}, false}, // Empty slice is not zero value
		{"non-empty slice", []int{1, 2, 3}, false},

		// Map tests
		{"nil map", func() map[string]int { var m map[string]int; return m }(), true},
		{"empty map", map[string]int{}, false}, // Empty map is not zero value
		{"non-empty map", map[string]int{"key": 42}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := ValueOf(tt.value)
			result := v.IsZero()
			if result != tt.expected {
				t.Errorf("IsZero() = %v, expected %v for %T(%v)", result, tt.expected, tt.value, tt.value)
			}
		})
	}
}

// TestIsZeroStruct tests IsZero for struct types
func TestIsZeroStruct(t *testing.T) {
	// Define test structs
	type SimpleStruct struct {
		A int
		B string
	}

	type NestedStruct struct {
		S SimpleStruct
		C bool
	}

	tests := []struct {
		name     string
		value    any
		expected bool
	}{
		// Simple struct tests
		{"struct all zero", SimpleStruct{A: 0, B: ""}, true},
		{"struct A non-zero", SimpleStruct{A: 42, B: ""}, false},
		{"struct B non-zero", SimpleStruct{A: 0, B: "hello"}, false},
		{"struct both non-zero", SimpleStruct{A: 42, B: "hello"}, false},

		// Nested struct tests
		{"nested struct all zero", NestedStruct{S: SimpleStruct{A: 0, B: ""}, C: false}, true},
		{"nested struct S.A non-zero", NestedStruct{S: SimpleStruct{A: 42, B: ""}, C: false}, false},
		{"nested struct S.B non-zero", NestedStruct{S: SimpleStruct{A: 0, B: "hello"}, C: false}, false},
		{"nested struct C non-zero", NestedStruct{S: SimpleStruct{A: 0, B: ""}, C: true}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := ValueOf(tt.value)
			result := v.IsZero()
			if result != tt.expected {
				t.Errorf("IsZero() = %v, expected %v for %T%+v", result, tt.expected, tt.value, tt.value)
			}
		})
	}
}

// TestIsZeroEdgeCases tests edge cases for IsZero
func TestIsZeroEdgeCases(t *testing.T) {
	// Test with interface{} containing zero values
	var i any = 0
	v := ValueOf(i)
	if !v.IsZero() {
		t.Errorf("IsZero() should return true for interface{} containing 0")
	}

	i = ""
	v = ValueOf(i)
	if !v.IsZero() {
		t.Errorf("IsZero() should return true for interface{} containing empty string")
	}

	i = 42
	v = ValueOf(i)
	if v.IsZero() {
		t.Errorf("IsZero() should return false for interface{} containing 42")
	}

	// Test with nil interface
	var nilInterface any
	v = ValueOf(nilInterface)
	if !v.IsZero() {
		t.Errorf("IsZero() should return true for nil interface{}")
	}
}

// TestIsZeroUnsupportedTypes tests behavior with unsupported types
func TestIsZeroUnsupportedTypes(t *testing.T) {
	// For unsupported types, we can't easily create test values without internal manipulation
	// The default case in IsZero returns false, which is the safe behavior
	// This test ensures the method doesn't panic on edge cases

	// Test with a zero Value (should be true)
	v := Value{}
	result := v.IsZero()
	if !result {
		t.Errorf("IsZero() should return true for zero Value, got false")
	}
}

// BenchmarkIsZero benchmarks the IsZero method
func BenchmarkIsZero(b *testing.B) {
	// Benchmark with different types
	testValues := []any{
		0,                        // int
		"",                       // string
		false,                    // bool
		3.14,                     // float64
		[]int{1, 2},              // slice
		map[string]int{"key": 1}, // map
		struct{ A int }{A: 42},   // struct
	}

	for _, val := range testValues {
		v := ValueOf(val)
		b.Run("IsZero", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = v.IsZero()
			}
		})
	}
}
