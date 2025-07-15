package tinyreflect

import (
	"testing"
)

// TestStruct is a helper struct for testing reflection operations.
type TestStruct struct {
	Name string
	Age  int
}

func TestSetString(t *testing.T) {
	// The struct we want to modify
	data := TestStruct{Name: "initial", Age: 30}
	const newName = "changed"

	// Get a Value of the struct. Since we want to modify the original 'data'
	// variable, we need to work with a pointer to it.
	v := ValueOf(&data)

	// The ValueOf(&data) returns a Ptr. We need to get the element it points to
	// before we can access its fields. This is what Elem() does.
	structVal, err := v.Elem()
	if err != nil {
		t.Fatalf("Failed to get element from pointer value: %v", err)
	}

	// Find the field we want to modify.
	// For this test, we know "Name" is the first field (index 0).
	nameField, err := structVal.Field(0)
	if err != nil {
		t.Fatalf("Failed to get field 'Name': %v", err)
	}

	// Now, set the new string value on the 'Name' field.
	nameField.SetString(newName)

	// Finally, verify that the original struct's field has been updated.
	if data.Name != newName {
		t.Errorf("SetString failed: expected Name to be %q, but got %q", newName, data.Name)
	}
}
