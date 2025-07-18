package tinyreflect

import (
	"testing"

	. "github.com/cdvelop/tinystring"
)

func TestDecodeExactScenario(t *testing.T) {
	// This replicates the exact scenario from decoder.go and scanner.go
	type simpleStruct struct {
		Name      string
		Timestamp int64
		Payload   []byte
		Ssid      []uint32
	}

	// Test scenario like in decoder: pointer to struct
	s := &simpleStruct{}

	// decoder.go line 46: rv := tinyreflect.Indirect(tinyreflect.ValueOf(v))
	rv := Indirect(ValueOf(s))

	// Check that rv has a valid type
	typ := rv.Type()
	if typ == nil {
		t.Fatal("rv.Type() returned nil - this is the 'value type nil' error")
	}

	t.Logf("rv.Type() = %p, Kind: %v", typ, typ.Kind())

	// decoder.go line 52: scanToCache(rv.Type(), d.schemas)
	// scanner.go line 20: scanTypeWithTinyReflect(t)
	// This should work for struct types

	// But let's test what happens when we have a pointer type
	// sometimes the Indirect might not work correctly
	originalV := ValueOf(s)
	originalTyp := originalV.Type()

	t.Logf("Original ValueOf(s).Type() = %p, Kind: %v", originalTyp, originalTyp.Kind())

	// scanner.go line 42: ptrType := t.PtrType()
	if originalTyp.Kind() == K.Pointer {
		ptrType := originalTyp.PtrType()
		if ptrType == nil {
			t.Fatal("PtrType() returned nil for pointer type")
		}
		t.Logf("PtrType() = %p, Elem = %p", ptrType, ptrType.Elem)

		if ptrType.Elem == nil {
			t.Fatal("PtrType.Elem is nil - this could be the issue")
		}
		t.Logf("PtrType.Elem Kind: %v", ptrType.Elem.Kind())
	}

	// Test the exact conditions from getElementType
	testGetElementType(t, originalTyp)
	testGetElementType(t, typ)
}

func testGetElementType(t *testing.T, typ *Type) {
	t.Logf("Testing getElementType for Type: %p, Kind: %v", typ, typ.Kind())

	kind := typ.Kind()
	switch kind {
	case K.Pointer:
		ptrType := typ.PtrType()
		if ptrType == nil {
			t.Errorf("getElementType: PtrType() returned nil for pointer type")
			return
		}
		if ptrType.Elem == nil {
			t.Errorf("getElementType: PtrType.Elem is nil")
			return
		}
		t.Logf("getElementType: Success - Elem Kind: %v", ptrType.Elem.Kind())
	case K.Array:
		arrayType := typ.ArrayType()
		if arrayType == nil {
			t.Errorf("getElementType: ArrayType() returned nil for array type")
			return
		}
		t.Logf("getElementType: Array element type: %v", arrayType.Elem.Kind())
	case K.Slice:
		sliceType := typ.SliceType()
		if sliceType == nil {
			t.Errorf("getElementType: SliceType() returned nil for slice type")
			return
		}
		t.Logf("getElementType: Slice element type: %v", sliceType.Elem.Kind())
	default:
		t.Logf("getElementType: Not a container type - Kind: %v", kind)
	}
}
