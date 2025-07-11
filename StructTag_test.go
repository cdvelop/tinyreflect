package tinyreflect_test

import (
	"fmt"
	"testing"

	"github.com/cdvelop/tinyreflect"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestGetFieldName(t *testing.T) {
	p := Person{"Cesar", 30}
	to := tinyreflect.TypeOf(p)

	numField, err := to.NumField()
	if err != nil {
		t.Fatalf("NumField() failed: %v", err)
	}

	if numField != 2 {
		t.Fatalf("expected 2 fields, but got %d", numField)
	}

	expectedFields := []struct {
		Name string
		Type string
		Tag  string
	}{
		{Name: "Name", Type: "string", Tag: `json:"name"`},
		{Name: "Age", Type: "int", Tag: `json:"age"`},
	}

	// Iterate over the struct fields
	for i := 0; i < numField; i++ {
		field, err := to.Field(i)
		if err != nil {
			t.Fatalf("Field(%d) failed: %v", i, err)
		}

		// Test field name
		if field.Name.String() != expectedFields[i].Name {
			t.Errorf("Field %d: expected name %s, got %s", i, expectedFields[i].Name, field.Name.String())
		}

		// Test field type
		if field.Typ.String() != expectedFields[i].Type {
			t.Errorf("Field %d: expected type %s, got %s", i, expectedFields[i].Type, field.Typ.String())
		}

		// Test field tag
		if string(field.Tag()) != expectedFields[i].Tag {
			t.Errorf("Field %d: expected tag %s, got %s", i, expectedFields[i].Tag, string(field.Tag()))
		}

		fmt.Printf("Field %d: %s (Type: %s, Tag: '%s')\n", i, field.Name.String(), field.Typ.String(), string(field.Tag()))
	}
}
