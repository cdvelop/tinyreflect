package tinyreflect_test

import (
	"testing"

	"github.com/cdvelop/tinyreflect"
)

// Test_StructIteration iterates over the fields of a struct to verify names and values.
func Test_StructIteration(t *testing.T) {
	// Test struct with public and private fields
	type sampleStruct struct {
		Name     string
		Age      int
		Active   bool
		Balance  float64
		private1 int
		_hidden  string
	}

	// Instance of the struct with values
	instance := sampleStruct{
		Name:     "Test Name",
		Age:      30,
		Active:   true,
		Balance:  123.45,
		private1: 42,
		_hidden:  "secret",
	}

	// Reflected value of the instance
	v := tinyreflect.ValueOf(instance)

	// Map to verify expected results
	expectedFields := map[string]interface{}{
		"Name":     "Test Name",
		"Age":      int(30),
		"Active":   true,
		"Balance":  float64(123.45),
		"private1": 42,
		"_hidden":  "secret",
	}

	typ := v.Type()
	numFields, err := typ.NumField()
	if err != nil {
		t.Fatalf("NumField() returned an unexpected error: %v", err)
	}

	// Iterate over the struct fields
	for i := 0; i < numFields; i++ {
		// Get the field name by index
		fieldName, err := typ.NameByIndex(i)
		if err != nil {
			t.Errorf("NameByIndex(%d) returned an unexpected error: %v", i, err)
			continue
		}

		// Get the field value by index
		fieldValue, err := v.Field(i)
		if err != nil {
			t.Errorf("Field(%d) returned an unexpected error: %v", i, err)
			continue
		}

		// Verify that the field exists in the expected values map
		expectedValue, ok := expectedFields[fieldName]
		if !ok {
			t.Errorf("Unexpected field '%s' in the test struct", fieldName)
			continue
		}

		// Get the field value and compare it
		value, err := fieldValue.Interface()
		if err != nil {
			t.Errorf("field.Interface() returned an unexpected error for field '%s': %v", fieldName, err)
			continue
		}

		if value != expectedValue {
			t.Errorf("Incorrect value for field '%s'. Expected: %v, Got: %v", fieldName, expectedValue, value)
		}

		// Remove the field from the map to ensure all fields are evaluated
		delete(expectedFields, fieldName)
	}

	// If the map is not empty, it means not all expected fields were found
	if len(expectedFields) > 0 {
		for name := range expectedFields {
			t.Errorf("Expected field '%s' was not found during iteration", name)
		}
	}
}
