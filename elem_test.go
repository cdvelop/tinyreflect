package tinyreflect_test

import (
	"testing"

	"github.com/cdvelop/tinyreflect"
)

func TestElem(t *testing.T) {
	tr := tinyreflect.New()

	// Test basic pointer dereference
	x := 42
	ptr := &x
	v := tr.ValueOf(ptr)
	elem, err := v.Elem()
	if err != nil {
		t.Fatalf("Elem() failed: %v", err)
	}
	if elem.Type() == nil {
		t.Fatal("Elem() returned nil type")
	}
	if elem.Kind().String() != "int" {
		t.Errorf("Expected Kind 'int', got '%s'", elem.Kind())
	}

	// Test nil pointer
	var nilPtr *int
	vNil := tr.ValueOf(nilPtr)
	elemNil, err := vNil.Elem()
	if err != nil {
		t.Fatalf("Elem() on nil pointer failed: %v", err)
	}
	if !elemNil.IsZero() {
		t.Error("Expected zero value for nil pointer elem")
	}

	// Test error case - non-pointer
	vInt := tr.ValueOf(42)
	_, err = vInt.Elem()
	if err == nil {
		t.Error("Expected error when calling Elem() on non-pointer")
	}
}

func TestElemStruct(t *testing.T) {
	tr := tinyreflect.New()
	type TestStruct struct {
		X int
		Y string
	}

	s := &TestStruct{X: 10, Y: "hello"}
	v := tr.ValueOf(s)
	elem, err := v.Elem()
	if err != nil {
		t.Fatalf("Elem() failed: %v", err)
	}

	if elem.Kind().String() != "struct" {
		t.Errorf("Expected Kind 'struct', got '%s'", elem.Kind())
	}

	numFields, err := elem.NumField()
	if err != nil {
		t.Fatalf("NumField() failed: %v", err)
	}
	if numFields != 2 {
		t.Errorf("Expected 2 fields, got %d", numFields)
	}
}