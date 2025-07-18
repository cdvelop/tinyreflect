package tinyreflect

import (
	"testing"
)

// TestIndirectTypeMethod tests that Indirect + Type operations work correctly
// This ensures the core reflection functionality works for various value types
func TestIndirectTypeMethod(t *testing.T) {
	t.Run("rv.Type() should not return nil", func(t *testing.T) {
		testCases := []interface{}{
			42,
			"hello",
			true,
			[]int{1, 2, 3},
			struct{ X int }{42},
			&struct{ X int }{42},
		}

		for _, tc := range testCases {
			// Test the core Indirect + ValueOf + Type pattern
			rv := Indirect(ValueOf(tc))
			typ := rv.Type()

			if typ == nil {
				t.Errorf("rv.Type() returned nil for value %v", tc)
			} else {
				t.Logf("rv.Type() returned %p for value %v", typ, tc)
			}
		}
	})

	t.Run("Test specific struct case", func(t *testing.T) {
		// Test a common struct scenario
		x := struct {
			Name string
			Age  int
		}{"John", 25}

		// Test direct struct
		rv := Indirect(ValueOf(x))
		typ := rv.Type()

		if typ == nil {
			t.Error("Struct case: rv.Type() returned nil")
		} else {
			t.Logf("Struct case: rv.Type() returned %p", typ)
		}

		// Test pointer to struct
		rv2 := Indirect(ValueOf(&x))
		typ2 := rv2.Type()

		if typ2 == nil {
			t.Error("Pointer to struct case: rv.Type() returned nil")
		} else {
			t.Logf("Pointer to struct case: rv.Type() returned %p", typ2)
		}
	})
}
