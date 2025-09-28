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

// PersonWithName implements StructNamer interface with a different name
type PersonWithName struct {
	Age int
}

func (PersonWithName) StructName() string {
	return "Person"
}

func TestTypeOfBasicFunctionality(t *testing.T) {
	tests := []struct {
		name  string
		value any
		desc  string
	}{
		{
			name:  "struct_with_structnamer",
			value: UserWithName{},
			desc:  "Struct implementing StructNamer should return correct Type",
		},
		{
			name:  "struct_without_structnamer",
			value: UserWithoutName{},
			desc:  "Struct without StructNamer should return correct Type",
		},
		{
			name:  "person_with_structnamer",
			value: PersonWithName{},
			desc:  "Different struct implementing StructNamer should return correct Type",
		},
		{
			name:  "pointer_to_struct_with_structnamer",
			value: &UserWithName{},
			desc:  "Pointer to struct with StructNamer should return pointer Type",
		},
		{
			name:  "pointer_to_struct_without_structnamer",
			value: &UserWithoutName{},
			desc:  "Pointer to struct without StructNamer should return pointer Type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typ := tinyreflect.TypeOf(tt.value)

			if typ == nil {
				t.Errorf("%s: TypeOf returned nil", tt.desc)
				return
			}

			// Verify that TypeOf works correctly without caching
			if tt.name == "pointer_to_struct_with_structnamer" || tt.name == "pointer_to_struct_without_structnamer" {
				if typ.Kind().String() != "ptr" {
					t.Errorf("%s: expected Pointer kind, got %v", tt.desc, typ.Kind())
				}
			} else {
				if typ.Kind().String() != "struct" {
					t.Errorf("%s: expected Struct kind, got %v", tt.desc, typ.Kind())
				}
			}

			t.Logf("✓ %s -> TypeOf works correctly", tt.name)
		})
	}
}

func TestTypeNameForNonStructs(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		expected string
	}{
		{name: "int_type", value: 42, expected: "int"},
		{name: "string_type", value: "hello", expected: "string"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typ := tinyreflect.TypeOf(tt.value)
			if typ.Name() != tt.expected {
				t.Errorf("Expected name '%s', got '%s'", tt.expected, typ.Name())
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
