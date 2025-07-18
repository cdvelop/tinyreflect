package tinyreflect

import (
	"testing"
)

func TestStructFieldTypDebug(t *testing.T) {
	// This is the exact type from the failing test
	type simpleStruct struct {
		Name      string
		Timestamp int64
		Payload   []byte
		Ssid      []uint32
	}

	s := &simpleStruct{}
	rv := Indirect(ValueOf(s))
	typ := rv.Type()

	if typ == nil {
		t.Fatal("typ is nil")
	}

	t.Logf("Type: %p, Kind: %v", typ, typ.Kind())

	// Get number of fields
	numFields, err := typ.NumField()
	if err != nil {
		t.Fatalf("NumField failed: %v", err)
	}

	t.Logf("NumFields: %d", numFields)

	// Check each field
	for i := 0; i < numFields; i++ {
		field, err := typ.Field(i)
		if err != nil {
			t.Fatalf("Field(%d) failed: %v", i, err)
		}

		fieldName := field.Name.String()
		t.Logf("Field %d: Name=%s, Typ=%p", i, fieldName, field.Typ)

		if field.Typ == nil {
			t.Fatalf("Field %d (%s) has nil Typ - this is the 'value type nil' error!", i, fieldName)
		}

		t.Logf("Field %d (%s) has Typ Kind: %v", i, fieldName, field.Typ.Kind())

		// This is what the scanner does - call scanTypeWithTinyReflect on field.Typ
		// If field.Typ is nil, this will cause the "value type nil" error
		t.Logf("Field %d (%s) would call scanTypeWithTinyReflect successfully", i, fieldName)
	}

	t.Logf("All fields have valid Typ - no 'value type nil' error should occur")
}
