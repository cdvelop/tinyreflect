package tinyreflect_test

import (
	"testing"

	"github.com/cdvelop/tinyreflect"
)

func TestStructFieldTypes(t *testing.T) {
	tr := tinyreflect.New()
	type TestStruct struct {
		Name   string
		Age    int
		Active bool
		Scores []float64
	}

	typ := tr.TypeOf(TestStruct{})
	if typ == nil {
		t.Fatal("TypeOf returned nil")
	}

	if typ.Kind().String() != "struct" {
		t.Fatalf("Expected struct kind, got %v", typ.Kind())
	}

	numFields, err := typ.NumField()
	if err != nil {
		t.Fatalf("NumField failed: %v", err)
	}
	t.Logf("Struct has %d fields", numFields)

	for i := 0; i < numFields; i++ {
		field, err := typ.Field(i)
		if err != nil {
			t.Fatalf("Field(%d) failed: %v", i, err)
		}

		fieldName := field.Name.String()
		t.Logf("Field %d: Name=%s, Typ=%p", i, fieldName, field.Typ)

		if field.Typ == nil {
			t.Errorf("Field %d (%s) has nil Typ", i, fieldName)
		} else {
			t.Logf("Field %d (%s) has Typ Kind: %v", i, fieldName, field.Typ.Kind())
		}
	}
}

func TestEmbedded(t *testing.T) {
	tr := tinyreflect.New()
	type E struct {
		e int
	}
	type S struct {
		E
		s int
	}

	typ := tr.TypeOf(S{})
	if typ == nil {
		t.Fatal("TypeOf returned nil")
	}

	// Test embedded field
	field0, err := typ.Field(0)
	if err != nil {
		t.Fatalf("Field(0) failed: %v", err)
	}
	if !field0.Embedded() {
		t.Error("Embedded for embedded field 'E': expected true, got false")
	}

	// Test non-embedded field
	field1, err := typ.Field(1)
	if err != nil {
		t.Fatalf("Field(1) failed: %v", err)
	}
	if field1.Embedded() {
		t.Error("Embedded for non-embedded field 's': expected false, got true")
	}
}