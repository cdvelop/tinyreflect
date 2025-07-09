package tinyreflect

import (
	"testing"

	. "github.com/cdvelop/tinystring"
)

// TestEssentialFunctionality verifies all core supported types work correctly
// This test covers the exact types listed as supported in README.md
func TestEssentialFunctionality(t *testing.T) {
	t.Run("Basic Types", func(t *testing.T) {
		// Test all basic types from README.md
		tests := []struct {
			name     string
			value    any
			expected Kind
		}{
			{"string", "hello", KString},
			{"bool", true, KBool},
			{"int", int(42), KInt},
			{"int8", int8(8), KInt8},
			{"int16", int16(16), KInt16},
			{"int32", int32(32), KInt32},
			{"int64", int64(64), KInt64},
			{"uint", uint(42), KUint},
			{"uint8", uint8(8), KUint8},
			{"uint16", uint16(16), KUint16},
			{"uint32", uint32(32), KUint32},
			{"uint64", uint64(64), KUint64},
			{"float32", float32(3.14), KFloat32},
			{"float64", float64(3.14), KFloat64},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				v := refValueOf(test.value)
				if v.err != nil {
					t.Errorf("Error for %s: %v", test.name, v.err)
				}
				if v.refKind() != test.expected {
					t.Errorf("%s: got Kind %v, want %v", test.name, v.refKind(), test.expected)
				}
			})
		}
	})

	t.Run("All Basic Slices", func(t *testing.T) {
		// Test all slice types from README.md
		tests := []struct {
			name     string
			value    any
			expected Kind
			length   int
		}{
			{"[]string", []string{"a", "b"}, KSliceStr, 2},
			{"[]bool", []bool{true, false}, KSlice, 2},
			{"[]byte", []byte{1, 2, 3}, KByte, 3},
			{"[]int", []int{1, 2}, KSlice, 2},
			{"[]int8", []int8{1, 2}, KSlice, 2},
			{"[]int16", []int16{1, 2}, KSlice, 2},
			{"[]int32", []int32{1, 2}, KSlice, 2},
			{"[]int64", []int64{1, 2}, KSlice, 2},
			{"[]uint", []uint{1, 2}, KSlice, 2},
			{"[]uint8", []uint8{1, 2}, KByte, 2}, // []uint8 is equivalent to []byte
			{"[]uint16", []uint16{1, 2}, KSlice, 2},
			{"[]uint32", []uint32{1, 2}, KSlice, 2},
			{"[]uint64", []uint64{1, 2}, KSlice, 2},
			{"[]float32", []float32{1.1, 2.2}, KSlice, 2},
			{"[]float64", []float64{1.1, 2.2}, KSlice, 2},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				v := refValueOf(test.value)
				if v.err != nil {
					t.Errorf("Error for %s: %v", test.name, v.err)
				}
				if v.refKind() != test.expected {
					t.Errorf("%s: got Kind %v, want %v", test.name, v.refKind(), test.expected)
				}
				if v.refLen() != test.length {
					t.Errorf("%s: got length %d, want %d", test.name, v.refLen(), test.length)
				}

				// Test element access
				elem := v.refIndex(0)
				if elem.err != nil {
					t.Errorf("%s: element access error: %v", test.name, elem.err)
				}
			})
		}
	})

	t.Run("Structs with Supported Fields", func(t *testing.T) {
		// Test struct with all supported field types
		type TestStruct struct {
			Name   string
			Age    int
			Height float64
			Active bool
		}

		s := TestStruct{
			Name:   "John",
			Age:    30,
			Height: 180.5,
			Active: true,
		}

		v := refValueOf(s)
		if v.err != nil {
			t.Errorf("Struct error: %v", v.err)
		}
		if v.refKind() != KStruct {
			t.Errorf("Struct: got Kind %v, want %v", v.refKind(), KStruct)
		}

		// Verify we can access fields (basic functionality)
		numFields := v.refNumField()
		if numFields != 4 {
			t.Errorf("Struct: got %d fields, want 4", numFields)
		}
	})

	t.Run("Maps with Supported Types", func(t *testing.T) {
		tests := []struct {
			name  string
			value any
		}{
			{"map[string]string", map[string]string{"key": "value"}},
			{"map[string]int", map[string]int{"key": 42}},
			{"map[int]string", map[int]string{1: "value"}},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				v := refValueOf(test.value)
				if v.err != nil {
					t.Errorf("Error for %s: %v", test.name, v.err)
				}
				if v.refKind() != KMap {
					t.Errorf("%s: got Kind %v, want %v", test.name, v.refKind(), KMap)
				}
			})
		}
	})

	t.Run("Pointers to Supported Types", func(t *testing.T) {
		str := "hello"
		num := 42
		flag := true
		pi := 3.14

		tests := []struct {
			name  string
			value any
		}{
			{"*string", &str},
			{"*int", &num},
			{"*bool", &flag},
			{"*float64", &pi},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				v := refValueOf(test.value)
				if v.err != nil {
					t.Errorf("Error for %s: %v", test.name, v.err)
				}
				if v.refKind() != KPointer {
					t.Errorf("%s: got Kind %v, want %v", test.name, v.refKind(), KPointer)
				}

				// Test pointer dereferencing
				elem := v.refElem()
				if elem.err != nil {
					t.Errorf("%s: refElem error: %v", test.name, elem.err)
				}
			})
		}
	})

	t.Run("Unsupported Types Return Error", func(t *testing.T) {
		// These should return errors according to README.md
		unsupported := []any{
			make(chan int),     // chan not supported
			func() {},          // func not supported
			complex(1, 2),      // complex64 not supported
			complex128(1 + 2i), // complex128 not supported
			uintptr(0x123),     // uintptr not supported (internal use only)
			[3]int{1, 2, 3},    // arrays not supported (only slices)
		}

		for i, value := range unsupported {
			v := refValueOf(value)
			if v.err == nil {
				t.Errorf("Unsupported type %d should return error, but didn't: %T", i, value)
			}
		}
	})
}
