package tinyreflect_test

import (
	"testing"

	"github.com/cdvelop/tinyreflect"
)

// TestIndirectTypeMethod tests that Indirect + Type operations work correctly.
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
			rv := tinyreflect.Indirect(tinyreflect.ValueOf(tc))
			typ := rv.Type()

			if typ == nil {
				t.Errorf("rv.Type() returned nil for value %v", tc)
			} else {
				t.Logf("rv.Type() returned %v for value %v", typ, tc)
			}
		}
	})

	t.Run("Test specific struct case", func(t *testing.T) {
		x := struct {
			Name string
			Age  int
		}{"John", 25}

		// Test direct struct
		rv := tinyreflect.Indirect(tinyreflect.ValueOf(x))
		typ := rv.Type()

		if typ == nil {
			t.Error("Struct case: rv.Type() returned nil")
		} else {
			t.Logf("Struct case: rv.Type() returned %v", typ)
		}

		// Test pointer to struct
		rv2 := tinyreflect.Indirect(tinyreflect.ValueOf(&x))
		typ2 := rv2.Type()

		if typ2 == nil {
			t.Error("Pointer to struct case: rv.Type() returned nil")
		} else {
			t.Logf("Pointer to struct case: rv.Type() returned %v", typ2)
		}
	})
}