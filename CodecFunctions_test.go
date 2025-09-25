package tinyreflect_test

import (
	"testing"

	"github.com/cdvelop/tinyreflect"
)

func TestMakeSlice(t *testing.T) {
	tr := tinyreflect.New()

	// Test success case
	sliceType := tr.TypeOf([]int{})
	v, err := tr.MakeSlice(sliceType, 5, 10)
	if err != nil {
		t.Errorf("MakeSlice success: unexpected error: %v", err)
	}
	if v.Kind().String() != "slice" {
		t.Errorf("MakeSlice success: expected kind Slice, got %s", v.Kind())
	}
	if l, _ := v.Len(); l != 5 {
		t.Errorf("MakeSlice success: expected len 5, got %d", l)
	}
	if c, _ := v.Cap(); c != 10 {
		t.Errorf("MakeSlice success: expected cap 10, got %d", c)
	}

	// Test nil type
	_, err = tr.MakeSlice(nil, 0, 0)
	if err == nil {
		t.Error("MakeSlice with nil type: expected an error, but got nil")
	}

	// Test non-slice type
	intType := tr.TypeOf(0)
	_, err = tr.MakeSlice(intType, 0, 0)
	if err == nil {
		t.Error("MakeSlice with non-slice type: expected an error, but got nil")
	}
}

func TestNewValue(t *testing.T) {
	tr := tinyreflect.New()
	typ := tr.TypeOf(0) // type int

	v := tr.NewValue(typ)
	if v.Kind().String() != "ptr" {
		t.Fatalf("NewValue should return a pointer, but got %s", v.Kind())
	}

	elem, err := v.Elem()
	if err != nil {
		t.Fatalf("Elem() on NewValue result failed: %v", err)
	}

	if elem.Kind().String() != "int" {
		t.Errorf("NewValue's element should be Int, but got %s", elem.Kind())
	}

	if !elem.IsZero() {
		t.Error("NewValue should point to a zero value")
	}
}

func TestIndirect(t *testing.T) {
	tr := tinyreflect.New()

	// Test with a non-pointer
	vInt := tr.ValueOf(123)
	indirectVInt := tr.Indirect(vInt)
	val, _ := indirectVInt.Int()
	if val != 123 {
		t.Errorf("Indirect on non-pointer should return the same value, got %d", val)
	}

	// Test with a pointer
	i := 456
	vPtr := tr.ValueOf(&i)
	indirectVPrt := tr.Indirect(vPtr)
	if indirectVPrt.Kind().String() != "int" {
		t.Errorf("Indirect on pointer should return the element's kind, got %s", indirectVPrt.Kind())
	}
	val, _ = indirectVPrt.Int()
	if val != 456 {
		t.Errorf("Indirect on pointer should return the element's value, got %d", val)
	}

	// Test with a nil pointer
	var nilPtr *int
	vNilPtr := tr.ValueOf(nilPtr)
	indirectVNilPtr := tr.Indirect(vNilPtr)
	if !indirectVNilPtr.IsZero() {
		t.Error("Indirect on a nil pointer should return a zero value")
	}
}