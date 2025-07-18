package tinyreflect

import (
	"testing"
	"unsafe"
)

func TestTypeMethod(t *testing.T) {
	t.Run("Basic types", func(t *testing.T) {
		testCases := []struct {
			name  string
			value interface{}
		}{
			{"int", 42},
			{"string", "hello"},
			{"bool", true},
			{"float32", float32(3.14)},
			{"float64", 3.14},
			{"slice", []int{1, 2, 3}},
			{"struct", struct{ X int }{42}},
			{"pointer", &[]int{1, 2, 3}},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				v := ValueOf(tc.value)

				// Test that typ_ is not nil
				if v.typ_ == nil {
					t.Errorf("ValueOf(%v): typ_ is nil", tc.value)
				}

				// Test that Type() returns non-nil
				typ := v.Type()
				if typ == nil {
					t.Errorf("ValueOf(%v).Type() returned nil", tc.value)
				}

				// Test that we can get the kind
				if typ != nil {
					kind := typ.Kind()
					t.Logf("ValueOf(%v).Type().Kind() = %v", tc.value, kind)
				}
			})
		}
	})

	t.Run("Indirect test", func(t *testing.T) {
		x := 42
		ptr := &x

		// Test direct pointer
		v1 := ValueOf(ptr)
		if v1.typ_ == nil {
			t.Error("ValueOf(ptr): typ_ is nil")
		}
		typ1 := v1.Type()
		if typ1 == nil {
			t.Error("ValueOf(ptr).Type() returned nil")
		}

		// Test indirect
		v2 := Indirect(v1)
		if v2.typ_ == nil {
			t.Error("Indirect(ValueOf(ptr)): typ_ is nil")
		}
		typ2 := v2.Type()
		if typ2 == nil {
			t.Error("Indirect(ValueOf(ptr)).Type() returned nil")
		}
	})

	t.Run("Nil values", func(t *testing.T) {
		// Test nil interface
		v1 := ValueOf(nil)
		if v1.typ_ != nil {
			t.Error("ValueOf(nil): typ_ should be nil")
		}
		typ1 := v1.Type()
		if typ1 != nil {
			t.Error("ValueOf(nil).Type() should return nil")
		}

		// Test nil pointer
		var ptr *int
		v2 := ValueOf(ptr)
		if v2.typ_ == nil {
			t.Error("ValueOf(nil pointer): typ_ is nil")
		}
		typ2 := v2.Type()
		if typ2 == nil {
			t.Error("ValueOf(nil pointer).Type() returned nil")
		}
	})
}

func TestEmptyInterface(t *testing.T) {
	t.Run("EmptyInterface layout", func(t *testing.T) {
		x := 42
		i := interface{}(x)

		// Test that we can cast to EmptyInterface
		e := (*EmptyInterface)(unsafe.Pointer(&i))
		if e.Type == nil {
			t.Error("EmptyInterface.Type is nil")
		}
		if e.Data == nil {
			t.Error("EmptyInterface.Data is nil")
		}

		t.Logf("EmptyInterface: Type=%p, Data=%p", e.Type, e.Data)
	})
}

func TestUnpackEface(t *testing.T) {
	t.Run("unpackEface basic", func(t *testing.T) {
		testCases := []interface{}{
			42,
			"hello",
			true,
			[]int{1, 2, 3},
		}

		for _, tc := range testCases {
			v := unpackEface(tc)

			if v.typ_ == nil {
				t.Errorf("unpackEface(%v): typ_ is nil", tc)
			}

			typ := v.Type()
			if typ == nil {
				t.Errorf("unpackEface(%v).Type() returned nil", tc)
			}

			t.Logf("unpackEface(%v): typ_=%p, Type()=%p", tc, v.typ_, typ)
		}
	})
}

func TestTypeStability(t *testing.T) {
	t.Run("Type pointer stability", func(t *testing.T) {
		x := 42
		v := ValueOf(x)

		// Get type multiple times
		typ1 := v.Type()
		typ2 := v.Type()
		typ3 := v.typ()

		if typ1 == nil || typ2 == nil || typ3 == nil {
			t.Error("Type methods returned nil")
		}

		if typ1 != typ2 {
			t.Error("Type() should return consistent pointer")
		}

		if typ1 != typ3 {
			t.Error("Type() and typ() should return same pointer")
		}
	})
}
