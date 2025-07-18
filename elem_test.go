package tinyreflect

import (
	"testing"

	. "github.com/cdvelop/tinystring"
)

func TestElem(t *testing.T) {
	// Test basic pointer dereference
	x := 42
	ptr := &x

	v := ValueOf(ptr)
	elem, err := v.Elem()
	if err != nil {
		t.Fatalf("Elem() failed: %v", err)
	}

	// Check that the element has the correct type
	if elem.Type() == nil {
		t.Fatal("Elem() returned nil type")
	}

	if elem.Type().Kind() != K.Int {
		t.Errorf("Expected Kind %v, got %v", K.Int, elem.Type().Kind())
	}

	// Test nil pointer
	var nilPtr *int
	vNil := ValueOf(nilPtr)
	elemNil, err := vNil.Elem()
	if err != nil {
		t.Fatalf("Elem() on nil pointer failed: %v", err)
	}

	// Should return zero value for nil pointer
	if elemNil.Type() != nil {
		t.Error("Expected nil type for nil pointer elem")
	}

	// Test error case - non-pointer
	vInt := ValueOf(42)
	_, err = vInt.Elem()
	if err == nil {
		t.Error("Expected error when calling Elem() on non-pointer")
	}
}

func TestElemStruct(t *testing.T) {
	// Test with struct pointer
	type TestStruct struct {
		X int
		Y string
	}

	s := &TestStruct{X: 10, Y: "hello"}
	v := ValueOf(s)
	elem, err := v.Elem()
	if err != nil {
		t.Fatalf("Elem() failed: %v", err)
	}

	if elem.Type().Kind() != K.Struct {
		t.Errorf("Expected Kind %v, got %v", K.Struct, elem.Type().Kind())
	}

	// Test accessing fields
	numFields, err := elem.NumField()
	if err != nil {
		t.Fatalf("NumField() failed: %v", err)
	}

	if numFields != 2 {
		t.Errorf("Expected 2 fields, got %d", numFields)
	}
}
