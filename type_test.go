package tinyreflect_test

import (
	"testing"

	"github.com/cdvelop/tinyreflect"
)

func TestTypeMethod(t *testing.T) {

	t.Run("Basic types", func(t *testing.T) {
		testCases := []struct {
			name  string
			value any
		}{
			{"int", 42},
			{"string", "hello"},
			{"bool", true},
			{"slice", []int{1, 2, 3}},
			{"struct", struct{ X int }{42}},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				v := tinyreflect.ValueOf(tc.value)
				typ := v.Type()
				if typ == nil {
					t.Errorf("ValueOf(%v).Type() returned nil", tc.value)
				}
			})
		}
	})

	t.Run("Indirect on Pointers", func(t *testing.T) {
		x := 42
		ptr := &x

		// Test direct pointer
		v1 := tinyreflect.ValueOf(ptr)
		typ1 := v1.Type()
		if typ1 == nil {
			t.Fatal("ValueOf(ptr).Type() returned nil")
		}
		if typ1.Kind().String() != "ptr" {
			t.Errorf("Expected pointer kind, got %v", typ1.Kind())
		}

		// Test after indirect
		v2 := tinyreflect.Indirect(v1)
		typ2 := v2.Type()
		if typ2 == nil {
			t.Fatal("Indirect(ValueOf(ptr)).Type() returned nil")
		}
		if typ2.Kind().String() != "int" {
			t.Errorf("Expected int kind after indirect, got %v", typ2.Kind())
		}
	})

	t.Run("Nil values", func(t *testing.T) {
		// Test nil interface
		v1 := tinyreflect.ValueOf(nil)
		typ1 := v1.Type()
		if typ1 != nil {
			t.Errorf("ValueOf(nil).Type() should return nil, but got %v", typ1)
		}

		// Test nil pointer
		var ptr *int
		v2 := tinyreflect.ValueOf(ptr)
		typ2 := v2.Type()
		if typ2 == nil {
			t.Error("ValueOf(nil pointer).Type() returned nil")
		}
	})
}

func TestTypeStability(t *testing.T) {
	x := 42
	v := tinyreflect.ValueOf(x)

	// Get type multiple times
	typ1 := v.Type()
	typ2 := v.Type()

	if typ1 == nil || typ2 == nil {
		t.Fatal("Type() returned nil")
	}

	if typ1 != typ2 {
		t.Error("Type() should return a consistent pointer on subsequent calls")
	}
}
