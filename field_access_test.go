package tinyreflect

import (
	"testing"
)

// TestFieldAccess tests field access functionality through Value.Field()
// This provides coverage for field access on both direct structs and pointers to structs
func TestFieldAccess(t *testing.T) {
	type TestStruct struct {
		Name string
		ID   int
	}

	// Create a test struct
	s := TestStruct{Name: "test", ID: 42}

	// Create a Value from it
	v := ValueOf(s)

	t.Logf("Value type: %v", v.Type())
	t.Logf("Value kind: %v", v.Kind())

	// Try to access the first field
	field0, err := v.Field(0)
	if err != nil {
		t.Errorf("Error accessing field 0: %v", err)
		return
	}

	t.Logf("Field 0 accessed successfully: %v", field0.Type())

	// Now test with a pointer to the struct
	p := &s
	pv := ValueOf(p)

	t.Logf("Pointer value type: %v", pv.Type())
	t.Logf("Pointer value kind: %v", pv.Kind())

	// Use Indirect to get the struct
	iv := Indirect(pv)

	t.Logf("After Indirect - type: %v", iv.Type())
	t.Logf("After Indirect - kind: %v", iv.Kind())

	// Try to access the first field
	field0_indirect, err := iv.Field(0)
	if err != nil {
		t.Errorf("Error accessing field 0 via indirect: %v", err)
		return
	}

	t.Logf("Field 0 via indirect accessed successfully: %v", field0_indirect.Type())
}
