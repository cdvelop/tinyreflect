package tinyreflect_test

import (
	"testing"

	"github.com/cdvelop/tinyreflect"
)

// TestAnonymousStruct tests behavior with anonymous structs
func TestAnonymousStruct(t *testing.T) {
	// Anonymous struct
	anon := struct {
		Field1 string
		Field2 int
	}{
		Field1: "test",
		Field2: 42,
	}

	typ := tinyreflect.TypeOf(anon)
	name := typ.Name()

	t.Logf("Anonymous struct name: '%s'", name)

	// Anonymous structs typically have empty names in Go reflection
	// But we need to handle this case appropriately
	if name == "" {
		t.Log("Anonymous struct correctly returns empty name")
	} else {
		t.Logf("Anonymous struct name: '%s' (implementation specific)", name)
	}

	// Verify it's still a struct type
	if typ.Kind().String() != "struct" {
		t.Errorf("Expected kind 'struct', got '%s'", typ.Kind().String())
	}
}

// TestNonStructTypes tests that Name() works correctly for non-struct types
func TestNonStructTypes(t *testing.T) {
	tests := []struct {
		name         string
		value        interface{}
		expectedKind string
		expectedName string
	}{
		{
			name:         "string type",
			value:        "hello",
			expectedKind: "string",
			expectedName: "string",
		},
		{
			name:         "int type",
			value:        42,
			expectedKind: "int",
			expectedName: "int",
		},
		{
			name:         "bool type",
			value:        true,
			expectedKind: "bool",
			expectedName: "bool",
		},
		{
			name:         "float64 type",
			value:        3.14,
			expectedKind: "float64",
			expectedName: "float64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typ := tinyreflect.TypeOf(tt.value)

			actualKind := typ.Kind().String()
			actualName := typ.Name()

			t.Logf("Value: %v, Kind: %s, Name: %s", tt.value, actualKind, actualName)

			if actualKind != tt.expectedKind {
				t.Errorf("Expected kind '%s', got '%s'", tt.expectedKind, actualKind)
			}

			if actualName != tt.expectedName {
				t.Errorf("Expected name '%s', got '%s'", tt.expectedName, actualName)
			}
		})
	}
}
