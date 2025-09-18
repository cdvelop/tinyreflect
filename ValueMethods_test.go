package tinyreflect

import (
	"testing"

	. "github.com/cdvelop/tinystring"
)

func TestLen(t *testing.T) {
	// Test with a slice
	slice := []int{1, 2, 3}
	vSlice := ValueOf(slice)
	lenSlice, err := vSlice.Len()
	if err != nil {
		t.Errorf("Len failed for slice: %v", err)
	}
	if lenSlice != 3 {
		t.Errorf("Len for slice: expected 3, got %d", lenSlice)
	}

	// Test with an array
	arr := [3]int{1, 2, 3}
	vArr := ValueOf(arr)
	lenArr, err := vArr.Len()
	if err != nil {
		t.Errorf("Len failed for array: %v", err)
	}
	if lenArr != 3 {
		t.Errorf("Len for array: expected 3, got %d", lenArr)
	}

	// Test with a string
	str := "hello"
	vStr := ValueOf(str)
	lenStr, err := vStr.Len()
	if err != nil {
		t.Errorf("Len failed for string: %v", err)
	}
	if lenStr != 5 {
		t.Errorf("Len for string: expected 5, got %d", lenStr)
	}

	// Test with an invalid type
	i := 123
	vInt := ValueOf(i)
	_, err = vInt.Len()
	if err == nil {
		t.Error("Len for int: expected an error, but got nil")
	}

	// Test with an array with nil type
	arrWithNilType := [3]int{1, 2, 3}
	vArrWithNilType := ValueOf(arrWithNilType)
	vArrWithNilType.typ_ = nil // Manually set typ_ to nil
	_, err = vArrWithNilType.Len()
	if err == nil {
		t.Error("Len for array with nil type: expected an error, but got nil")
	}
}

func TestCap(t *testing.T) {
	// Test with a slice
	slice := make([]int, 3, 5)
	vSlice := ValueOf(slice)
	capSlice, err := vSlice.Cap()
	if err != nil {
		t.Errorf("Cap failed for slice: %v", err)
	}
	if capSlice != 5 {
		t.Errorf("Cap for slice: expected 5, got %d", capSlice)
	}

	// Test with an array
	arr := [3]int{1, 2, 3}
	vArr := ValueOf(arr)
	capArr, err := vArr.Cap()
	if err != nil {
		t.Errorf("Cap failed for array: %v", err)
	}
	if capArr != 3 {
		t.Errorf("Cap for array: expected 3, got %d", capArr)
	}

	// Test with an invalid type
	i := 123
	vInt := ValueOf(i)
	_, err = vInt.Cap()
	if err == nil {
		t.Error("Cap for int: expected an error, but got nil")
	}

	// Test with an array with nil type
	arrWithNilType := [3]int{1, 2, 3}
	vArrWithNilType := ValueOf(arrWithNilType)
	vArrWithNilType.typ_ = nil // Manually set typ_ to nil
	_, err = vArrWithNilType.Cap()
	if err == nil {
		t.Error("Cap for array with nil type: expected an error, but got nil")
	}
}

func TestIndex(t *testing.T) {
	// Test with a slice
	slice := []int{10, 20, 30}
	vSlice := ValueOf(slice)
	elem, err := vSlice.Index(1)
	if err != nil {
		t.Errorf("Index failed for slice: %v", err)
	}
	if val, _ := elem.Int(); val != 20 {
		t.Errorf("Index for slice: expected 20, got %d", val)
	}

	// Test with an array
	arr := [3]int{10, 20, 30}
	vArr := ValueOf(arr)
	elem, err = vArr.Index(2)
	if err != nil {
		t.Errorf("Index failed for array: %v", err)
	}
	if val, _ := elem.Int(); val != 30 {
		t.Errorf("Index for array: expected 30, got %d", val)
	}

	// Test with a string
	str := "abc"
	vStr := ValueOf(str)
	elem, err = vStr.Index(0)
	if err != nil {
		t.Errorf("Index failed for string: %v", err)
	}
	if val, _ := elem.Uint(); val != 'a' {
		t.Errorf("Index for string: expected 'a', got %c", val)
	}

	// Test out of range
	_, err = vSlice.Index(3)
	if err == nil {
		t.Error("Index out of range: expected an error, but got nil")
	}

	// Test with invalid type
	i := 123
	vInt := ValueOf(i)
	_, err = vInt.Index(0)
	if err == nil {
		t.Error("Index for int: expected an error, but got nil")
	}

	// Test with array with nil type
	vArrWithNilType := ValueOf(arr)
	vArrWithNilType.typ_ = nil
	_, err = vArrWithNilType.Index(0)
	if err == nil {
		t.Error("Index for array with nil type: expected an error, but got nil")
	}

	// Test with slice with nil type
	vSliceWithNilType := ValueOf(slice)
	vSliceWithNilType.typ_ = nil
	_, err = vSliceWithNilType.Index(0)
	if err == nil {
		t.Error("Index for slice with nil type: expected an error, but got nil")
	}

	// Test string out of range
	str = "abc"
	vStr = ValueOf(str)
	_, err = vStr.Index(3)
	if err == nil {
		t.Error("Index string out of range: expected an error, but got nil")
	}
}

func TestIsNil(t *testing.T) {
	// Test with a nil slice
	var slice []int
	vSlice := ValueOf(slice)
	isNil, err := vSlice.IsNil()
	if err != nil {
		t.Errorf("IsNil failed for nil slice: %v", err)
	}
	if !isNil {
		t.Error("IsNil for nil slice: expected true, got false")
	}

	// Test with a non-nil slice
	slice = []int{1, 2, 3}
	vSlice = ValueOf(slice)
	isNil, err = vSlice.IsNil()
	if err != nil {
		t.Errorf("IsNil failed for non-nil slice: %v", err)
	}
	if isNil {
		t.Error("IsNil for non-nil slice: expected false, got true")
	}

	// Test with a nil pointer
	var ptr *int
	vPtr := ValueOf(ptr)
	isNil, err = vPtr.IsNil()
	if err != nil {
		t.Errorf("IsNil failed for nil pointer: %v", err)
	}
	if !isNil {
		t.Error("IsNil for nil pointer: expected true, got false")
	}

	// Test with a non-nil pointer
	i := 123
	ptr = &i
	vPtr = ValueOf(ptr)
	isNil, err = vPtr.IsNil()
	if err != nil {
		t.Errorf("IsNil failed for non-nil pointer: %v", err)
	}
	if isNil {
		t.Error("IsNil for non-nil pointer: expected false, got true")
	}

	// Test with a nil interface
	var iface interface{}
	vIface := ValueOf(iface)
	_, err = vIface.IsNil()
	if err == nil {
		t.Error("IsNil for nil interface: expected an error, but got nil")
	}

	// Test with a non-nil interface that holds a non-pointer value
	iface = 123
	vIface = ValueOf(iface)
	_, err = vIface.IsNil()
	if err == nil {
		t.Error("IsNil for non-nil interface: expected an error, but got nil")
	}

	// Test with an invalid type
	vInt := ValueOf(123)
	_, err = vInt.IsNil()
	if err == nil {
		t.Error("IsNil for int: expected an error, but got nil")
	}

	// Test with an indirect pointer
	var p *int
	vIndir := ValueOf(&p)
	elem, _ := vIndir.Elem()
	isNil, err = elem.IsNil()
	if err != nil {
		t.Errorf("IsNil failed for indirect pointer: %v", err)
	}
	if !isNil {
		t.Error("IsNil for indirect pointer: expected true, got false")
	}

	// Test with a non-nil indirect pointer
	i = 123
	p = &i
	vIndir = ValueOf(&p)
	elem, _ = vIndir.Elem()
	isNil, err = elem.IsNil()
	if err != nil {
		t.Errorf("IsNil failed for non-nil indirect pointer: %v", err)
	}
	if isNil {
		t.Error("IsNil for non-nil indirect pointer: expected false, got true")
	}
}

func TestAddr(t *testing.T) {
	// Test with an addressable value
	i := 123
	v := ValueOf(&i)
	elem, _ := v.Elem()
	addr, err := elem.Addr()
	if err != nil {
		t.Errorf("Addr failed for addressable value: %v", err)
	}
	if addr.Kind() != K.Pointer {
		t.Errorf("Addr for addressable value: expected kind Pointer, got %s", addr.Kind())
	}

	// Test with a non-addressable value
	v = ValueOf(123)
	_, err = v.Addr()
	if err == nil {
		t.Error("Addr for non-addressable value: expected an error, but got nil")
	}

	// Test with a nil type
	v = ValueOf(&i)
	elem, _ = v.Elem()
	elem.typ_ = nil
	_, err = elem.Addr()
	if err == nil {
		t.Error("Addr with nil type: expected an error, but got nil")
	}
}

func TestSet(t *testing.T) {
	// Test with compatible types
	i1, i2 := 123, 456
	v1 := ValueOf(&i1)
	v2 := ValueOf(&i2)
	elem1, _ := v1.Elem()
	elem2, _ := v2.Elem()
	err := elem1.Set(elem2)
	if err != nil {
		t.Errorf("Set failed for compatible types: %v", err)
	}
	if i1 != 456 {
		t.Errorf("Set for compatible types: expected %d, got %d", 456, i1)
	}

	// Test with incompatible types
	s := "hello"
	v3 := ValueOf(&s)
	elem3, _ := v3.Elem()
	err = elem1.Set(elem3)
	if err == nil {
		t.Error("Set for incompatible types: expected an error, but got nil")
	}

	// Test with zero values
	var z1, z2 int
	vZ1 := ValueOf(&z1)
	vZ2 := ValueOf(z2)
	elemZ1, _ := vZ1.Elem()
	err = elemZ1.Set(vZ2)
	if err != nil {
		t.Errorf("Set with zero value: %v", err)
	}

	// Test with slices
	slice1 := []int{1, 2}
	slice2 := []int{3, 4}
	vSlice1 := ValueOf(&slice1)
	vSlice2 := ValueOf(&slice2)
	elemSlice1, _ := vSlice1.Elem()
	elemSlice2, _ := vSlice2.Elem()
	err = elemSlice1.Set(elemSlice2)
	if err != nil {
		t.Errorf("Set failed for slices: %v", err)
	}
	if slice1[0] != 3 {
		t.Errorf("Set for slices: expected %d, got %d", 3, slice1[0])
	}

	// Test with pointers
	p1 := &i1
	p2 := &i2
	vP1 := ValueOf(&p1)
	vP2 := ValueOf(&p2)
	elemP1, _ := vP1.Elem()
	elemP2, _ := vP2.Elem()
	err = elemP1.Set(elemP2)
	if err != nil {
		t.Errorf("Set failed for pointers: %v", err)
	}
	if *p1 != *p2 {
		t.Errorf("Set for pointers: expected %d, got %d", *p2, *p1)
	}
}

func TestSetErrors(t *testing.T) {
	// Test with nil type in destination
	i1 := 123
	v1 := ValueOf(&i1)
	elem1, _ := v1.Elem()
	elem1.typ_ = nil
	err := elem1.Set(elem1)
	if err == nil {
		t.Error("Set with nil type in destination: expected an error, but got nil")
	}

	// Test with nil type in source
	v1 = ValueOf(&i1)
	elem1, _ = v1.Elem()
	v2 := ValueOf(&i1)
	elem2, _ := v2.Elem()
	elem2.typ_ = nil
	err = elem1.Set(elem2)
	if err == nil {
		t.Error("Set with nil type in source: expected an error, but got nil")
	}

	// Test with incompatible slice types
	slice1 := []int{1, 2}
	slice2 := []byte{3, 4}
	vSlice1 := ValueOf(&slice1)
	vSlice2 := ValueOf(&slice2)
	elemSlice1, _ := vSlice1.Elem()
	elemSlice2, _ := vSlice2.Elem()
	err = elemSlice1.Set(elemSlice2)
	if err == nil {
		t.Error("Set with incompatible slice types: expected an error, but got nil")
	}

	// Test with incompatible pointer types
	s := "hello"
	p1 := &i1
	p2 := &s
	vP1 := ValueOf(&p1)
	vP2 := ValueOf(&p2)
	elemP1, _ := vP1.Elem()
	elemP2, _ := vP2.Elem()
	err = elemP1.Set(elemP2)
	if err == nil {
		t.Error("Set with incompatible pointer types: expected an error, but got nil")
	}
}

func TestSetZeroSizeAndNonIndir(t *testing.T) {
	// Test with zero-sized type
	type zero struct{}
	z1 := zero{}
	z2 := zero{}
	vZ1 := ValueOf(&z1)
	vZ2 := ValueOf(&z2)
	elemZ1, _ := vZ1.Elem()
	elemZ2, _ := vZ2.Elem()
	err := elemZ1.Set(elemZ2)
	if err != nil {
		t.Errorf("Set with zero-sized type: %v", err)
	}

	// Test with non-indir source and dest
	i1 := 123
	i2 := 456
	v1 := ValueOf(i1)
	v2 := ValueOf(i2)
	v1.flag |= flagAddr  // make it addressable
	v2.flag |= flagAddr  // make it addressable
	err = v1.Set(v2)
	if err != nil {
		t.Errorf("Set with non-indir values: %v", err)
	}
}
