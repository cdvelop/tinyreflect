package tinyreflect

import (
	"testing"
)

func TestPtrTypeMethod(t *testing.T) {
	// Test with a pointer type
	x := 42
	ptr := &x

	v := ValueOf(ptr)
	typ := v.Type()

	if typ == nil {
		t.Fatal("ValueOf(ptr).Type() returned nil")
	}

	t.Logf("Type: %p, Kind: %v", typ, typ.Kind())

	// Test PtrType method
	ptrType := typ.PtrType()
	if ptrType == nil {
		t.Errorf("PtrType() returned nil for pointer type")
	} else {
		t.Logf("PtrType() returned %p", ptrType)
		if ptrType.Elem == nil {
			t.Error("PtrType.Elem is nil")
		} else {
			t.Logf("PtrType.Elem: %p, Kind: %v", ptrType.Elem, ptrType.Elem.Kind())
		}
	}

	// Test with struct pointer (like in the failing test)
	type simpleStruct struct {
		Name string
	}

	s := &simpleStruct{}
	v2 := ValueOf(s)
	typ2 := v2.Type()

	if typ2 == nil {
		t.Fatal("ValueOf(struct_ptr).Type() returned nil")
	}

	t.Logf("Struct ptr Type: %p, Kind: %v", typ2, typ2.Kind())

	ptrType2 := typ2.PtrType()
	if ptrType2 == nil {
		t.Errorf("PtrType() returned nil for struct pointer type")
	} else {
		t.Logf("Struct ptr PtrType() returned %p", ptrType2)
		if ptrType2.Elem == nil {
			t.Error("Struct ptr PtrType.Elem is nil")
		} else {
			t.Logf("Struct ptr PtrType.Elem: %p, Kind: %v", ptrType2.Elem, ptrType2.Elem.Kind())
		}
	}
}
