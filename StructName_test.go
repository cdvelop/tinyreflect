package tinyreflect_test

import (
	"testing"

	"github.com/cdvelop/tinyreflect"
)

func TestTypeName(t *testing.T) {
	// Anonymous struct for testing
	anon := struct {
		Field1 string
		Field2 int
	}{"test", 42}

	tests := []struct {
		name         string
		value        interface{}
		expectedName string
	}{
		{"Anonymous Struct", anon, "struct"},
		{"string", "hello", "string"},
		{"int", 42, "int"},
		{"bool", true, "bool"},
		{"float64", 3.14, "float64"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typ := tinyreflect.TypeOf(tt.value)
			if typ == nil {
				t.Fatalf("TypeOf returned nil for value %v", tt.value)
			}

			actualName := typ.Name()
			if actualName != tt.expectedName {
				t.Errorf("Expected name '%s', got '%s'", tt.expectedName, actualName)
			}
		})
	}
}