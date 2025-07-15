package tinyreflect_test

import (
	"testing"

	"github.com/cdvelop/tinyreflect"
)

// Test struct types with different names for name extraction testing
type Customer struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Item struct {
	Title string  `json:"title"`
	Price float64 `json:"price"`
}

type Organization struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

// TestStructName validates that Type.Name() returns the actual struct name
func TestStructName(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{
			name:     "Customer struct name",
			value:    Customer{},
			expected: "Customer",
		},
		{
			name:     "Item struct name",
			value:    Item{},
			expected: "Item",
		},
		{
			name:     "Organization struct name",
			value:    Organization{},
			expected: "Organization",
		},
		{
			name:     "Customer pointer name",
			value:    &Customer{},
			expected: "Customer", // Should return the dereferenced type name
		},
		{
			name:     "Item pointer name",
			value:    &Item{},
			expected: "Item",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typ := tinyreflect.TypeOf(tt.value)
			
			// Handle pointer types - get the element type
			if typ.Kind().String() == "ptr" {
				// For pointer types, we need to get the element type
				// This will be implemented in Type.Elem() method
				t.Logf("Pointer type detected, getting element type...")
				// For now, skip pointer tests until Elem() is implemented
				t.Skip("Pointer support not yet implemented - need Type.Elem() method")
				return
			}

			// Test the Name() method
			actualName := typ.Name()
			t.Logf("Type: %T, Kind: %s, Name: %s", tt.value, typ.Kind().String(), actualName)

			if actualName != tt.expected {
				t.Errorf("Expected struct name '%s', got '%s'", tt.expected, actualName)
			}

			// Additional validation: ensure it's a struct type
			if typ.Kind().String() != "struct" {
				t.Errorf("Expected kind 'struct', got '%s'", typ.Kind().String())
			}
		})
	}
}

// TestStructNameConsistency validates that same struct types have same names
func TestStructNameConsistency(t *testing.T) {
	// Different instances of same struct type
	customer1 := Customer{}
	customer2 := Customer{Name: "John", Age: 30}
	customer3 := &Customer{Name: "Jane", Age: 25}

	type1 := tinyreflect.TypeOf(customer1)
	type2 := tinyreflect.TypeOf(customer2)
	type3 := tinyreflect.TypeOf(*customer3) // Dereference pointer

	name1 := type1.Name()
	name2 := type2.Name()
	name3 := type3.Name()

	t.Logf("Name1: %s", name1)
	t.Logf("Name2: %s", name2)
	t.Logf("Name3: %s", name3)

	// All should have the same name
	if name1 != name2 {
		t.Errorf("Inconsistent names: %s != %s", name1, name2)
	}
	if name1 != name3 {
		t.Errorf("Inconsistent names: %s != %s", name1, name3)
	}
	if name2 != name3 {
		t.Errorf("Inconsistent names: %s != %s", name2, name3)
	}

	// Name should be "Customer"
	expectedName := "Customer"
	if name1 != expectedName {
		t.Errorf("Expected name '%s', got '%s'", expectedName, name1)
	}
}

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
		name          string
		value         interface{}
		expectedKind  string
		expectedName  string
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
