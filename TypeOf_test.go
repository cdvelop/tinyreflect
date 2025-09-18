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

func TestTypeMethods(t *testing.T) {
	// Test SliceType
	slice := []int{}
	typ := tinyreflect.TypeOf(slice)
	if typ.SliceType() == nil {
		t.Error("SliceType: expected non-nil")
	}

	// Test ArrayType
	arr := [3]int{}
	typ = tinyreflect.TypeOf(arr)
	if typ.ArrayType() == nil {
		t.Error("ArrayType: expected non-nil")
	}

	// Test PtrType
	var p *int
	typ = tinyreflect.TypeOf(p)
	if typ.PtrType() == nil {
		t.Error("PtrType: expected non-nil")
	}

	// Test StructID
	s := UserWithName{}
	typ = tinyreflect.TypeOf(s)
	if typ.StructID() == 0 {
		t.Error("StructID: expected non-zero")
	}
	i := 123
	typ = tinyreflect.TypeOf(i)
	if typ.StructID() != 0 {
		t.Error("StructID on non-struct: expected 0")
	}

	// Test StructType
	typ = tinyreflect.TypeOf(s)
	if typ.StructType() == nil {
		t.Error("StructType: expected non-nil")
	}
	typ = tinyreflect.TypeOf(i)
	if typ.StructType() != nil {
		t.Error("StructType on non-struct: expected nil")
	}

	// Test Field error
	_, err := typ.Field(0)
	if err == nil {
		t.Error("Field on non-struct: expected an error")
	}

	// Test NumField error
	_, err = typ.NumField()
	if err == nil {
		t.Error("NumField on non-struct: expected an error")
	}

	// Test NameByIndex error
	_, err = typ.NameByIndex(0)
	if err == nil {
		t.Error("NameByIndex on non-struct: expected an error")
	}

	// Test SliceType on non-slice
	if typ.SliceType() != nil {
		t.Error("SliceType on non-slice: expected nil")
	}

	// Test ArrayType on non-array
	if typ.ArrayType() != nil {
		t.Error("ArrayType on non-array: expected nil")
	}

	// Test PtrType on non-pointer
	if typ.PtrType() != nil {
		t.Error("PtrType on non-pointer: expected nil")
	}

	// Test Field with out of range index
	typ = tinyreflect.TypeOf(s)
	_, err = typ.Field(2)
	if err == nil {
		t.Error("Field with out of range index: expected an error")
	}
	_, err = typ.Field(-1)
	if err == nil {
		t.Error("Field with negative index: expected an error")
	}

	// Test NameByIndex with out of range index
	_, err = typ.NameByIndex(2)
	if err == nil {
		t.Error("NameByIndex with out of range index: expected an error")
	}
	_, err = typ.NameByIndex(-1)
	if err == nil {
		t.Error("NameByIndex with negative index: expected an error")
	}

	// Test Elem on non-elem type
	if typ.Elem() != nil {
		t.Error("Elem on non-elem type: expected nil")
	}
}
