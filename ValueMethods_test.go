package tinyreflect_test

import (
	"testing"

	"github.com/cdvelop/tinyreflect"
)

func TestLenAndCap(t *testing.T) {
	tr := tinyreflect.New()
	testCases := []struct {
		name       string
		value      interface{}
		len        int
		cap        int
		lenWantErr bool
		capWantErr bool
	}{
		{"Slice", []int{1, 2, 3}, 3, 3, false, false},
		{"Slice with Cap", make([]int, 3, 5), 3, 5, false, false},
		{"Array", [3]int{1, 2, 3}, 3, 3, false, false},
		{"String", "hello", 5, 0, false, true}, // Len is ok, Cap is not.
		{"Int", 123, 0, 0, true, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := tr.ValueOf(tc.value)

			// Test Len()
			length, err := v.Len()
			if (err != nil) != tc.lenWantErr {
				t.Errorf("Len() error = %v, wantErr %v", err, tc.lenWantErr)
			}
			if !tc.lenWantErr && length != tc.len {
				t.Errorf("Len() = %v, want %v", length, tc.len)
			}

			// Test Cap()
			capacity, err := v.Cap()
			if (err != nil) != tc.capWantErr {
				t.Errorf("Cap() error = %v, wantErr %v", err, tc.capWantErr)
			}
			if !tc.capWantErr && capacity != tc.cap {
				t.Errorf("Cap() = %v, want %v", capacity, tc.cap)
			}
		})
	}
}

func TestIndex(t *testing.T) {
	tr := tinyreflect.New()
	slice := []int{10, 20, 30}
	vSlice := tr.ValueOf(slice)
	elem, err := vSlice.Index(1)
	if err != nil {
		t.Fatalf("Index failed for slice: %v", err)
	}
	if val, _ := elem.Int(); val != 20 {
		t.Errorf("Index for slice: expected 20, got %d", val)
	}

	_, err = vSlice.Index(3)
	if err == nil {
		t.Error("Index out of range: expected an error, but got nil")
	}
}

func TestIsNil(t *testing.T) {
	tr := tinyreflect.New()
	var nilSlice []int
	var nilPtr *int
	nonNilSlice := []int{1}
	i := 123
	nonNilPtr := &i

	testCases := []struct {
		name    string
		value   interface{}
		isNil   bool
		wantErr bool
	}{
		{"Nil Slice", nilSlice, true, false},
		{"Non-nil Slice", nonNilSlice, false, false},
		{"Nil Pointer", nilPtr, true, false},
		{"Non-nil Pointer", nonNilPtr, false, false},
		{"Int", 123, false, true},
		{"Nil Interface", (interface{})(nil), true, true}, // IsNil is not for interface value itself
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := tr.ValueOf(tc.value)
			isNil, err := v.IsNil()
			if (err != nil) != tc.wantErr {
				t.Errorf("IsNil() error = %v, wantErr %v", err, tc.wantErr)
			}
			if !tc.wantErr && isNil != tc.isNil {
				t.Errorf("IsNil() = %v, want %v", isNil, tc.isNil)
			}
		})
	}
}

func TestAddr(t *testing.T) {
	tr := tinyreflect.New()
	i := 123
	v := tr.ValueOf(&i)
	elem, _ := v.Elem()
	addr, err := elem.Addr()
	if err != nil {
		t.Fatalf("Addr failed for addressable value: %v", err)
	}
	if addr.Kind().String() != "ptr" {
		t.Errorf("Addr: expected kind Pointer, got %s", addr.Kind())
	}

	v = tr.ValueOf(123)
	_, err = v.Addr()
	if err == nil {
		t.Error("Addr for non-addressable value: expected an error, but got nil")
	}
}

func TestSet(t *testing.T) {
	tr := tinyreflect.New()
	i1, i2 := 123, 456
	s1, s2 := "hello", "world"

	v1 := tr.ValueOf(&i1)
	v2 := tr.ValueOf(&i2)
	v3 := tr.ValueOf(&s1)
	v4 := tr.ValueOf(&s2)

	elem1, _ := v1.Elem()
	elem2, _ := v2.Elem()
	elem3, _ := v3.Elem()
	elem4, _ := v4.Elem()

	// Test compatible types
	if err := elem1.Set(elem2); err != nil {
		t.Fatalf("Set failed for compatible types: %v", err)
	}
	if i1 != 456 {
		t.Errorf("Set for compatible types: expected %d, got %d", 456, i1)
	}

	// Test incompatible types
	if err := elem1.Set(elem3); err == nil {
		t.Error("Set for incompatible types: expected an error, but got nil")
	}

	// Test setting string
	if err := elem3.Set(elem4); err != nil {
		t.Fatalf("Set failed for strings: %v", err)
	}
	if s1 != "world" {
		t.Errorf("Set for strings: expected 'world', got %s", s1)
	}
}