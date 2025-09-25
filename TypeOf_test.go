package tinyreflect_test

import (
	"fmt"
	"strings"
	"sync"
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

// A thread-safe buffer to capture log messages.
type concurrentBuffer struct {
	b  strings.Builder
	mu sync.Mutex
}

func (b *concurrentBuffer) Write(p []byte) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.Write(p)
}

func (b *concurrentBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.String()
}

func TestStructNameCaching(t *testing.T) {
	tests := []struct {
		name        string
		value       interface{}
		expectedLog string
		desc        string
	}{
		{
			name:        "struct_with_structnamer",
			value:       UserWithName{},
			expectedLog: "cached schema for struct User",
			desc:        "Struct implementing StructNamer should be cached with its custom name",
		},
		{
			name:        "struct_without_structnamer",
			value:       UserWithoutName{},
			expectedLog: "cached schema for struct struct",
			desc:        "Struct without StructNamer should be cached with 'struct' as name",
		},
		{
			name:        "person_with_structnamer",
			value:       PersonWithName{},
			expectedLog: "cached schema for struct Person",
			desc:        "Different struct implementing StructNamer should return its custom name",
		},
		{
			name:        "pointer_to_struct_with_structnamer",
			value:       &UserWithName{},
			expectedLog: "cached schema for struct User",
			desc:        "Pointer to struct with StructNamer should also be cached correctly",
		},
		{
			name:        "pointer_to_struct_without_structnamer",
			value:       &UserWithoutName{},
			expectedLog: "cached schema for struct struct",
			desc:        "Pointer to struct without StructNamer should also be cached correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf concurrentBuffer
			logger := func(msgs ...any) {
				fmt.Fprintln(&buf, msgs...)
			}

			tr := tinyreflect.New(logger)
			tr.TypeOf(tt.value) // This should trigger the caching

			logOutput := buf.String()
			if !strings.Contains(logOutput, tt.expectedLog) {
				t.Errorf("%s: expected log to contain '%s', but got '%s'", tt.desc, tt.expectedLog, logOutput)
			} else {
				t.Logf("✓ %s -> Log contains '%s'", tt.name, tt.expectedLog)
			}
		})
	}
}

func TestTypeNameForNonStructs(t *testing.T) {
	tr := tinyreflect.New()
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{name: "int_type", value: 42, expected: "int"},
		{name: "string_type", value: "hello", expected: "string"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typ := tr.TypeOf(tt.value)
			if typ.Name() != tt.expected {
				t.Errorf("Expected name '%s', got '%s'", tt.expected, typ.Name())
			}
		})
	}
}

func TestTypeMethods(t *testing.T) {
	tr := tinyreflect.New()

	// Test SliceType
	slice := []int{}
	typ := tr.TypeOf(slice)
	if typ.SliceType() == nil {
		t.Error("SliceType: expected non-nil")
	}

	// Test ArrayType
	arr := [3]int{}
	typ = tr.TypeOf(arr)
	if typ.ArrayType() == nil {
		t.Error("ArrayType: expected non-nil")
	}

	// Test PtrType
	var p *int
	typ = tr.TypeOf(p)
	if typ.PtrType() == nil {
		t.Error("PtrType: expected non-nil")
	}

	// Test StructID
	s := UserWithName{}
	typ = tr.TypeOf(s)
	if typ.StructID() == 0 {
		t.Error("StructID: expected non-zero")
	}
	i := 123
	typ = tr.TypeOf(i)
	if typ.StructID() != 0 {
		t.Error("StructID on non-struct: expected 0")
	}

	// Test StructType
	typ = tr.TypeOf(s)
	if typ.StructType() == nil {
		t.Error("StructType: expected non-nil")
	}
	typ = tr.TypeOf(i)
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
	typ = tr.TypeOf(s)
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