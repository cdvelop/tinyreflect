package tinyreflect

import (
	"testing"

	. "github.com/cdvelop/tinystring"
)

func TestMakeSlice(t *testing.T) {
	// Test success case
	sliceType := TypeOf([]int{})
	v, err := MakeSlice(sliceType, 5, 10)
	if err != nil {
		t.Errorf("MakeSlice success: unexpected error: %v", err)
	}
	if v.Kind() != K.Slice {
		t.Errorf("MakeSlice success: expected kind Slice, got %s", v.Kind())
	}
	if l, _ := v.Len(); l != 5 {
		t.Errorf("MakeSlice success: expected len 5, got %d", l)
	}
	if c, _ := v.Cap(); c != 10 {
		t.Errorf("MakeSlice success: expected cap 10, got %d", c)
	}

	// Test nil type
	_, err = MakeSlice(nil, 0, 0)
	if err == nil {
		t.Error("MakeSlice with nil type: expected an error, but got nil")
	}

	// Test non-slice type
	intType := TypeOf(0)
	_, err = MakeSlice(intType, 0, 0)
	if err == nil {
		t.Error("MakeSlice with non-slice type: expected an error, but got nil")
	}

	// Test negative len
	_, err = MakeSlice(sliceType, -1, 0)
	if err == nil {
		t.Error("MakeSlice with negative len: expected an error, but got nil")
	}

	// Test negative cap
	_, err = MakeSlice(sliceType, 0, -1)
	if err == nil {
		t.Error("MakeSlice with negative cap: expected an error, but got nil")
	}

	// Test len > cap
	_, err = MakeSlice(sliceType, 1, 0)
	if err == nil {
		t.Error("MakeSlice with len > cap: expected an error, but got nil")
	}

	// Test zero-sized element
	zeroSliceType := TypeOf([]struct{}{})
	v, err = MakeSlice(zeroSliceType, 5, 10)
	if err != nil {
		t.Errorf("MakeSlice with zero-sized element: unexpected error: %v", err)
	}
}
