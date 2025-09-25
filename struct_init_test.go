package tinyreflect_test

import (
	"testing"

	"github.com/cdvelop/tinyreflect"
)

func TestStructFieldTypInitialization(t *testing.T) {
	tr := tinyreflect.New()
	type TestStruct struct {
		Name      string
		Timestamp int64
		Payload   []byte
		Ssid      []uint32
	}

	s := TestStruct{}
	typ := tr.TypeOf(s)

	if typ == nil {
		t.Fatal("TypeOf returned nil")
	}
	t.Logf("Type: %v, Kind: %v", typ, typ.Kind())

	if typ.Kind().String() != "struct" {
		t.Fatalf("Expected struct, got %v", typ.Kind())
	}

	numFields, err := typ.NumField()
	if err != nil {
		t.Fatalf("NumField() failed: %v", err)
	}
	t.Logf("NumField() returned: %d", numFields)

	if numFields != 4 {
		t.Fatalf("Expected 4 fields, got %d", numFields)
	}

	for i := 0; i < numFields; i++ {
		field, err := typ.Field(i)
		if err != nil {
			t.Fatalf("Field(%d) failed: %v", i, err)
		}
		if field.Typ == nil {
			t.Errorf("❌ Field %d (%s) has nil Typ!", i, field.Name)
		} else {
			t.Logf("✅ Field %d (%s) has Typ: %v, Kind: %v", i, field.Name, field.Typ, field.Typ.Kind())
		}
	}
}