package tinyreflect

import (
	"testing"
	"unsafe"

	. "github.com/cdvelop/tinystring"
)

func TestSliceType(t *testing.T) {
	// Test with a slice of ints
	slice := []int{}
	v := ValueOf(slice)
	st := (*SliceType)(unsafe.Pointer(v.Type()))

	// Test Element
	elem := st.Element()
	if elem.Kind() != K.Int {
		t.Errorf("Element: expected kind Int, got %s", elem.Kind())
	}

}

func TestArrayType(t *testing.T) {
	// Test with an array of ints
	arr := [3]int{}
	v := ValueOf(arr)
	at := (*ArrayType)(unsafe.Pointer(v.Type()))

	// Test Element
	elem := at.Element()
	if elem.Kind() != K.Int {
		t.Errorf("Element: expected kind Int, got %s", elem.Kind())
	}

	// Test Length
	if at.Length() != 3 {
		t.Errorf("Length: expected 3, got %d", at.Length())
	}
}
