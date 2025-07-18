package tinyreflect

import (
	"testing"
	"unsafe"
)

func TestStructFieldTypInitialization(t *testing.T) {
	type TestStruct struct {
		Name      string
		Timestamp int64
		Payload   []byte
		Ssid      []uint32
	}

	s := TestStruct{}
	typ := TypeOf(s)

	if typ == nil {
		t.Fatal("TypeOf returned nil")
	}

	t.Logf("Type: %p, Kind: %v", typ, typ.Kind())

	// Check if it's a struct
	if typ.Kind().String() != "struct" {
		t.Fatalf("Expected struct, got %v", typ.Kind())
	}

	// Cast to StructType to access Fields directly
	st := (*StructType)(unsafe.Pointer(typ))
	t.Logf("StructType.Fields length: %d", len(st.Fields))

	// Check each field
	for i, field := range st.Fields {
		t.Logf("Field %d: Name=%s, Typ=%p", i, field.Name, field.Typ)
		if field.Typ == nil {
			t.Errorf("❌ Field %d (%s) has nil Typ!", i, field.Name)
		} else {
			t.Logf("✅ Field %d (%s) has Typ: %p, Kind: %v", i, field.Name, field.Typ, field.Typ.Kind())
		}
	}

	// Also test Field() method
	numFields, err := typ.NumField()
	if err != nil {
		t.Fatalf("NumField() failed: %v", err)
	}
	t.Logf("NumField() returned: %d", numFields)

	for i := 0; i < numFields; i++ {
		field, err := typ.Field(i)
		if err != nil {
			t.Fatalf("Field(%d) failed: %v", i, err)
		}
		t.Logf("Field(%d) via method: Name=%s, Typ=%p", i, field.Name, field.Typ)
	}
}
