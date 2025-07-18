package tinyreflect

import (
	"testing"

	. "github.com/cdvelop/tinystring"
)

func TestStructFieldTypes(t *testing.T) {
	type TestStruct struct {
		Name   string
		Age    int
		Active bool
		Scores []float64
	}

	typ := TypeOf(TestStruct{})
	if typ == nil {
		t.Fatal("TypeOf returned nil")
	}

	if typ.Kind() != K.Struct {
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
