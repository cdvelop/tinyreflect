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

	// Get an addressable reflected value of the instance by taking a pointer and getting the element.
	v, err := tinyreflect.ValueOf(&instance).Elem()
	if err != nil {
		t.Fatalf("Failed to get addressable value: %v", err)
	}

	// Map to verify expected results
	expectedFields := map[string]any{
		"Name":     "Test Name",
		"Age":      30,
		"Active":   true,
		"Balance":  123.45,
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
		fieldName, err := typ.NameByIndex(i)
		if err != nil {
			t.Errorf("NameByIndex(%d) returned an unexpected error: %v", i, err)
			continue
		}

		fieldValue, err := v.Field(i)
		if err != nil {
			t.Errorf("Field(%d) returned an unexpected error: %v", i, err)
			continue
		}

		expectedValue, ok := expectedFields[fieldName]
		if !ok {
			t.Errorf("Unexpected field '%s' in the test struct", fieldName)
			continue
		}

		value, err := fieldValue.Interface()
		if err != nil {
			t.Errorf("field.Interface() returned an unexpected error for field '%s': %v", fieldName, err)
			continue
		}

		// Convert integer types for comparison to avoid type mismatches with any
		switch v := value.(type) {
		case int:
			if exp, ok := expectedValue.(int); !ok || v != exp {
				t.Errorf("Incorrect value for field '%s'. Expected: %v, Got: %v", fieldName, expectedValue, value)
			}
		default:
			if value != expectedValue {
				t.Errorf("Incorrect value for field '%s'. Expected: %v, Got: %v", fieldName, expectedValue, value)
			}
		}
	}
}
