package tinyreflect_test

import (
	"testing"

	"github.com/cdvelop/tinyreflect"
)

// UserWithName implements StructNamer interface
type UserWithName struct {
	ID   int
	Name string
}

func (UserWithName) StructName() string {
	return "User"
}

// UserWithoutName does NOT implement StructNamer interface
type UserWithoutName struct {
	ID   int
	Name string
}

// PersonWithName implements StructNamer interface with different name
type PersonWithName struct {
	Age int
}

func (PersonWithName) StructName() string {
	return "Person"
}

func TestTypeNameForStruct(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
		desc     string
	}{
		{
			name:     "struct_with_structnamer",
			value:    UserWithName{},
			expected: "User",
			desc:     "Struct implementing StructNamer should return custom name",
		},
		{
			name:     "struct_without_structnamer",
			value:    UserWithoutName{},
			expected: "struct",
			desc:     "Struct NOT implementing StructNamer should return 'struct'",
		},
		{
			name:     "person_with_structnamer",
			value:    PersonWithName{},
			expected: "Person",
			desc:     "Different struct implementing StructNamer should return its custom name",
		},
		{
			name:     "pointer_to_struct_with_structnamer",
			value:    &UserWithName{},
			expected: "User",
			desc:     "Pointer to struct implementing StructNamer should return custom name via Elem()",
		},
		{
			name:     "pointer_to_struct_without_structnamer",
			value:    &UserWithoutName{},
			expected: "struct",
			desc:     "Pointer to struct NOT implementing StructNamer should return 'struct' via Elem()",
		},
		{
			name:     "int_type",
			value:    42,
			expected: "int",
			desc:     "Non-struct types should return kind name",
		},
		{
			name:     "string_type",
			value:    "hello",
			expected: "string",
			desc:     "Non-struct types should return kind name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var actualName string

			typ := tinyreflect.TypeOf(tt.value)
			if typ == nil {
				t.Fatalf("TypeOf(%v) returned nil", tt.value)
			}

			// Handle pointer types
			if typ.Kind().String() == "ptr" {
				elem := typ.Elem()
				if elem == nil {
					t.Fatalf("Elem() returned nil for pointer type")
				}
				actualName = elem.Name()
			} else {
				actualName = typ.Name()
			}

			// Debug: check if actualName is somehow nil (shouldn't happen with string)
			if actualName == "" {
				t.Logf("DEBUG: actualName is empty string for test %s", tt.name)
			}

			if actualName != tt.expected {
				t.Errorf("%s: expected name '%s', got '%s'", tt.desc, tt.expected, actualName)
			} else {
				t.Logf("✓ %s: %s -> '%s'", tt.name, tt.desc, actualName)
			}
		})
	}
}
