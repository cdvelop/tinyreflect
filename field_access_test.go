package tinyreflect_test

import (
	"testing"

	"github.com/cdvelop/tinyreflect"
)

// TestFieldAccess tests field access functionality through Value.Field()
// on both direct structs and pointers to structs.
func TestFieldAccess(t *testing.T) {
	tr := tinyreflect.New()
	type TestStruct struct {
		Name string
		ID   int
	}

	s := TestStruct{Name: "test", ID: 42}

	// Test on a direct struct value
	v := tr.ValueOf(s)
	t.Logf("Value kind: %v", v.Kind())
	field0, err := v.Field(0)
	if err != nil {
		t.Errorf("Error accessing field 0 on direct struct: %v", err)
		return
	}
	if name := field0.String(); name != "test" {
		t.Errorf("Expected field 0 value 'test', got '%s'", name)
	}
	t.Logf("Field 0 on direct struct accessed successfully.")

	// Test on a pointer to a struct
	p := &s
	pv := tr.ValueOf(p)
	t.Logf("Pointer value kind: %v", pv.Kind())

	// Use Indirect to get the struct from the pointer
	iv := tr.Indirect(pv)
	t.Logf("After Indirect - kind: %v", iv.Kind())

	field0Indirect, err := iv.Field(0)
	if err != nil {
		t.Errorf("Error accessing field 0 via indirect: %v", err)
		return
	}
	if name := field0Indirect.String(); name != "test" {
		t.Errorf("Expected field 0 value 'test' via indirect, got '%s'", name)
	}
	t.Logf("Field 0 via indirect accessed successfully.")
}