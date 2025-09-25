package tinyreflect_test

import (
	"testing"

	"github.com/cdvelop/tinyreflect"
)

func TestSliceType(t *testing.T) {
	slice := []int{}
	typ := tinyreflect.TypeOf(slice)

	st := typ.SliceType()
	if st == nil {
		t.Fatal("SliceType() returned nil for a slice type")
	}

	elem := st.Element()
	if elem.Kind().String() != "int" {
		t.Errorf("Element: expected kind 'int', got '%s'", elem.Kind())
	}
}

func TestArrayType(t *testing.T) {
	arr := [3]int{}
	typ := tinyreflect.TypeOf(arr)

	at := typ.ArrayType()
	if at == nil {
		t.Fatal("ArrayType() returned nil for an array type")
	}

	elem := at.Element()
	if elem.Kind().String() != "int" {
		t.Errorf("Element: expected kind 'int', got '%s'", elem.Kind())
	}

	if at.Len != 3 {
		t.Errorf("Len: expected 3, got %d", at.Len)
	}
}